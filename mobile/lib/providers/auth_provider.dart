import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../core/network/api_exception.dart';
import '../core/storage/session_storage.dart';
import '../models/user.dart';

class AuthState {
  final User? user;
  final bool loading;
  const AuthState({required this.user, required this.loading});

  AuthState copyWith({User? user, bool? loading, bool clearUser = false}) => AuthState(
        user: clearUser ? null : (user ?? this.user),
        loading: loading ?? this.loading,
      );

  bool get isAuthenticated => user != null;
}

/// Mirrors frontend/lib/auth.tsx: seeds from the cached user for an instant,
/// flash-free start, then revalidates against the server in the background.
/// Only a genuine 401 clears the session — any other failure (offline, 5xx)
/// keeps the cached user so the app stays usable.
class AuthController extends Notifier<AuthState> {
  @override
  AuthState build() {
    _bootstrap();
    return const AuthState(user: null, loading: true);
  }

  Future<void> _bootstrap() async {
    final cached = await SessionStorage.instance.getStoredUser();
    // A saved user plus a token is enough to open the app: stop "loading" now
    // so the splash hands over immediately, and let refresh() revalidate in the
    // background. Waiting on it here held the splash for a full round trip —
    // close to a minute when the API is waking from idle. If the token turns
    // out to be dead, refresh() clears the session and the router sends the
    // user to sign in.
    if (cached != null && await SessionStorage.instance.getToken() != null) {
      state = AuthState(user: cached, loading: false);
    }
    await refresh();
  }

  Future<void> refresh() async {
    final token = await SessionStorage.instance.getToken();
    if (token == null) {
      state = const AuthState(user: null, loading: false);
      return;
    }
    try {
      final user = await ApiClient.instance.me();
      await SessionStorage.instance.setStoredUser(user);
      state = state.copyWith(user: user, loading: false);
    } on ApiException catch (e) {
      if (e.status == 401) {
        await SessionStorage.instance.clearSession();
        state = const AuthState(user: null, loading: false);
      } else {
        state = state.copyWith(loading: false);
      }
    } catch (_) {
      state = state.copyWith(loading: false);
    }
  }

  Future<void> login(String email, String password) async {
    final (token, user) = await ApiClient.instance.login(email, password);
    await SessionStorage.instance.setToken(token);
    await SessionStorage.instance.setStoredUser(user);
    state = state.copyWith(user: user, loading: false);
  }

  Future<void> register(String email, String password, String name) async {
    final (token, user) = await ApiClient.instance.register(email, password, name);
    await SessionStorage.instance.setToken(token);
    await SessionStorage.instance.setStoredUser(user);
    state = state.copyWith(user: user, loading: false);
  }

  Future<void> logout() async {
    await SessionStorage.instance.clearSession();
    state = const AuthState(user: null, loading: false);
  }

  void setUser(User user) {
    SessionStorage.instance.setStoredUser(user);
    state = state.copyWith(user: user);
  }
}

final authProvider = NotifierProvider<AuthController, AuthState>(AuthController.new);
