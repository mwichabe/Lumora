import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:lumora_mobile/features/learn/learn_screen.dart';
import 'package:lumora_mobile/models/lesson.dart';
import 'package:lumora_mobile/models/user.dart';
import 'package:lumora_mobile/providers/auth_provider.dart';
import 'package:lumora_mobile/providers/learn_provider.dart';

final _user = User.fromJson({'id': 1, 'email': 'learner@example.com', 'name': 'Learner', 'targetLanguage': 'de'});

/// A course with long, real-world skill titles in both units, and every state
/// a node can be in: completed, current and locked.
final _skills = [
  for (final (i, s) in [
    ('A1 · Grundlagen', 'Artikel & Nomen', true, true),
    ('A1 · Grundlagen', 'Präsens: sein & haben und die regelmäßigen Verben', true, false),
    ('A1 · Grundlagen', 'Zahlen, Uhrzeit und Datum', false, false),
    ('A2 · Aufbaustufe — Alltag, Reisen und Gespräche im Café', 'Perfekt mit haben und sein', false, false),
    ('A2 · Aufbaustufe — Alltag, Reisen und Gespräche im Café', 'Dativ', false, false),
  ].indexed)
    Skill.fromJson({
      'id': i + 1,
      'unit': s.$1,
      'title': s.$2,
      'icon': 'Hand',
      'color': '#6C3FC5',
      'unlocked': s.$3,
      'completed': s.$4,
      'lessonCount': 4,
      'completedCount': s.$4 ? 4 : 1,
      'lessons': [
        {'id': (i + 1) * 10},
      ],
    }),
];

class _SignedIn extends AuthController {
  @override
  AuthState build() => AuthState(user: _user, loading: false);
}

class _Loaded extends LearnController {
  @override
  LearnState build() => LearnState(skills: _skills, listening: const [], reading: const [], loading: false, error: false);
}

Future<void> _openRoadmap(WidgetTester tester, double width) async {
  tester.view.physicalSize = Size(width, 2400);
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);

  await tester.pumpWidget(ProviderScope(
    overrides: [
      authProvider.overrideWith(_SignedIn.new),
      learnProvider.overrideWith(_Loaded.new),
    ],
    child: const MaterialApp(home: LearnScreen()),
  ));
  await tester.tap(find.text('Roadmap'));
  await tester.pumpAndSettle();
}

void main() {
  // The narrowest phones in use are 320 logical pixels wide.
  for (final width in [320.0, 360.0, 412.0]) {
    testWidgets('roadmap fits a ${width.toInt()}px-wide phone', (tester) async {
      await _openRoadmap(tester, width);

      // A RenderFlex overflow would already have failed the test; also check
      // that every label and node sits fully on screen.
      for (final s in _skills) {
        final title = find.text(s.title);
        expect(title, findsOneWidget, reason: s.title);
        final rect = tester.getRect(title);
        expect(rect.left, greaterThanOrEqualTo(0), reason: '${s.title} starts off-screen');
        expect(rect.right, lessThanOrEqualTo(width), reason: '${s.title} runs off-screen');
      }
      for (final icon in [Icons.lock_rounded, Icons.check_rounded]) {
        final icons = find.byIcon(icon);
        for (var i = 0; i < icons.evaluate().length; i++) {
          final rect = tester.getRect(icons.at(i));
          expect(rect.left >= 0 && rect.right <= width, isTrue, reason: '$icon off-screen at $rect');
        }
      }
    });
  }

  testWidgets('roadmap shows each node\'s status like the web', (tester) async {
    await _openRoadmap(tester, 360);
    expect(find.text('4/4 lessons'), findsOneWidget); // completed
    expect(find.text('1/4 lessons'), findsOneWidget); // current
    expect(find.text('Finish the previous skill'), findsNWidgets(3)); // locked
    expect(find.byIcon(Icons.check_rounded), findsOneWidget);
    expect(find.text('Fluency'), findsOneWidget);
  });
}
