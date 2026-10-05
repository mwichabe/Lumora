import 'package:flutter_test/flutter_test.dart';
import 'package:lumora_mobile/core/grading.dart';
import 'package:lumora_mobile/core/languages.dart';

void main() {
  test('Italian is an available course on the CEFR scale', () {
    expect(kAvailableLanguages.map((l) => l.code), contains('it'));
    expect(languageMeta('it').nativeName, 'Italiano');
    expect(levelDisplay('B2', 'it'), 'B2');
    expect(countsCharacters('it'), isFalse);
  });

  test('typed Italian answers: accents and apostrophes', () {
    // A missing accent still counts, with a note.
    expect(gradeTyped('Lei e medico', 'Lei è medico').verdict, Verdict.almost);
    // Elision with either apostrophe.
    expect(gradeTyped('Ho vent’anni', "Ho vent'anni").verdict, Verdict.correct);
    // The wrong auxiliary is pointed out.
    final wrong = gradeTyped('Ho andato a Roma', 'Sono andato a Roma');
    expect(wrong.verdict, Verdict.incorrect);
    expect(wrong.issues.single.yours, 'ho');
    expect(wrong.issues.single.correct, 'sono');
  });
}
