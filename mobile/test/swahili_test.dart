import 'package:flutter_test/flutter_test.dart';
import 'package:lumora_mobile/core/grading.dart';
import 'package:lumora_mobile/core/languages.dart';

void main() {
  test('Swahili is an available course on the CEFR scale', () {
    expect(kAvailableLanguages.map((l) => l.code), contains('sw'));
    expect(languageMeta('sw').nativeName, 'Kiswahili');
    expect(levelDisplay('B1', 'sw'), 'B1');
    expect(hasNativeLevels('sw'), isFalse);
    expect(countsCharacters('sw'), isFalse); // written with spaces: count words
  });

  test("typed Swahili answers keep ng' and grade like other Latin-script languages", () {
    expect(gradeTyped("Ng'ombe anakula", "Ng'ombe anakula.").verdict, Verdict.correct);
    expect(gradeTyped('Ng’ombe anakula', "Ng'ombe anakula").verdict, Verdict.correct); // curly apostrophe
    expect(gradeTyped('Ninajifunza Kiswahil', 'Ninajifunza Kiswahili').verdict, Verdict.almost); // one-letter typo
    final wrong = gradeTyped('Nilisoma kitabu', 'Ninasoma kitabu');
    expect(wrong.verdict, Verdict.incorrect);
    expect(wrong.issues.single.yours, 'nilisoma');
  });
}
