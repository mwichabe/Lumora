import 'package:flutter_test/flutter_test.dart';
import 'package:lumora_mobile/core/languages.dart';
import 'package:lumora_mobile/core/voices.dart';

void main() {
  group('Japanese text', () {
    test('counts each kana and kanji (Japanese has no spaces)', () {
      expect(countWords('きょうはねつがあります。'), 11);
      expect(countWords('会議を変更してください'), 11);
      expect(countWords('コーヒーをください'), 9);
    });

    test('scores a Japanese phrase by character, not as one word', () {
      expect(scorePronunciation('ありがとうございます', 'ありがとうございます'), 100);
      // Close but not identical still earns partial credit.
      final partial = scorePronunciation('ありがとうございます', 'ありがとう');
      expect(partial, greaterThan(20));
      expect(partial, lessThan(100));
    });
  });

  group('JLPT levels', () {
    test('show the JLPT level next to CEFR', () {
      expect(levelDisplay('A1', 'ja'), 'JLPT N5 · A1');
      expect(levelDisplay('C1', 'ja'), 'JLPT N1 · C1');
      expect(levelDisplay('C2', 'ja'), 'Beyond N1 · C2');
      expect(levelDisplay('FINAL', 'ja'), 'JLPT N1+ (Advanced)');
      expect(hasNativeLevels('ja'), isTrue);
      expect(countsCharacters('ja'), isTrue);
      expect(hasNativeLevels('de'), isFalse);
    });

    test('Japanese is an available course', () {
      expect(kAvailableLanguages.map((l) => l.code), contains('ja'));
      expect(languageMeta('ja').nativeName, '日本語');
    });
  });
}
