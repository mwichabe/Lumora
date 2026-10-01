import 'dart:async';
import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:lumora_mobile/core/network/dio_client.dart';
import 'package:lumora_mobile/core/router/app_router.dart';
import 'package:lumora_mobile/core/storage/session_storage.dart';
import 'package:lumora_mobile/main.dart';
import 'package:lumora_mobile/models/user.dart';
import 'package:lumora_mobile/providers/auth_provider.dart';

/// Answers every request with 401, each one only when the test releases it —
/// so a response can be made to arrive after the session has changed.
class _Held401 implements HttpClientAdapter {
  final pending = <({String? auth, Completer<void> release})>[];

  @override
  Future<ResponseBody> fetch(RequestOptions options, Stream<Uint8List>? body, Future<void>? cancel) async {
    final release = Completer<void>();
    pending.add((auth: options.headers['Authorization'] as String?, release: release));
    await release.future;
    return ResponseBody.fromString(jsonEncode({'error': 'invalid token'}), 401,
        headers: {Headers.contentTypeHeader: [Headers.jsonContentType]});
  }

  @override
  void close({bool force = false}) {}
}

final _user = User.fromJson({'id': 1, 'email': 'a@b.c', 'name': 'A', 'targetLanguage': 'de'});

void main() {
  late Map<String, String> secure;

  setUp(() {
    SharedPreferences.setMockInitialValues({});
    secure = {};
    TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
      const MethodChannel('plugins.it_nomads.com/flutter_secure_storage'),
      (call) async {
        final args = (call.arguments as Map).cast<String, dynamic>();
        switch (call.method) {
          case 'write':
            secure[args['key'] as String] = args['value'] as String;
          case 'delete':
            secure.remove(args['key']);
          case 'read':
            return secure[args['key']];
        }
        return null;
      },
    );
  });

  group('401 handling', () {
    late _Held401 server;
    late int signedOut;

    setUp(() async {
      server = _Held401();
      DioClient.instance.dio.httpClientAdapter = server;
      signedOut = 0;
      DioClient.instance.onUnauthorized = () => signedOut++;
      await SessionStorage.instance.clearSession();
    });

    Future<void> request() => DioClient.instance.dio.get('/api/home').then((_) {}, onError: (_) {});

    Future<void> waitForRequests(int n) async {
      while (server.pending.length < n) {
        await Future<void>.delayed(Duration.zero);
      }
    }

    test('a 401 for the current session signs the user out', () async {
      await SessionStorage.instance.setToken('token-A');
      final done = request();
      await waitForRequests(1);
      server.pending.single.release.complete();
      await done;

      expect(signedOut, 1);
      expect(await SessionStorage.instance.getToken(), isNull);
    });

    test('a 401 that arrives after signing back in leaves the new session alone', () async {
      // Signed in as A, a request goes out...
      await SessionStorage.instance.setToken('token-A');
      final stale = request();
      // ...and a background poll fires while signed out...
      await SessionStorage.instance.clearSession();
      final signedOutPoll = request();
      await waitForRequests(2);
      expect(server.pending[1].auth, isNull);

      // ...then the user signs in again before either answer arrives.
      await SessionStorage.instance.setToken('token-B');
      for (final p in server.pending) {
        p.release.complete();
      }
      await Future.wait([stale, signedOutPoll]);

      expect(signedOut, 0, reason: 'neither 401 was about the current session');
      expect(await SessionStorage.instance.getToken(), 'token-B');
      expect(secure['lumora_token'], 'token-B', reason: 'the new token must survive on disk too');
    });

    test('a wrong password on the login form is not a session expiry', () async {
      final done = request(); // no token: like POST /api/auth/login
      await waitForRequests(1);
      server.pending.single.release.complete();
      await done;
      expect(signedOut, 0);
    });
  });

  testWidgets('an expired session lands on sign-in, not back on home', (tester) async {
    tester.view.physicalSize = const Size(1200, 2600);
    tester.view.devicePixelRatio = 1;
    addTearDown(tester.view.reset);
    // The test font draws every glyph as a full-width block, so the welcome
    // screen's long button labels overflow here (they don't with real fonts).
    // Layout isn't what this test is about.
    final onError = FlutterError.onError;
    FlutterError.onError = (details) {
      if (!details.toString().contains('A RenderFlex overflowed')) onError?.call(details);
    };
    addTearDown(() => FlutterError.onError = onError);

    final container = ProviderContainer(overrides: [authProvider.overrideWith(_SignedIn.new)]);
    addTearDown(container.dispose);
    await tester.pumpWidget(UncontrolledProviderScope(container: container, child: const LumoraApp()));
    final router = container.read(routerProvider);
    router.go('/learn');
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }

    DioClient.instance.onUnauthorized!();
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }

    expect(container.read(authProvider).isAuthenticated, isFalse);
    expect(_location(router), '/welcome');
    await tester.pumpWidget(const SizedBox());
  });
}

String _location(GoRouter r) => r.routerDelegate.currentConfiguration.last.matchedLocation;

class _SignedIn extends AuthController {
  @override
  AuthState build() => AuthState(user: _user, loading: false);
}
