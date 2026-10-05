import 'package:flutter_test/flutter_test.dart';
import 'package:lumora_mobile/core/grading.dart';
import 'package:lumora_mobile/core/languages.dart';
import 'package:lumora_mobile/core/voices.dart';

void main() {
  test('Hindi is an available course, written with spaces (counted in words)', () {
    expect(kAvailableLanguages.map((l) => l.code), contains('hi'));
    expect(languageMeta('hi').nativeName, 'हिन्दी');
    expect(countsCharacters('hi'), isFalse);
    expect(countWords('मेरा नाम अन्ना है। मैं केन्या से हूँ।'), 8);
  });

  test('Devanagari vowel signs are never ignored when grading', () {
    // किताब vs कताब differ only in a vowel sign — that's a wrong answer.
    expect(gradeTyped('कताब', 'किताब').verdict, isNot(Verdict.correct));
    // The danda is punctuation.
    expect(gradeTyped('मैं ठीक हूँ', 'मैं ठीक हूँ।').verdict, Verdict.correct);
  });

  test('the copy check works on Hindi', () {
    const ex = 'मैं नैरोबी में रहती हूँ। सुबह थोड़ी सर्दी होती है, लेकिन दोपहर में धूप निकलती है।';
    expect(copiesExample(ex, ex), isTrue);
    expect(copiesExample('मेरा घर मुंबई में है और वहाँ हमेशा गर्मी रहती है।', ex), isFalse);
  });
}
