import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:lumora_mobile/core/network/dio_client.dart';
import 'package:lumora_mobile/core/storage/session_storage.dart';
import 'package:lumora_mobile/features/practice/practice_run_screen.dart';
import 'package:lumora_mobile/models/user.dart';
import 'package:lumora_mobile/providers/auth_provider.dart';
import 'package:lumora_mobile/providers/home_provider.dart';

final _userJson = {'id': 1, 'email': 'a@b.c', 'name': 'A', 'targetLanguage': 'de'};

/// Serves canned JSON per path and records every request it sees.
class _Api implements HttpClientAdapter {
  final requests = <String>[];

  /// The language the fake server reports the learner is on.
  String lang = 'de';
  final bodies = <String, List<dynamic>>{};

  @override
  Future<ResponseBody> fetch(RequestOptions o, Stream<Uint8List>? body, Future<void>? cancel) async {
    requests.add('${o.method} ${o.path}');
    if (o.data != null) bodies.putIfAbsent(o.path, () => []).add(o.data);
    final Object json = switch (o.path) {
      '/api/practice/reading' => {
          'sessions': [
            {
              'id': 7,
              'title': 'Im Café',
              'lines': [
                {'id': 1, 'text': 'Anna trinkt jeden Morgen einen Kaffee.'},
                {'id': 2, 'text': 'Heute bestellt sie einen Tee.'},
              ],
              'questions': [
                {'question': 'What does Anna order today?', 'options': ['Coffee', 'Tea'], 'correctAnswer': 'Tea'},
                {'question': 'When does Anna drink coffee?', 'options': ['Every morning', 'At night'], 'correctAnswer': 'Every morning'},
              ],
            },
          ],
        },
      '/api/practice' => {
          'vocab': [],
          'mistakes': [
            {'id': 3, 'prompt': 'Listen and choose the meaning', 'question': 'der Mann', 'correctAnswer': 'the man'},
          ],
          'listeningCount': 0,
          'readingCount': 0,
        },
      '/api/practice/complete' => {'xpEarned': 12, 'user': _userJson},
      '/api/home' => {'user': {..._userJson, 'targetLanguage': lang}, 'quests': []},
      _ => <String, dynamic>{},
    };
    return ResponseBody.fromString(jsonEncode(json), 200, headers: {Headers.contentTypeHeader: [Headers.jsonContentType]});
  }

  @override
  void close({bool force = false}) {}
}

class _SignedIn extends AuthController {
  @override
  AuthState build() => AuthState(user: User.fromJson(_userJson), loading: false);
}

Future<_Api> _pumpRun(WidgetTester tester, String mode) async {
  tester.view.physicalSize = const Size(400, 900);
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
  final api = _Api();
  DioClient.instance.dio.httpClientAdapter = api;
  // The token lookup is cached in memory. Settle it inside this test's fake
  // clock: one cached by another test (or by setUp, which runs on the real
  // clock) is never delivered here, and requests would silently stall.
  SessionStorage.instance.clearSession();
  await tester.pumpWidget(ProviderScope(
    overrides: [authProvider.overrideWith(_SignedIn.new)],
    child: MaterialApp(home: PracticeRunScreen(mode: mode)),
  ));
  for (var i = 0; i < 5; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }
  return api;
}

void main() {
  setUp(() async {
    SharedPreferences.setMockInitialValues({});
    final messenger = TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;
    messenger.setMockMethodCallHandler(const MethodChannel('plugins.it_nomads.com/flutter_secure_storage'), (_) async => null);
    messenger.setMockMethodCallHandler(const MethodChannel('flutter_tts'), (_) async => 1);
    await SessionStorage.instance.clearSession();
  });

  testWidgets('reading comprehension shows the passage, not just the questions', (tester) async {
    final api = await _pumpRun(tester, 'reading');

    expect(find.text('Im Café'), findsOneWidget);
    expect(find.text('Anna trinkt jeden Morgen einen Kaffee.'), findsOneWidget);
    expect(find.text('Heute bestellt sie einen Tee.'), findsOneWidget);
    expect(find.text('1. What does Anna order today?'), findsOneWidget);

    // Answer one right and one wrong, then submit.
    await tester.tap(find.text('Tea'));
    await tester.pump();
    await tester.scrollUntilVisible(find.text('At night'), 100);
    await tester.tap(find.text('At night'));
    await tester.pump();
    await tester.tap(find.text('Submit answers'));
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }

    // The miss is saved for Review Mistakes and recapped before the score.
    expect(api.bodies['/api/mistakes'], hasLength(1));
    expect((api.bodies['/api/mistakes']!.single as Map)['question'], 'When does Anna drink coffee?');
    expect(find.text('Review your mistakes'), findsOneWidget);
    await tester.tap(find.text('See my score'));
    await tester.pump();
    expect(find.text('50%'), findsOneWidget);
  });

  testWidgets('a saved listening mistake shows its question when reviewed', (tester) async {
    await _pumpRun(tester, 'mistakes');
    // It used to be hidden (and read aloud) because its prompt mentions "listen".
    expect(find.text('der Mann'), findsOneWidget);
    expect(find.text('the man'), findsOneWidget);
  });

  test('home refetches when the learner switches language', () async {
    final api = _Api();
    DioClient.instance.dio.httpClientAdapter = api;
    final container = ProviderContainer(overrides: [authProvider.overrideWith(_SignedIn.new)]);
    addTearDown(container.dispose);
    container.listen(homeProvider, (_, _) {});

    await container.read(homeProvider.future);
    expect(api.requests.where((r) => r == 'GET /api/home'), hasLength(1));

    api.lang = 'fr';
    container.read(authProvider.notifier).setUser(User.fromJson({..._userJson, 'targetLanguage': 'fr'}));
    await container.read(homeProvider.future);
    expect(api.requests.where((r) => r == 'GET /api/home'), hasLength(2),
        reason: 'the "continue where you left off" card must follow the new language');
  });
}
