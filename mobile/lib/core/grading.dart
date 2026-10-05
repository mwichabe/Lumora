/// Grading for typed answers in lessons — a line-for-line port of
/// frontend/lib/grading.ts (keep the two in step; they share test cases).
///
///  * Capitals and punctuation never matter.
///  * A missing/wrong accent, or a one-letter slip in one word, still counts
///    as correct ("almost"), with a note showing the right spelling.
///  * Anything else is wrong, and `issues` says exactly what: which words were
///    wrong, missing, extra, or out of order.
library;

import 'dart:math';

enum Verdict { correct, almost, incorrect }

enum IssueKind { wrong, missing, extra, order, accent, typo }

class GradeIssue {
  final IssueKind kind;
  final String? yours;
  final String? correct;
  const GradeIssue(this.kind, {this.yours, this.correct});
}

class Grade {
  final Verdict verdict;
  final String note;
  final List<GradeIssue> issues;
  final String expected;
  const Grade(this.verdict, this.note, this.issues, this.expected);
}

final _punct = RegExp(r'[.,!?¡¿;:"«»„“”()…。，！？、；：]');

String _normalise(String s) => s
    .toLowerCase()
    .replaceAll(RegExp('[’`´]'), "'")
    .replaceAll(_punct, ' ')
    .replaceAll(RegExp(r'\s+'), ' ')
    .trim();

const _accents = {
  'á': 'a', 'à': 'a', 'â': 'a', 'ä': 'a', 'ã': 'a', 'å': 'a', 'ā': 'a', 'ǎ': 'a',
  'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e', 'ē': 'e', 'ě': 'e',
  'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i', 'ī': 'i', 'ǐ': 'i',
  'ó': 'o', 'ò': 'o', 'ô': 'o', 'ö': 'o', 'õ': 'o', 'ō': 'o', 'ǒ': 'o',
  'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u', 'ū': 'u', 'ǔ': 'u', 'ǖ': 'u', 'ǘ': 'u', 'ǚ': 'u', 'ǜ': 'u',
  'ñ': 'n', 'ç': 'c', 'ý': 'y', 'ÿ': 'y',
};

/// Accent-insensitive form (the TS version uses Unicode NFD; this table
/// covers every accented letter the courses use).
String _fold(String s) {
  final b = StringBuffer();
  for (final ch in s.split('')) {
    b.write(_accents[ch] ?? ch);
  }
  return b.toString().replaceAll('ß', 'ss').replaceAll('œ', 'oe').replaceAll('æ', 'ae');
}

List<String> _words(String s) => s.isEmpty ? [] : s.split(' ');

int _editDistance(String a, String b) {
  var prev = List<int>.generate(b.length + 1, (j) => j);
  for (var i = 1; i <= a.length; i++) {
    final cur = <int>[i];
    for (var j = 1; j <= b.length; j++) {
      cur.add([prev[j] + 1, cur[j - 1] + 1, prev[j - 1] + (a[i - 1] == b[j - 1] ? 0 : 1)].reduce(min));
    }
    prev = cur;
  }
  return prev[b.length];
}

/// Word-level diff (LCS on accent-folded words) → issues.
List<GradeIssue> _diff(List<String> yours, List<String> correct) {
  final fy = yours.map(_fold).toList();
  final fc = correct.map(_fold).toList();
  final m = fy.length, n = fc.length;
  final dp = List.generate(m + 1, (_) => List.filled(n + 1, 0));
  for (var i = m - 1; i >= 0; i--) {
    for (var j = n - 1; j >= 0; j--) {
      dp[i][j] = fy[i] == fc[j] ? dp[i + 1][j + 1] + 1 : max(dp[i + 1][j], dp[i][j + 1]);
    }
  }
  final issues = <GradeIssue>[];
  var extra = <String>[], missing = <String>[];
  void flush() {
    // Pair a removed word with an added one as "wrong word", left to right.
    while (extra.isNotEmpty && missing.isNotEmpty) {
      issues.add(GradeIssue(IssueKind.wrong, yours: extra.removeAt(0), correct: missing.removeAt(0)));
    }
    for (final w in extra) {
      issues.add(GradeIssue(IssueKind.extra, yours: w));
    }
    for (final w in missing) {
      issues.add(GradeIssue(IssueKind.missing, correct: w));
    }
    extra = [];
    missing = [];
  }

  var i = 0, j = 0;
  while (i < m || j < n) {
    if (i < m && j < n && fy[i] == fc[j]) {
      flush();
      if (yours[i] != correct[j]) issues.add(GradeIssue(IssueKind.accent, yours: yours[i], correct: correct[j]));
      i++;
      j++;
    } else if (j < n && (i >= m || dp[i][j + 1] >= dp[i + 1][j])) {
      missing.add(correct[j++]);
    } else {
      extra.add(yours[i++]);
    }
  }
  flush();
  return issues;
}

/// Grades a typed answer against the expected one.
Grade gradeTyped(String answer, String expected) {
  final a = _normalise(answer);
  final e = _normalise(expected);
  if (a.isEmpty) return Grade(Verdict.incorrect, 'Type your answer.', const [], expected);
  if (a == e) return Grade(Verdict.correct, '', const [], expected);

  final aw = _words(a), ew = _words(e);

  if (_fold(a) == _fold(e)) {
    final issues = _diff(aw, ew).where((x) => x.kind == IssueKind.accent).toList();
    return Grade(Verdict.almost, 'Watch the accents: $expected', issues, expected);
  }

  // One small slip in one word of a same-length answer.
  if (aw.length == ew.length) {
    final differing = [
      for (var k = 0; k < aw.length; k++)
        if (_fold(aw[k]) != _fold(ew[k])) (aw[k], ew[k]),
    ];
    if (differing.length == 1) {
      final (w, c) = differing.first;
      if (c.length >= 4 && _editDistance(_fold(w), _fold(c)) <= 1) {
        return Grade(Verdict.almost, 'Small typo: “$w” should be “$c”.',
            [GradeIssue(IssueKind.typo, yours: w, correct: c)], expected);
      }
    }
  }

  final sortedA = (aw.map(_fold).toList()..sort()).join(' ');
  final sortedE = (ew.map(_fold).toList()..sort()).join(' ');
  if (sortedA == sortedE) {
    return Grade(Verdict.incorrect, 'Right words, wrong order.', const [GradeIssue(IssueKind.order)], expected);
  }

  final issues = _diff(aw, ew).where((x) => x.kind != IssueKind.accent).toList();
  final note = issues.isEmpty ? 'Not quite.' : (issues.length == 1 ? 'One thing to fix:' : '${issues.length} things to fix:');
  return Grade(Verdict.incorrect, note, issues, expected);
}

/// One readable line per issue, for the feedback panel.
List<String> issueLines(List<GradeIssue> issues) => [
      for (final x in issues.take(5))
        switch (x.kind) {
          IssueKind.wrong => '“${x.yours}” should be “${x.correct}”',
          IssueKind.missing => 'Missing: “${x.correct}”',
          IssueKind.extra => 'Not needed: “${x.yours}”',
          IssueKind.accent => 'Accent: “${x.yours}” → “${x.correct}”',
          IssueKind.typo => 'Typo: “${x.yours}” → “${x.correct}”',
          IssueKind.order => 'Right words, wrong order',
        },
    ];

// Words, or single characters in Chinese/Japanese (no spaces to split on).
final _tok = RegExp(
    r"[\p{Script=Han}\p{Script=Hiragana}\p{Script=Katakana}ー]|[^\s\p{Script=Han}\p{Script=Hiragana}\p{Script=Katakana}ー\p{P}\p{S}]+",
    unicode: true);

/// Client-side copy check for free writing, matching the server's rule
/// (backend/controllers/lesson_writing.go `copiesText`).
bool copiesExample(String answer, String example) {
  List<String> toks(String t) => _tok
      .allMatches(t.toLowerCase())
      .map((m) => m.group(0)!.replaceAll(RegExp(r"^['’-]+|['’-]+$"), ''))
      .where((w) => w.isNotEmpty)
      .toList();
  final a = toks(answer);
  final s = toks(example);
  if (a.isEmpty || s.length < 3) return false;
  final aj = ' ${a.join(' ')} ', sj = ' ${s.join(' ')} ';
  if (aj.contains(sj)) return true;
  if (a.length >= 5 && sj.contains(aj)) return true;
  if (a.length < 4) return false;
  final grams = <String>{for (var i = 0; i + 3 <= s.length; i++) s.sublist(i, i + 3).join(' ')};
  var total = 0, hit = 0;
  for (var i = 0; i + 3 <= a.length; i++) {
    total++;
    if (grams.contains(a.sublist(i, i + 3).join(' '))) hit++;
  }
  return total > 0 && hit / total >= 0.6;
}
