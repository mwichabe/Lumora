/**
 * Grading for typed answers in lessons. Mirrored line for line in
 * mobile/lib/core/grading.dart — keep the two in step (they share test cases).
 *
 * - Capitals and punctuation never matter.
 * - A missing/wrong accent, or a one-letter slip in one word, still counts as
 *   correct ("almost"), with a note showing the right spelling.
 * - Anything else is wrong, and `issues` says exactly what: which words were
 *   wrong, missing, extra, or out of order.
 */

export type Verdict = "correct" | "almost" | "incorrect";

export interface GradeIssue {
  kind: "wrong" | "missing" | "extra" | "order" | "accent" | "typo";
  yours?: string;
  correct?: string;
}

export interface Grade {
  verdict: Verdict;
  /** One line for the learner, when there's something to say. */
  note: string;
  issues: GradeIssue[];
  expected: string;
}

const PUNCT = /[.,!?¡¿;:"«»„“”()…。，！？、；：]/g;

function normalise(s: string): string {
  return s
    .normalize("NFC")
    .toLowerCase()
    .replace(/[’`´]/g, "'")
    .replace(PUNCT, " ")
    .replace(/\s+/g, " ")
    .trim();
}

function fold(s: string): string {
  return s
    .normalize("NFD")
    .replace(/\p{M}/gu, "")
    .replace(/ß/g, "ss")
    .replace(/œ/g, "oe")
    .replace(/æ/g, "ae");
}

const words = (s: string) => (s ? s.split(" ") : []);

function editDistance(a: string, b: string): number {
  const m = a.length;
  const n = b.length;
  let prev = Array.from({ length: n + 1 }, (_, j) => j);
  for (let i = 1; i <= m; i++) {
    const cur = [i];
    for (let j = 1; j <= n; j++) {
      cur[j] = Math.min(
        prev[j] + 1,
        cur[j - 1] + 1,
        prev[j - 1] + (a[i - 1] === b[j - 1] ? 0 : 1)
      );
    }
    prev = cur;
  }
  return prev[n];
}

/** Word-level diff (LCS on accent-folded words) → issues. */
function diff(yours: string[], correct: string[]): GradeIssue[] {
  const fy = yours.map(fold);
  const fc = correct.map(fold);
  const m = fy.length;
  const n = fc.length;
  const dp: number[][] = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0));
  for (let i = m - 1; i >= 0; i--) {
    for (let j = n - 1; j >= 0; j--) {
      dp[i][j] = fy[i] === fc[j] ? dp[i + 1][j + 1] + 1 : Math.max(dp[i + 1][j], dp[i][j + 1]);
    }
  }
  const issues: GradeIssue[] = [];
  let extra: string[] = [];
  let missing: string[] = [];
  const flush = () => {
    // Pair a removed word with an added one as "wrong word", left to right.
    while (extra.length && missing.length) {
      issues.push({ kind: "wrong", yours: extra.shift(), correct: missing.shift() });
    }
    extra.forEach((w) => issues.push({ kind: "extra", yours: w }));
    missing.forEach((w) => issues.push({ kind: "missing", correct: w }));
    extra = [];
    missing = [];
  };
  let i = 0;
  let j = 0;
  while (i < m || j < n) {
    if (i < m && j < n && fy[i] === fc[j]) {
      flush();
      if (yours[i] !== correct[j]) issues.push({ kind: "accent", yours: yours[i], correct: correct[j] });
      i++;
      j++;
    } else if (j < n && (i >= m || dp[i][j + 1] >= dp[i + 1][j])) {
      missing.push(correct[j++]);
    } else {
      extra.push(yours[i++]);
    }
  }
  flush();
  return issues;
}

function gradeAgainst(answer: string, expected: string): Grade {
  const a = normalise(answer);
  const e = normalise(expected);
  const base = { issues: [] as GradeIssue[], expected };
  if (!a) return { ...base, verdict: "incorrect", note: "Type your answer." };
  if (a === e) return { ...base, verdict: "correct", note: "" };

  const aw = words(a);
  const ew = words(e);

  if (fold(a) === fold(e)) {
    const issues = diff(aw, ew).filter((x) => x.kind === "accent");
    return { issues, expected, verdict: "almost", note: `Watch the accents: ${expected}` };
  }

  // One small slip in one word of a same-length answer.
  if (aw.length === ew.length) {
    const differing = aw
      .map((w, k) => ({ w, c: ew[k] }))
      .filter(({ w, c }) => fold(w) !== fold(c));
    if (differing.length === 1) {
      const { w, c } = differing[0];
      if (c.length >= 4 && editDistance(fold(w), fold(c)) <= 1) {
        return {
          issues: [{ kind: "typo", yours: w, correct: c }],
          expected,
          verdict: "almost",
          note: `Small typo: “${w}” should be “${c}”.`,
        };
      }
    }
  }

  const sortedA = aw.map(fold).sort().join(" ");
  const sortedE = ew.map(fold).sort().join(" ");
  if (sortedA === sortedE) {
    return {
      issues: [{ kind: "order" }],
      expected,
      verdict: "incorrect",
      note: "Right words, wrong order.",
    };
  }

  const issues = diff(aw, ew).filter((x) => x.kind !== "accent");
  return { issues, expected, verdict: "incorrect", note: describe(issues) };
}

/** Grades a typed answer against the expected one. */
export function gradeTyped(answer: string, expected: string): Grade {
  return gradeAgainst(answer, expected);
}

/** One readable line per issue, for the feedback panel. */
export function issueLines(issues: GradeIssue[]): string[] {
  return issues.slice(0, 5).map((x) => {
    switch (x.kind) {
      case "wrong":
        return `“${x.yours}” should be “${x.correct}”`;
      case "missing":
        return `Missing: “${x.correct}”`;
      case "extra":
        return `Not needed: “${x.yours}”`;
      case "accent":
        return `Accent: “${x.yours}” → “${x.correct}”`;
      case "typo":
        return `Typo: “${x.yours}” → “${x.correct}”`;
      case "order":
        return "Right words, wrong order";
    }
  });
}

function describe(issues: GradeIssue[]): string {
  if (!issues.length) return "Not quite.";
  const n = issues.length;
  return n === 1 ? "One thing to fix:" : `${n} things to fix:`;
}

/**
 * Client-side copy check for free writing, matching the server's rule
 * (backend/controllers/lesson_writing.go `copiesText`) so the learner is told
 * straight away; the server still checks.
 */
export function copiesExample(answer: string, example: string): boolean {
  // Words, or single characters in Chinese/Japanese (no spaces to split on).
  const toks = (s: string) =>
    s
      .toLowerCase()
      .match(/[\p{Script=Han}\p{Script=Hiragana}\p{Script=Katakana}ー]|[^\s\p{Script=Han}\p{Script=Hiragana}\p{Script=Katakana}ー\p{P}\p{S}]+/gu)
      ?.map((t) => t.replace(/^['’-]+|['’-]+$/g, ""))
      .filter(Boolean) ?? [];
  const a = toks(answer);
  const s = toks(example);
  if (!a.length || s.length < 3) return false;
  const aj = ` ${a.join(" ")} `;
  const sj = ` ${s.join(" ")} `;
  if (aj.includes(sj)) return true;
  if (a.length >= 5 && sj.includes(aj)) return true;
  if (a.length < 4) return false;
  const grams = new Set<string>();
  for (let i = 0; i + 3 <= s.length; i++) grams.add(s.slice(i, i + 3).join(" "));
  let total = 0;
  let hit = 0;
  for (let i = 0; i + 3 <= a.length; i++) {
    total++;
    if (grams.has(a.slice(i, i + 3).join(" "))) hit++;
  }
  return total > 0 && hit / total >= 0.6;
}
