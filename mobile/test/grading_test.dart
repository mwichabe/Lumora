import 'package:flutter_test/flutter_test.dart';
import 'package:lumora_mobile/core/grading.dart';

// Same cases as the web grader (frontend/lib/grading.ts).
void main() {
  test('verdicts', () {
    const cases = [
      ('Ich bin glücklich', 'Ich bin glücklich', Verdict.correct),
      ('ich bin glücklich!', 'Ich bin glücklich.', Verdict.correct),
      ('Ich bin glucklich', 'Ich bin glücklich', Verdict.almost),
      ('Hemos vivido aqui', 'Hemos vivido aquí', Verdict.almost),
      ('Ich bin glücklih', 'Ich bin glücklich', Verdict.almost),
      ('Du bist mude', 'Du bist müde', Verdict.almost),
      ('bin ich glücklich', 'Ich bin glücklich', Verdict.incorrect),
      ('Ich bist glücklich', 'Ich bin glücklich', Verdict.incorrect),
      ('Ich glücklich', 'Ich bin glücklich', Verdict.incorrect),
      ('Ich bin sehr glücklich', 'Ich bin glücklich', Verdict.incorrect),
      ('', 'Ich bin glücklich', Verdict.incorrect),
      ('strasse', 'Straße', Verdict.almost),
      ('bin', 'bin', Verdict.correct),
      ('bist', 'bin', Verdict.incorrect),
    ];
    for (final (a, e, want) in cases) {
      expect(gradeTyped(a, e).verdict, want, reason: '"$a" vs "$e"');
    }
  });

  test('points out the error', () {
    final wrong = gradeTyped('Ich bist glücklich', 'Ich bin glücklich').issues;
    expect(wrong.single.kind, IssueKind.wrong);
    expect(wrong.single.yours, 'bist');
    expect(wrong.single.correct, 'bin');
    final missing = gradeTyped('Ich glücklich', 'Ich bin glücklich').issues;
    expect(missing.single.kind, IssueKind.missing);
    expect(missing.single.correct, 'bin');
  });

  test('copied example is detected', () {
    const ex = 'Sehr geehrte Damen und Herren, ich komme am Freitag an. Haben Sie ein Zimmer frei?';
    expect(copiesExample(ex, ex), isTrue);
    expect(copiesExample('ich komme am Freitag an. Haben Sie ein Zimmer frei?', ex), isTrue);
    expect(copiesExample('Hallo, mein Name ist Peter und ich brauche ein Einzelzimmer.', ex), isFalse);
  });
}
