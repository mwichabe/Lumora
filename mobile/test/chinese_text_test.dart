import 'package:flutter_test/flutter_test.dart';

import 'package:lumora_mobile/core/languages.dart';
import 'package:lumora_mobile/core/voices.dart';

void main() {
  group('countWords', () {
    test('counts each Chinese character (Chinese has no spaces)', () {
      expect(countWords('我叫小明。'), 4);
      expect(countWords('你好，我是肯尼亚人！'), 8);
    });
    test('still counts words in space-separated languages', () {
      expect(countWords('Hola, ¿cómo estás?'), 3);
      expect(countWords('  '), 0);
    });
    test('handles mixed text', () {
      expect(countWords('I like 饺子'), 4);
    });
  });

  group('scorePronunciation', () {
    test('ignores Chinese punctuation', () {
      expect(scorePronunciation('你好，我叫小明。', '你好我叫小明'), 100);
    });
    test('gives partial credit per character for a near miss', () {
      // One character of four wrong. Scored as one whole "word" this used to
      // be ~30 — a fail — however close the learner was.
      final score = scorePronunciation('我叫小明', '我叫小红');
      expect(score, greaterThanOrEqualTo(70));
      expect(score, lessThan(100));
    });
    test('is unchanged for Spanish', () {
      expect(scorePronunciation('Buenos días', 'buenos dias'), 100);
    });
  });

  group('levelDisplay', () {
    test('shows the HSK level for Mandarin', () {
      expect(levelDisplay('A1', 'zh'), 'HSK 1 · A1');
      expect(levelDisplay('C2', 'zh'), 'HSK 6 · C2');
      expect(levelDisplay('FINAL', 'zh'), 'HSK 7–9 (Advanced)');
    });
    test('keeps CEFR for other languages', () {
      expect(levelDisplay('B2', 'de'), 'B2');
      expect(levelDisplay('FINAL', 'es'), 'Final Mastery');
    });
  });

  test('Mandarin is offered as a course', () {
    expect(kAvailableLanguages.map((l) => l.code), contains('zh'));
    expect(languageMeta('zh').name, 'Mandarin');
  });
}
