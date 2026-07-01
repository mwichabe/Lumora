"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { X, Volume2, Mic, Check } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { MistakesReview, ReviewItem } from "@/components/MistakesReview";
import { SpeakerChip } from "@/components/Speaker";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import {
  speakAs,
  speakSequence,
  stopSpeaking,
  recognizeSpeech,
  scorePronunciation,
  speechRecognitionSupported,
} from "@/lib/voices";
import type {
  VocabItem,
  Mistake,
  ListeningSession,
  ReadingSession,
} from "@/lib/types";

type Mode =
  | "mix"
  | "quiz"
  | "listen"
  | "speak"
  | "mistakes"
  | "listening"
  | "reading";

interface Drill {
  kind: "choose" | "speak";
  prompt: string;
  question: string;
  correct: string;
  options?: string[];
  listen?: boolean;
  speaker?: string;
  mistakeId?: number;
}

const MODE_TITLE: Record<Mode, string> = {
  mix: "Daily Mix",
  quiz: "Vocabulary Quiz",
  listen: "Listening Drill",
  speak: "Speaking Practice",
  mistakes: "Review Mistakes",
  listening: "Listening Comprehension",
  reading: "Reading Comprehension",
};

const shuffle = <T,>(a: T[]) => [...a].sort(() => Math.random() - 0.5);
const sample = <T,>(a: T[], n: number) => shuffle(a).slice(0, n);

/**
 * Picks `n` vocab items the learner hasn't practised recently, so sessions don't
 * repeat. Recently-seen words are remembered per language in localStorage and
 * the set is cleared once the whole pool has been cycled through.
 */
function pickFresh(vocab: VocabItem[], n: number, lang: string): VocabItem[] {
  const key = `lumora_practice_seen_${lang}`;
  let seen: string[] = [];
  try {
    seen = JSON.parse(localStorage.getItem(key) || "[]");
  } catch {
    seen = [];
  }
  const seenSet = new Set(seen);
  let fresh = vocab.filter((v) => !seenSet.has(v.word));
  // Not enough new words left → cycle: start the pool over.
  if (fresh.length < n) {
    seen = [];
    fresh = vocab;
  }
  const chosen = sample(fresh, Math.min(n, fresh.length));
  const updated = [...seen, ...chosen.map((c) => c.word)].slice(-Math.max(n, 200));
  try {
    localStorage.setItem(key, JSON.stringify(updated));
  } catch {
    /* ignore */
  }
  return chosen;
}

function buildDrills(
  mode: Mode,
  vocab: VocabItem[],
  mistakes: Mistake[],
  lang: string
): Drill[] {
  const translations = Array.from(new Set(vocab.map((v) => v.translation)));
  const words = Array.from(new Set(vocab.map((v) => v.word)));

  const options = (correct: string, pool: string[]) =>
    shuffle([correct, ...sample(pool.filter((x) => x && x !== correct), 3)]);

  const quizDrill = (v: VocabItem, listen: boolean): Drill => ({
    kind: "choose",
    listen,
    prompt: listen ? "Listen and choose the meaning" : "What does this mean?",
    question: v.word,
    correct: v.translation,
    options: options(v.translation, translations),
    speaker: v.speaker,
  });
  const speakDrill = (v: VocabItem): Drill => ({
    kind: "speak",
    prompt: "Say it out loud",
    question: v.example || v.word,
    correct: v.example || v.word,
    speaker: v.speaker,
  });

  if (mode === "mistakes") {
    const answers = Array.from(
      new Set([...translations, ...words, ...mistakes.map((m) => m.correctAnswer)])
    );
    return mistakes.slice(0, 12).map((m) => ({
      kind: "choose",
      prompt: m.prompt || "Choose the correct answer",
      question: m.question,
      correct: m.correctAnswer,
      options: options(m.correctAnswer, answers),
      mistakeId: m.id,
    }));
  }

  if (mode === "speak") {
    return pickFresh(vocab, 8, lang).map(speakDrill);
  }

  if (mode === "mix") {
    // A varied set of fresh words: quiz + listening + speaking.
    const picks = pickFresh(vocab, 10, lang);
    const drills: Drill[] = [
      ...picks.slice(0, 4).map((v) => quizDrill(v, false)),
      ...picks.slice(4, 7).map((v) => quizDrill(v, true)),
      ...picks.slice(7).map(speakDrill),
    ];
    return shuffle(drills);
  }

  // quiz / listen
  return pickFresh(vocab, 8, lang).map((v) => quizDrill(v, mode === "listen"));
}

export default function PracticeRunnerPage() {
  return (
    <Suspense fallback={null}>
      <RunnerDispatch />
    </Suspense>
  );
}

function RunnerDispatch() {
  const params = useSearchParams();
  const mode = (params.get("mode") as Mode) || "quiz";
  if (mode === "listening" || mode === "reading") {
    return <SessionRunner mode={mode} />;
  }
  return <DrillRunner mode={mode} />;
}

function DrillRunner({ mode }: { mode: Mode }) {
  const router = useRouter();
  const { user, setUser } = useAuth();
  const lang = user?.targetLanguage || "es";

  const [loading, setLoading] = useState(true);
  const [drills, setDrills] = useState<Drill[]>([]);
  const [idx, setIdx] = useState(0);
  const [correct, setCorrect] = useState(0);
  const [resolved, setResolved] = useState<number[]>([]);
  const [misses, setMisses] = useState<ReviewItem[]>([]);
  const [reviewing, setReviewing] = useState(false);
  const [done, setDone] = useState(false);
  const [xp, setXp] = useState(0);

  useEffect(() => {
    api
      .practice()
      .then((d) => setDrills(buildDrills(mode, d.vocab, d.mistakes, lang)))
      .catch(() => setDrills([]))
      .finally(() => setLoading(false));
    return () => stopSpeaking();
  }, [mode, lang]);

  const finish = useCallback(
    async (finalCorrect: number, resolvedIds: number[]) => {
      const earned = Math.min(50, 5 + finalCorrect * 3);
      setXp(earned);
      setReviewing(false);
      setDone(true);
      try {
        const r = await api.completePractice(earned);
        setUser(r.user);
      } catch {
        /* ignore */
      }
      if (resolvedIds.length) {
        api.resolveMistakes(resolvedIds).catch(() => {});
      }
    },
    [setUser]
  );

  function next(wasCorrect: boolean, mistakeId?: number) {
    const d = drills[idx];
    const nc = correct + (wasCorrect ? 1 : 0);
    const nr = wasCorrect && mistakeId ? [...resolved, mistakeId] : resolved;

    // On a miss: record it for Practice → Review Mistakes (skip drills that are
    // already mistakes), and collect it for this session's review recap.
    let nm = misses;
    if (!wasCorrect && d) {
      if (d.kind === "choose" && !d.mistakeId) {
        api
          .recordMistake({
            prompt: d.prompt,
            question: d.question,
            correctAnswer: d.correct,
          })
          .catch(() => {});
      }
      nm = [
        ...misses,
        {
          prompt: d.prompt,
          question: d.question,
          correctAnswer: d.correct,
          playText: d.mistakeId ? undefined : d.question,
          speaker: d.speaker,
        },
      ];
    }

    setCorrect(nc);
    setResolved(nr);
    setMisses(nm);

    if (idx + 1 < drills.length) {
      setIdx((i) => i + 1);
    } else if (nm.length > 0) {
      setFinalResult({ nc, nr });
      setReviewing(true); // study misses before the completion screen
    } else {
      finish(nc, nr);
    }
  }

  const [finalResult, setFinalResult] = useState<{
    nc: number;
    nr: number[];
  } | null>(null);

  if (loading) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={() => router.push("/practice")}>
        <div className="flex flex-1 items-center justify-center">
          <FoxMascot size={110} glow />
        </div>
      </Shell>
    );
  }

  if (reviewing && finalResult) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={() => router.push("/practice")}>
        <MistakesReview
          items={misses}
          finishLabel="Finish practice"
          onDone={() => finish(finalResult.nc, finalResult.nr)}
        />
      </Shell>
    );
  }

  if (done) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={() => router.push("/practice")}>
        <div className="flex flex-1 flex-col items-center justify-center text-center">
          <FoxMascot size={130} glow bounce />
          <h2 className="mt-4 text-heading-xl font-extrabold text-ink">
            Practice complete!
          </h2>
          <p className="mt-1 text-body-md text-slatey">
            You got {correct}/{drills.length} right.
          </p>
          <span className="mt-5 flex h-9 items-center gap-1.5 rounded-full bg-amber px-4 font-extrabold text-ink">
            <Check size={16} /> +{xp} XP
          </span>
          <div className="mt-8 w-full max-w-sm">
            <Button full onClick={() => router.push("/practice")}>
              Back to practice
            </Button>
          </div>
        </div>
      </Shell>
    );
  }

  if (drills.length === 0) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={() => router.push("/practice")}>
        <div className="flex flex-1 flex-col items-center justify-center text-center">
          <FoxMascot size={110} glow />
          <p className="mt-4 text-heading-sm font-extrabold text-ink">
            {mode === "mistakes" ? "No mistakes to review 🎉" : "Nothing to practise yet"}
          </p>
          <p className="mt-1 max-w-xs text-body-md text-slatey">
            {mode === "mistakes"
              ? "Answer some lessons — anything you miss shows up here."
              : "Start a lesson to unlock practice for this language."}
          </p>
          <div className="mt-6 w-full max-w-sm">
            <Button full variant="outline" onClick={() => router.push("/learn")}>
              Go to lessons
            </Button>
          </div>
        </div>
      </Shell>
    );
  }

  const drill = drills[idx];
  const progress = (idx / drills.length) * 100;

  return (
    <Shell
      title={MODE_TITLE[mode]}
      onClose={() => router.push("/practice")}
      progress={progress}
    >
      {drill.kind === "speak" ? (
        <SpeakDrill key={idx} drill={drill} onNext={next} />
      ) : (
        <ChooseDrill key={idx} drill={drill} onNext={next} />
      )}
    </Shell>
  );
}

function Shell({
  title,
  onClose,
  progress,
  children,
}: {
  title: string;
  onClose: () => void;
  progress?: number;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-[100dvh] w-full justify-center bg-cream lg:bg-[#eceaf3]">
      <div className="flex min-h-[100dvh] w-full max-w-2xl flex-col bg-cream lg:my-8 lg:min-h-[calc(100dvh-4rem)] lg:overflow-hidden lg:rounded-3xl lg:shadow-card-lg">
        <header className="flex items-center gap-3 px-4 pb-2 pt-12 lg:px-8 lg:pt-6">
          <button onClick={onClose} aria-label="Close practice">
            <X className="text-gray-500" />
          </button>
          {progress != null ? (
            <div className="h-2 flex-1 overflow-hidden rounded-full bg-gray-100">
              <motion.div
                className="h-full rounded-full bg-purple"
                animate={{ width: `${progress}%` }}
                transition={{ duration: 0.3 }}
              />
            </div>
          ) : (
            <h1 className="flex-1 text-heading-sm font-extrabold text-ink">{title}</h1>
          )}
        </header>
        <div className="flex flex-1 flex-col px-5 pb-8 pt-3 lg:px-8">{children}</div>
      </div>
    </div>
  );
}

function ChooseDrill({
  drill,
  onNext,
}: {
  drill: Drill;
  onNext: (correct: boolean, mistakeId?: number) => void;
}) {
  const [answer, setAnswer] = useState("");
  const [feedback, setFeedback] = useState<null | "correct" | "incorrect">(null);

  useEffect(() => {
    if (drill.listen) speakAs(drill.speaker || "Mira", drill.question);
  }, [drill]);

  function check() {
    if (!answer) return;
    setFeedback(answer === drill.correct ? "correct" : "incorrect");
  }

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-label-sm font-bold uppercase tracking-wide text-gray-500">
        {drill.prompt}
      </p>

      {drill.listen ? (
        <div className="mt-5 flex flex-col items-center gap-3 py-4">
          <SpeakerChip name={drill.speaker || "Mira"} />
          <button
            onClick={() => speakAs(drill.speaker || "Mira", drill.question)}
            className="flex h-16 w-16 items-center justify-center rounded-full bg-purple text-white shadow-float"
            aria-label="Play audio"
          >
            <Volume2 size={28} />
          </button>
          <p className="text-body-sm text-slatey">Tap to listen again</p>
        </div>
      ) : (
        <p className="mt-3 text-display-lg font-extrabold text-ink">{drill.question}</p>
      )}

      <div className="mt-5 space-y-2">
        {(drill.options || []).map((opt) => {
          const selected = answer === opt;
          const isCorrect = opt === drill.correct;
          let cls = "border-gray-100 bg-white";
          if (feedback && isCorrect) cls = "border-teal bg-teal-light";
          else if (feedback && selected && !isCorrect) cls = "border-coral bg-coral-light";
          else if (selected) cls = "border-purple bg-purple-light";
          return (
            <button
              key={opt}
              onClick={() => !feedback && setAnswer(opt)}
              className={`flex h-14 w-full items-center rounded-md border-2 px-4 text-left text-body-lg font-semibold transition ${cls}`}
            >
              {opt}
            </button>
          );
        })}
      </div>

      <div className="mt-auto pt-6">
        <AnimatePresence mode="wait">
          {feedback ? (
            <motion.div
              key="fb"
              initial={{ y: 16, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              className={`rounded-xl p-4 ${
                feedback === "correct" ? "bg-teal-light" : "bg-coral-light"
              }`}
            >
              <p
                className={`mb-3 font-extrabold ${
                  feedback === "correct" ? "text-teal" : "text-coral"
                }`}
              >
                {feedback === "correct" ? "Correct!" : `Answer: ${drill.correct}`}
              </p>
              <Button
                full
                variant={feedback === "correct" ? "primary" : "danger"}
                onClick={() => onNext(feedback === "correct", drill.mistakeId)}
              >
                Continue
              </Button>
            </motion.div>
          ) : (
            <Button key="check" full disabled={!answer} onClick={check}>
              Check
            </Button>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

function SpeakDrill({
  drill,
  onNext,
}: {
  drill: Drill;
  onNext: (correct: boolean) => void;
}) {
  const [listening, setListening] = useState(false);
  const [score, setScore] = useState<number | null>(null);
  const [heard, setHeard] = useState<string | null>(null);
  const [supported] = useState(() => speechRecognitionSupported());

  async function record() {
    if (listening) return;
    setScore(null);
    setHeard(null);
    setListening(true);
    try {
      const said = await recognizeSpeech();
      setHeard(said);
      setScore(scorePronunciation(drill.correct, said));
    } catch {
      setScore(null);
    } finally {
      setListening(false);
    }
  }

  const tone = score == null ? "" : score >= 80 ? "text-teal" : score >= 50 ? "text-amber" : "text-coral";

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-label-sm font-bold uppercase tracking-wide text-coral">
        Speaking practice
      </p>

      <div className="mt-6 flex flex-1 flex-col items-center justify-center text-center">
        <SpeakerChip name={drill.speaker || "Lumora"} className="mb-4" />
        <p className="text-heading-lg font-extrabold text-purple">
          &ldquo;{drill.question}&rdquo;
        </p>
        <div className="mt-5 flex items-center gap-3">
          <button
            onClick={() => speakAs(drill.speaker || "Lumora", drill.question)}
            className="flex h-12 w-12 items-center justify-center rounded-full bg-purple text-white shadow-float"
            aria-label="Hear it"
          >
            <Volume2 size={22} />
          </button>
          <button
            onClick={record}
            disabled={listening || !supported}
            className={`flex h-16 w-16 items-center justify-center rounded-full bg-coral text-white shadow-float disabled:opacity-40 ${
              listening ? "animate-pulse" : ""
            }`}
            aria-label="Tap to speak"
          >
            <Mic size={28} />
          </button>
        </div>
        <p className="mt-3 text-body-sm text-slatey">
          {!supported
            ? "Speech scoring needs Chrome/Edge — tap Skip."
            : listening
            ? "Listening… say the phrase"
            : "Tap the mic and say it out loud"}
        </p>

        {score != null && (
          <div className="mt-5 w-full max-w-xs rounded-xl bg-gray-50 p-3">
            <p className={`text-display-lg font-extrabold ${tone}`}>{score}%</p>
            {heard && <p className="text-body-sm text-slatey">You said: “{heard}”</p>}
          </div>
        )}
      </div>

      <div className="mt-auto pt-6">
        {score == null ? (
          <Button full variant="outline" onClick={() => onNext(false)}>
            Skip
          </Button>
        ) : (
          <Button full onClick={() => onNext(score >= 60)}>
            Continue
          </Button>
        )}
      </div>
    </div>
  );
}

/* ---------- listening / reading comprehension ---------- */

function SessionRunner({ mode }: { mode: "listening" | "reading" }) {
  const router = useRouter();
  const { user, setUser } = useAuth();
  const lang = user?.targetLanguage || "es";

  const [loading, setLoading] = useState(true);
  const [sessions, setSessions] = useState<(ListeningSession | ReadingSession)[]>([]);
  const [current, setCurrent] = useState<ListeningSession | ReadingSession | null>(
    null
  );
  const [answers, setAnswers] = useState<Record<number, string>>({});
  const [submitted, setSubmitted] = useState(false);
  const [reviewing, setReviewing] = useState(false);
  const [reviewItems, setReviewItems] = useState<ReviewItem[]>([]);
  const [score, setScore] = useState(0);
  const [xp, setXp] = useState(0);
  const [played, setPlayed] = useState(false);

  // Rotate through unlocked sessions so practice never feels static — avoid
  // repeating the one just done (remembered per mode + language).
  const pick = useCallback(
    (list: (ListeningSession | ReadingSession)[]) => {
      if (!list.length) return null;
      const key = `lumora_practice_sess_${mode}_${lang}`;
      let last = 0;
      try {
        last = Number(localStorage.getItem(key) || 0);
      } catch {
        last = 0;
      }
      const pool = list.length > 1 ? list.filter((s) => s.id !== last) : list;
      const chosen = pool[Math.floor(Math.random() * pool.length)];
      try {
        localStorage.setItem(key, String(chosen.id));
      } catch {
        /* ignore */
      }
      return chosen;
    },
    [mode, lang]
  );

  useEffect(() => {
    const fetcher =
      mode === "listening" ? api.practiceListening() : api.practiceReading();
    fetcher
      .then((d) => {
        const list = d.sessions || [];
        setSessions(list);
        setCurrent(pick(list));
      })
      .catch(() => setSessions([]))
      .finally(() => setLoading(false));
    return () => stopSpeaking();
  }, [mode, pick]);

  function practiceAnother() {
    setAnswers({});
    setSubmitted(false);
    setReviewing(false);
    setReviewItems([]);
    setScore(0);
    setXp(0);
    setPlayed(false);
    stopSpeaking();
    setCurrent(pick(sessions));
  }

  const questions = current?.questions || [];
  const allAnswered =
    questions.length > 0 && questions.every((_, i) => answers[i]);

  async function submit() {
    const correct = questions.filter(
      (q, i) => answers[i] === q.correctAnswer
    ).length;
    const pct = questions.length
      ? Math.round((correct / questions.length) * 100)
      : 0;
    setScore(pct);
    setSubmitted(true);
    stopSpeaking();
    // wrong answers feed the Review Mistakes pile + this session's recap
    const wrong: ReviewItem[] = [];
    questions.forEach((q, i) => {
      if (answers[i] && answers[i] !== q.correctAnswer) {
        api
          .recordMistake({
            prompt: q.prompt || MODE_TITLE[mode],
            question: q.question,
            correctAnswer: q.correctAnswer,
          })
          .catch(() => {});
        wrong.push({
          prompt: q.prompt || MODE_TITLE[mode],
          question: q.question,
          correctAnswer: q.correctAnswer,
        });
      }
    });
    if (wrong.length) {
      setReviewItems(wrong);
      setReviewing(true); // study misses before the score screen
    }
    const earned = Math.min(45, 8 + correct * 4);
    setXp(earned);
    try {
      const r = await api.completePractice(earned);
      setUser(r.user);
    } catch {
      /* ignore */
    }
  }

  const close = () => {
    stopSpeaking();
    router.push("/practice");
  };

  if (loading) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={close}>
        <div className="flex flex-1 items-center justify-center">
          <FoxMascot size={110} glow />
        </div>
      </Shell>
    );
  }

  if (!current) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={close}>
        <div className="flex flex-1 flex-col items-center justify-center text-center">
          <FoxMascot size={110} glow />
          <p className="mt-4 text-heading-sm font-extrabold text-ink">
            Nothing unlocked yet
          </p>
          <p className="mt-1 max-w-xs text-body-md text-slatey">
            Complete more lessons to unlock{" "}
            {mode === "listening" ? "listening conversations" : "reading passages"}{" "}
            for practice.
          </p>
          <div className="mt-6 w-full max-w-sm">
            <Button full variant="outline" onClick={() => router.push("/learn")}>
              Go to lessons
            </Button>
          </div>
        </div>
      </Shell>
    );
  }

  if (reviewing) {
    return (
      <Shell title={MODE_TITLE[mode]} onClose={close}>
        <MistakesReview
          items={reviewItems}
          finishLabel="See my score"
          onDone={() => setReviewing(false)}
        />
      </Shell>
    );
  }

  if (submitted) {
    const tone = score >= 80 ? "text-teal" : score >= 50 ? "text-amber" : "text-coral";
    return (
      <Shell title={MODE_TITLE[mode]} onClose={close}>
        <div className="flex flex-1 flex-col items-center justify-center text-center">
          <FoxMascot size={130} glow bounce />
          <h2 className="mt-4 text-heading-xl font-extrabold text-ink">
            {score >= 60 ? "Nicely done!" : "Keep going!"}
          </h2>
          <p className={`mt-2 text-display-lg font-extrabold ${tone}`}>{score}%</p>
          <p className="mt-1 text-body-md text-slatey">on “{current.title}”</p>
          <span className="mt-4 flex h-9 items-center gap-1.5 rounded-full bg-amber px-4 font-extrabold text-ink">
            <Check size={16} /> +{xp} XP
          </span>
          <div className="mt-8 w-full max-w-sm space-y-2">
            <Button full onClick={practiceAnother}>
              Practice another
            </Button>
            <Button full variant="outline" onClick={close}>
              Back to practice
            </Button>
          </div>
        </div>
      </Shell>
    );
  }

  const listeningLines =
    mode === "listening" ? (current as ListeningSession).lines || [] : [];
  const readingLines =
    mode === "reading" ? (current as ReadingSession).lines || [] : [];

  return (
    <Shell title={MODE_TITLE[mode]} onClose={close}>
      <div className="flex flex-1 flex-col">
        <p className="text-label-sm font-bold uppercase tracking-wide text-gray-500">
          {mode === "listening" ? "Listen and answer" : "Read and answer"}
        </p>
        <h2 className="mt-1 text-heading-md font-extrabold text-ink">
          {current.title}
        </h2>

        {mode === "listening" ? (
          <div className="mt-4">
            <div className="mb-3 flex flex-wrap gap-2">
              {Array.from(new Set(listeningLines.map((l) => l.character))).map(
                (speaker) => (
                  <SpeakerChip key={speaker} name={speaker} />
                )
              )}
            </div>
            <button
              onClick={() => {
                setPlayed(true);
                speakSequence(
                  listeningLines.map((l) => ({
                    character: l.character,
                    text: l.text,
                  }))
                );
              }}
              className="flex w-full items-center justify-center gap-2 rounded-full bg-purple py-3 font-extrabold text-white shadow-float"
            >
              <Volume2 size={20} />{" "}
              {played ? "Play again" : "Play the conversation"}
            </button>
            <p className="mt-2 text-center text-body-sm text-slatey">
              Replay as many times as you need.
            </p>
          </div>
        ) : (
          <div className="mt-4 space-y-2 rounded-2xl border border-gray-100 bg-white p-4 shadow-card">
            {readingLines.map((l) => (
              <p key={l.id} className="text-body-lg leading-relaxed text-ink/90">
                {l.text}
              </p>
            ))}
          </div>
        )}

        <div className="mt-5 space-y-4">
          {questions.map((q, qi) => (
            <div key={qi}>
              <p className="mb-2 text-body-md font-extrabold text-ink">
                {qi + 1}. {q.question}
              </p>
              <div className="space-y-2">
                {(q.options || []).map((opt) => {
                  const sel = answers[qi] === opt;
                  return (
                    <button
                      key={opt}
                      onClick={() => setAnswers((a) => ({ ...a, [qi]: opt }))}
                      className={`flex w-full items-center rounded-md border-2 px-4 py-3 text-left text-body-md font-semibold transition ${
                        sel
                          ? "border-purple bg-purple-light"
                          : "border-gray-100 bg-white"
                      }`}
                    >
                      {opt}
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>

        <div className="mt-6 pb-2">
          <Button full disabled={!allAnswered} onClick={submit}>
            Submit answers
          </Button>
        </div>
      </div>
    </Shell>
  );
}
