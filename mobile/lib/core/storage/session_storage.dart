import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../models/user.dart';

/// Mirrors the web client's localStorage session (frontend/lib/api.ts):
/// `lumora_token`, `lumora_user`, `lumora_last_route`. The JWT goes in secure
/// storage (Keychain/Keystore); the cached user + last route are plain prefs
/// since they hold nothing secret and are read synchronously at startup.
class SessionStorage {
  SessionStorage._();
  static final SessionStorage instance = SessionStorage._();

  static const _tokenKey = 'lumora_token';
  static const _userKey = 'lumora_user';
  static const _routeKey = 'lumora_last_route';

  final _secure = const FlutterSecureStorage();
  SharedPreferences? _prefs;

  Future<SharedPreferences> get _p async => _prefs ??= await SharedPreferences.getInstance();

  // The token is attached to every request, and a Keystore/Keychain read is a
  // slow platform-channel round trip (decryption on Android). Read it once and
  // serve it from memory afterwards; the Future is cached rather than the
  // value so requests fired together at startup share a single read.
  Future<String?>? _token;

  Future<String?> getToken() => _token ??= _secure.read(key: _tokenKey);

  Future<void> setToken(String token) {
    _token = Future.value(token);
    return _secure.write(key: _tokenKey, value: token);
  }

  Future<void> clearToken() {
    _token = Future.value(null);
    return _secure.delete(key: _tokenKey);
  }

  Future<User?> getStoredUser() async {
    final prefs = await _p;
    final raw = prefs.getString(_userKey);
    if (raw == null) return null;
    try {
      return User.fromJson(jsonDecode(raw) as Map<String, dynamic>);
    } catch (_) {
      return null;
    }
  }

  Future<void> setStoredUser(User? user) async {
    final prefs = await _p;
    if (user == null) {
      await prefs.remove(_userKey);
    } else {
      await prefs.setString(_userKey, jsonEncode(user.toJson()));
    }
  }

  Future<String?> getLastRoute() async => (await _p).getString(_routeKey);

  Future<void> setLastRoute(String path) async => (await _p).setString(_routeKey, path);

  Future<void> clearSession() async {
    await clearToken();
    await setStoredUser(null);
  }
}
