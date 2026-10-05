"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { X, Volume2, Check, Mic } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { SpeechBubble, HeartIndicator } from "@/components/widgets";
import { Button } from "@/components/Button";
import { MistakesReview, ReviewItem } from "@/components/MistakesReview";
import { SpeakerAvatar, SpeakerChip } from "@/components/Speaker";
import { OutOfHeartsModal } from "@/components/OutOfHeartsModal";
import { characterInfo } from "@/lib/characters";
import { useHearts } from "@/lib/hearts";
import { api } from "@/lib/api";
import {
  speakAs,
  stopSpeaking,
  recognizeSpeech,
  scorePronunciation,
  countWords,
  speechRecognitionSupported,
} from "@/lib/voices";
import { gradeTyped, issueLines, copiesExample } from "@/lib/grading";
import type { Lesson, Exercise, WritingCorrection } from "@/lib/types";

type Feedback = null | "correct" | "incorrect";

/** What the feedback bar explains beyond right/wrong. */
interface FeedbackDetail {
  /** Headline under "Nailed it!" / "Not quite" (e.g. "Watch the accents"). */
  note?: string;
  /** Pointed-out errors in a typed answer, one per line. */
  lines?: string[];
  /** Writing-coach corrections on free writing. */
  corrections?: WritingCorrection[];
  /** Hide "Answer: …" (free writing has no single answer). */
  hideAnswer?: boolean;
}

export default function LessonPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();

  const [lesson, setLesson] = useState<Lesson | null>(null);
  const [phase, setPhase] = useState<"vocab" | "practice" | "review">("vocab");
  const [idx, setIdx] = useState(0);
  const [misses, setMisses] = useState<ReviewItem[]>([]);
  const { hearts, secondsToNext, status: heartsStatus, lose, buy } = useHearts();
  const [outOfHearts, setOutOfHearts] = useState(false);
  const [buyingHearts, setBuyingHearts] = useState(false);
  const [answer, setAnswer] = useState("");
  const [feedback, setFeedback] = useState<Feedback>(null);
  const [correctCount, setCorrectCount] = useState(0);
  const [gradedCount, setGradedCount] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [detail, setDetail] = useState<FeedbackDetail>({});
  // Problems with a free-writing answer that must be fixed before it counts.
  const [writeProblems, setWriteProblems] = useState<string[]>([]);
  const [checking, setChecking] = useState(false);

  useEffect(() => {
    api
      .lesson(id)
      .then((d) => {
        setLesson(d.lesson);
        // Skip straight to practice if the lesson has no vocabulary.
        setPhase(d.lesson.vocab && d.lesson.vocab.length ? "vocab" : "practice");
      })
      .catch(() => {});
  }, [id]);

  // Starting (or continuing) with no hearts isn't allowed — show the gate.
  useEffect(() => {
    if (heartsStatus && heartsStatus.hearts <= 0) setOutOfHearts(true);
  }, [heartsStatus]);

  // Stop any speech when leaving the lesson.
  useEffect(() => () => stopSpeaking(), []);

  const exercises = lesson?.exercises || [];
  const ex: Exercise | undefined = exercises[idx];
  const total = exercises.length;
  const progress = total ? (idx / total) * 100 : 0;

  const isNarrative = ex?.type === "character";
  const isSpeak = ex?.type === "speak";
  const isWrite = ex?.type === "write";
  // Translate/fill now arrive with generated options, so render them as choices
  // too. Only fall back to typing if no options were provided.
  const hasOptions = !!(ex?.options && ex.options.length > 0);
  const needsChoice =
    !!ex &&
    !isNarrative &&
    !isSpeak &&
    !isWrite &&
    (["multiple_choice", "listen", "match"].includes(ex.type) || hasOptions);
  const needsTyping =
    !!ex && ["translate", "fill"].includes(ex.type) && !hasOptions;

  const wordCount = countWords(answer);

  const canCheck = useMemo(() => {
    if (!ex || feedback) return false;
    if (isNarrative || isSpeak) return true;
    if (checking) return false;
    return answer.trim().length > 0;
  }, [ex, feedback, answer, isNarrative, isSpeak, checking]);

  function normalise(s: string) {
    return s.trim().toLowerCase().replace(/[.,!¡¿?]/g, "");
  }

  function markCorrect(d: FeedbackDetail = {}) {
    setGradedCount((c) => c + 1);
    setCorrectCount((c) => c + 1);
    setDetail(d);
    setFeedback("correct");
  }

  async function check() {
    if (!ex) return;
    if (isNarrative || isSpeak) {
      advance();
      return;
    }
    if (isWrite) return checkWriting(ex);
    if (needsTyping) return checkTyped(ex);
    if (normalise(answer) === normalise(ex.correctAnswer)) markCorrect();
    else markWrong(ex);
  }

  /** Typed answers: point out exactly what's wrong, accept accent slips and
   *  one-letter typos, and ask the server about valid alternative wordings. */
  async function checkTyped(ex: Exercise) {
    const g = gradeTyped(answer, ex.correctAnswer);
    if (g.verdict !== "incorrect") {
      return markCorrect({ note: g.note || undefined, lines: issueLines(g.issues) });
    }
    if (ex.correctAnswer.trim().split(/\s+/).length >= 2) {
      setChecking(true);
      const second = await Promise.race([
        api
          .checkAnswer({
            lessonId: ex.lessonId,
            question: ex.question,
            expected: ex.correctAnswer,
            answer,
          })
          .catch(() => null),
        new Promise<null>((r) => setTimeout(() => r(null), 9000)),
      ]);
      setChecking(false);
      if (second?.available && second.correct) {
        return markCorrect({
          note: second.explanation || "That's another correct way to say it.",
          lines: [`Also correct: “${ex.correctAnswer}”`],
        });
      }
    }
    markWrong(ex, { note: g.note, lines: issueLines(g.issues) });
  }

  /** Free writing: rule checks (no copying the example, long enough, right
   *  language) must pass before it counts; then corrections are shown. */
  async function checkWriting(ex: Exercise) {
    if (copiesExample(answer, ex.correctAnswer)) {
      setWriteProblems([
        "An example cannot be used here — write your own answer in your own words.",
      ]);
      return;
    }
    setWriteProblems([]);
    setChecking(true);
    try {
      const r = await api.checkWriting(ex.id, answer);
      if (!r.ok) {
        setWriteProblems(r.problems.map((p) => p.message));
        return;
      }
      const d: FeedbackDetail = {
        note: r.summary,
        corrections: r.corrections ?? [],
        hideAnswer: true,
      };
      if (r.acceptable === false) markWrong(ex, d);
      else markCorrect(d);
    } catch {
      // Offline or server trouble: the local checks passed, so accept it.
      markCorrect({ note: "Saved. Detailed feedback isn't available right now.", hideAnswer: true });
    } finally {
      setChecking(false);
    }
  }

  function markWrong(ex: Exercise, d: FeedbackDetail = {}) {
    setGradedCount((c) => c + 1);
    setDetail(d);
    setFeedback("incorrect");
    // Spend a heart on the server; if that empties them, the lesson ends.
    lose().then((s) => {
      if (s && s.hearts <= 0) setOutOfHearts(true);
    });
    // Remember the miss so it shows up in Practice → Review Mistakes.
    api
      .recordMistake({
        prompt: ex.prompt || "Choose the correct answer",
        question: ex.question,
        correctAnswer: ex.correctAnswer,
      })
      .catch(() => {});
    // ...and collect it for the end-of-lesson review.
    setMisses((m) => [
      ...m,
      {
        prompt: ex.prompt,
        question: ex.question,
        correctAnswer: ex.correctAnswer,
        playText: ex.type === "listen" ? ex.question : undefined,
        speaker: ex.character,
      },
    ]);
  }

  function advance() {
    setFeedback(null);
    setDetail({});
    setWriteProblems([]);
    setAnswer("");
    if (idx + 1 < total) {
      setIdx((i) => i + 1);
    } else if (misses.length > 0) {
      setPhase("review"); // study misses before the completion screen
    } else {
      finish();
    }
  }

  async function finish() {
    setSubmitting(true);
    const accuracy =
      gradedCount > 0 ? Math.round((correctCount / gradedCount) * 100) : 100;
    try {
      const res = await api.completeLesson(id, accuracy);
      sessionStorage.setItem(
        "lumora_lesson_result",
        JSON.stringify({
          xp: res.xpEarned,
          accuracy: res.accuracy,
          firstClear: res.firstClear,
        })
      );
    } catch {
      sessionStorage.setItem(
        "lumora_lesson_result",
        JSON.stringify({ xp: lesson?.exercises?.length ?? 0, accuracy, firstClear: true })
      );
    } finally {
      router.replace(`/lesson/${id}/complete`);
    }
  }

  // Out of hearts (either on entry, or after a wrong answer) ends the lesson.
  const heartsModal = (
    <OutOfHeartsModal
      status={heartsStatus}
      secondsToNext={secondsToNext}
      note="You ran out of hearts, so this lesson has ended. Refill to try again now, or wait for a heart and restart it."
      closeLabel="Back to lessons"
      buying={buyingHearts}
      onBuy={() => {
        setBuyingHearts(true);
        buy();
      }}
      onClose={() => {
        stopSpeaking();
        router.push("/learn");
      }}
    />
  );

  if (!lesson) {
    return (
      <div className="flex min-h-[100dvh] items-center justify-center bg-cream">
        <FoxMascot size={120} glow />
      </div>
    );
  }

  return (
    <div className="flex min-h-[100dvh] w-full justify-center bg-cream lg:bg-[#eceaf3]">
      <div className="flex min-h-[100dvh] w-full max-w-2xl flex-col bg-cream lg:my-8 lg:min-h-[calc(100dvh-4rem)] lg:overflow-hidden lg:rounded-3xl lg:shadow-card-lg">
      {/* Top bar */}
      <header className="flex items-center gap-3 px-4 pb-2 pt-12 lg:px-8 lg:pt-6">
        <button onClick={() => router.push("/home")} aria-label="Close lesson">
          <X className="text-gray-500" />
        </button>
        <div className="h-2 flex-1 overflow-hidden rounded-full bg-gray-100">
          <motion.div
            className="h-full rounded-full bg-purple"
            animate={{ width: `${progress}%` }}
            transition={{ duration: 0.3 }}
          />
        </div>
        <HeartIndicator hearts={hearts} secondsToNext={secondsToNext} />
      </header>
      {outOfHearts && heartsModal}

      <div className="flex flex-1 flex-col px-5 pt-4 lg:px-8">
        {phase === "vocab" ? (
          <VocabPhase
            vocab={lesson.vocab || []}
            onDone={() => setPhase("practice")}
          />
        ) : phase === "review" ? (
          <MistakesReview
            items={misses}
            finishLabel={submitting ? "Finishing…" : "Finish lesson"}
            onDone={finish}
          />
        ) : (
        <>
        {ex && (
          <>
            {!isNarrative && (
              <p className="text-label-sm font-bold uppercase tracking-wide text-gray-500">
                {ex.prompt}
              </p>
            )}

            {/* Question area */}
            <div className="mt-3">
              {isNarrative ? (
                <NarrativeCard ex={ex} />
              ) : (
                <QuestionArea ex={ex} />
              )}
            </div>

            {/* Answer area */}
            {!isNarrative && (
              <div className="mt-5">
                {needsChoice && (
                  <ChoiceList
                    ex={ex}
                    answer={answer}
                    feedback={feedback}
                    onSelect={(v) => !feedback && setAnswer(v)}
                  />
                )}
                {needsTyping && (
                  <input
                    key={idx}
                    autoFocus
                    autoCapitalize="off"
                    autoCorrect="off"
                    spellCheck={false}
                    onKeyDown={(e) => e.key === "Enter" && canCheck && check()}
                    value={answer}
                    onChange={(e) => setAnswer(e.target.value)}
                    disabled={!!feedback}
                    placeholder="Type your answer…"
                    className="h-[52px] w-full rounded-lg border border-gray-100 bg-white px-4 text-body-lg outline-none focus:border-purple"
                  />
                )}
                {isSpeak && <SpeakControl phrase={ex.question} disabled={!!feedback} />}
                {isWrite && (
                  <WriteControl
                    value={answer}
                    onChange={(v) => {
                      setAnswer(v);
                      if (writeProblems.length) setWriteProblems([]);
                    }}
                    words={wordCount}
                    example={ex.correctAnswer}
                    problems={writeProblems}
                    disabled={!!feedback}
                  />
                )}
              </div>
            )}
          </>
        )}

        {/* Bottom action / feedback */}
        <div className="mt-auto pb-8 pt-6">
          <AnimatePresence mode="wait">
            {feedback ? (
              <FeedbackBar
                key="fb"
                feedback={feedback}
                correctAnswer={ex?.correctAnswer || ""}
                detail={detail}
                onContinue={advance}
              />
            ) : (
              <Button
                key="check"
                full
                disabled={!canCheck}
                loading={submitting || checking}
                onClick={check}
              >
                {isNarrative
                  ? "Continue"
                  : isSpeak
                  ? "I said it!"
                  : isWrite
                  ? "Check my writing"
                  : "Check"}
              </Button>
            )}
          </AnimatePresence>
        </div>
        </>
        )}
      </div>
      </div>
    </div>
  );
}

/* ---------- sub-components ---------- */

function VocabPhase({
  vocab,
  onDone,
}: {
  vocab: import("@/lib/types").VocabItem[];
  onDone: () => void;
}) {
  const [i, setI] = useState(0);
  const item = vocab[i];
  const last = i === vocab.length - 1;

  // Auto-play the word in its character's voice when the card appears.
  useEffect(() => {
    if (item) speakAs(item.speaker || "Lumora", item.word);
    return () => stopSpeaking();
  }, [item]);

  if (!item) return null;

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-label-sm font-bold uppercase tracking-wide text-gray-500">
        New words · {i + 1}/{vocab.length}
      </p>

      <div className="mt-4 flex flex-1 flex-col items-center justify-center">
        <AnimatePresence mode="wait">
          <motion.div
            key={item.id}
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -12 }}
            className="w-full max-w-md rounded-2xl border border-gray-100 bg-white p-6 text-center shadow-card-lg"
          >
            <p className="text-display-lg font-extrabold text-ink">{item.word}</p>
            <p className="mt-1 text-heading-sm font-bold text-purple">
              {item.translation}
            </p>

            <button
              onClick={() => speakAs(item.speaker || "Lumora", item.word)}
              className="mx-auto mt-4 flex h-12 w-12 items-center justify-center rounded-full bg-purple text-white shadow-float"
              aria-label="Hear the word"
            >
              <Volume2 size={22} />
            </button>

            {item.example && (
              <button
                onClick={() => speakAs(item.speaker || "Lumora", item.example)}
                className="mt-5 w-full rounded-xl bg-gray-50 p-3 text-left transition hover:bg-gray-100"
              >
                <span className="block text-body-md font-semibold text-ink">
                  {item.example}
                </span>
                <span className="mt-0.5 block text-body-sm text-slatey">
                  {item.exampleTranslation}
                </span>
              </button>
            )}
            {item.speaker && (
              <div className="mt-4 flex justify-center">
                <SpeakerChip name={item.speaker} />
              </div>
            )}
          </motion.div>
        </AnimatePresence>
      </div>

      {/* progress dots */}
      <div className="mb-4 mt-4 flex justify-center gap-1.5">
        {vocab.map((_, n) => (
          <span
            key={n}
            className={`h-1.5 rounded-full transition-all ${
              n === i ? "w-5 bg-purple" : "w-1.5 bg-gray-200"
            }`}
          />
        ))}
      </div>

      <div className="flex gap-3 pb-8">
        {i > 0 && (
          <Button
            variant="outline"
            className="flex-1"
            onClick={() => setI((n) => Math.max(0, n - 1))}
          >
            Back
          </Button>
        )}
        <Button
          full={i === 0}
          className={i > 0 ? "flex-1" : ""}
          onClick={() => (last ? onDone() : setI((n) => n + 1))}
        >
          {last ? "Start practice" : "Next word"}
        </Button>
      </div>
    </div>
  );
}

function NarrativeCard({ ex }: { ex: Exercise }) {
  // The character speaks their line aloud, in their own voice, on appear.
  useEffect(() => {
    if (ex.question) speakAs(ex.character || "Lumora", ex.question);
    return () => stopSpeaking();
  }, [ex.id, ex.character, ex.question]);

  return (
    <div className="flex flex-col items-center gap-3 pt-6">
      <CharacterAvatar name={ex.character} />
      <SpeechBubble className="max-w-[300px] text-body-lg">{ex.question}</SpeechBubble>
      <button
        onClick={() => speakAs(ex.character || "Lumora", ex.question)}
        className="flex items-center gap-1.5 rounded-full bg-purple-light px-3 py-1.5 text-label-lg font-bold text-purple"
        aria-label="Replay voice"
      >
        <Volume2 size={16} /> Replay
      </button>
    </div>
  );
}

function QuestionArea({ ex }: { ex: Exercise }) {
  if (ex.type === "listen") {
    return (
      <div className="flex flex-col items-center gap-3 py-4">
        <SpeakerChip name="Mira" />
        <button
          onClick={() => speakAs("Mira", ex.question)}
          className="flex h-16 w-16 items-center justify-center rounded-full bg-purple text-white shadow-float"
          aria-label="Play audio"
        >
          <Volume2 size={28} />
        </button>
        <p className="text-body-sm text-slatey">Tap to listen, then choose the meaning</p>
      </div>
    );
  }
  return (
    <p className="text-heading-lg font-bold leading-snug text-ink">{ex.question}</p>
  );
}

function ChoiceList({
  ex,
  answer,
  feedback,
  onSelect,
}: {
  ex: Exercise;
  answer: string;
  feedback: Feedback;
  onSelect: (v: string) => void;
}) {
  return (
    <div className="space-y-2">
      {(ex.options || []).map((opt) => {
        const selected = answer === opt;
        const isCorrect = opt === ex.correctAnswer;
        let cls = "border-gray-100 bg-white";
        if (feedback && isCorrect) cls = "border-teal bg-teal-light";
        else if (feedback && selected && !isCorrect) cls = "border-coral bg-coral-light";
        else if (selected) cls = "border-purple bg-purple-light";
        return (
          <button
            key={opt}
            onClick={() => onSelect(opt)}
            className={`flex h-14 w-full items-center rounded-md border-2 px-4 text-left text-body-lg font-semibold transition ${cls}`}
          >
            {opt}
          </button>
        );
      })}
    </div>
  );
}

function WriteControl({
  value,
  onChange,
  words,
  example,
  problems,
  disabled,
}: {
  value: string;
  onChange: (v: string) => void;
  words: number;
  example: string;
  /** Why the answer can't be accepted yet (copied example, too short…). */
  problems: string[];
  disabled: boolean;
}) {
  const [showExample, setShowExample] = useState(false);
  return (
    <div>
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        rows={6}
        placeholder="Write your answer here…"
        // Pasting the example back in is exactly what's not allowed.
        onPaste={(e) => {
          const pasted = e.clipboardData.getData("text");
          if (example && copiesExample(pasted, example)) e.preventDefault();
        }}
        aria-invalid={problems.length > 0}
        className={`w-full rounded-xl border bg-white p-4 text-body-lg outline-none transition focus:border-purple ${
          problems.length ? "border-coral" : "border-gray-100"
        }`}
      />
      {problems.length > 0 && (
        <ul className="mt-2 space-y-1 rounded-xl bg-coral-light px-4 py-3 text-body-sm font-semibold text-coral">
          {problems.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      )}
      <div className="mt-2 flex items-center justify-between">
        <span className="text-body-sm text-slatey">{words} words</span>
        {example && (
          <button
            type="button"
            onClick={() => setShowExample((s) => !s)}
            className="text-body-sm font-semibold text-teal"
          >
            {showExample ? "Hide example" : "Show example"}
          </button>
        )}
      </div>
      {showExample && example && (
        <pre className="mt-2 whitespace-pre-wrap rounded-xl bg-gray-50 p-3 text-body-sm text-slatey">
          {example}
        </pre>
      )}
    </div>
  );
}

function SpeakControl({ phrase, disabled }: { phrase: string; disabled: boolean }) {
  const [listening, setListening] = useState(false);
  const [heard, setHeard] = useState<string | null>(null);
  const [score, setScore] = useState<number | null>(null);
  const [supported] = useState(() => speechRecognitionSupported());

  async function record() {
    if (listening || disabled) return;
    setHeard(null);
    setScore(null);
    setListening(true);
    try {
      const said = await recognizeSpeech();
      setHeard(said);
      setScore(scorePronunciation(phrase, said));
    } catch {
      setHeard("");
      setScore(null);
    } finally {
      setListening(false);
    }
  }

  const tone =
    score == null ? "" : score >= 80 ? "text-teal" : score >= 50 ? "text-amber" : "text-coral";
  const label =
    score == null ? "" : score >= 80 ? "Excellent!" : score >= 50 ? "Good try!" : "Keep practising";

  return (
    <div className="flex flex-col items-center gap-3 py-2">
      <p className="text-label-sm font-bold uppercase tracking-wide text-coral">
        Speaking practice
      </p>
      <p className="text-heading-md font-bold text-purple">&ldquo;{phrase}&rdquo;</p>

      <div className="flex items-center gap-3">
        <button
          onClick={() => speakAs("Lumora", phrase)}
          disabled={disabled}
          className="flex h-12 w-12 items-center justify-center rounded-full bg-purple text-white shadow-float disabled:opacity-40"
          aria-label="Hear it"
        >
          <Volume2 size={22} />
        </button>
        <button
          onClick={record}
          disabled={disabled || listening || !supported}
          className={`flex h-16 w-16 items-center justify-center rounded-full text-white shadow-float disabled:opacity-40 ${
            listening ? "animate-pulse bg-coral" : "bg-coral"
          }`}
          aria-label="Tap to speak"
        >
          <Mic size={28} />
        </button>
      </div>

      {supported ? (
        <p className="text-body-sm text-slatey">
          {listening ? "Listening… say the phrase" : "Tap the mic and say it out loud"}
        </p>
      ) : (
        <p className="text-body-sm text-slatey">
          Speech scoring needs Chrome/Edge. Say it aloud, then tap “I said it!”
        </p>
      )}

      {/* Fluency measure */}
      {score != null && (
        <div className="mt-1 w-full max-w-xs rounded-xl bg-gray-50 p-3 text-center">
          <p className={`text-display-lg font-extrabold ${tone}`}>{score}%</p>
          <p className={`text-label-md font-bold uppercase tracking-wide ${tone}`}>
            {label} · fluency
          </p>
          {heard ? (
            <p className="mt-1 text-body-sm text-slatey">You said: “{heard}”</p>
          ) : (
            <p className="mt-1 text-body-sm text-slatey">
              Didn&apos;t catch that — try again.
            </p>
          )}
        </div>
      )}
    </div>
  );
}

function FeedbackBar({
  feedback,
  correctAnswer,
  detail,
  onContinue,
}: {
  feedback: Feedback;
  correctAnswer: string;
  detail: FeedbackDetail;
  onContinue: () => void;
}) {
  const correct = feedback === "correct";
  return (
    <motion.div
      initial={{ y: 24, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className={`rounded-xl p-4 ${correct ? "bg-teal-light" : "bg-coral-light"}`}
    >
      <div className="mb-3 flex items-center gap-2">
        <span
          className={`flex h-8 w-8 items-center justify-center rounded-full text-white ${
            correct ? "bg-teal" : "bg-coral"
          }`}
        >
          {correct ? <Check size={18} /> : <X size={18} />}
        </span>
        <div>
          <p className={`font-extrabold ${correct ? "text-teal" : "text-coral"}`}>
            {correct ? "Nailed it!" : "Not quite"}
          </p>
          {!correct && !detail.hideAnswer && (
            <p className="text-body-sm text-ink">
              Answer: <span className="font-bold">{correctAnswer}</span>
            </p>
          )}
        </div>
      </div>
      {detail.note && <p className="mb-2 text-body-sm text-ink">{detail.note}</p>}
      {detail.lines && detail.lines.length > 0 && (
        <ul className="mb-3 space-y-0.5 text-body-sm text-ink">
          {detail.lines.map((l) => (
            <li key={l}>• {l}</li>
          ))}
        </ul>
      )}
      {detail.corrections && detail.corrections.length > 0 && (
        <ul className="mb-3 max-h-56 space-y-2 overflow-y-auto">
          {detail.corrections.map((c, i) => (
            <li key={i} className="rounded-lg bg-white/70 px-3 py-2 text-body-sm">
              <span className="text-coral line-through">{c.original}</span>{" "}
              → <span className="font-bold text-teal">{c.correction}</span>
              <span className="block text-slatey">{c.explanation}</span>
            </li>
          ))}
        </ul>
      )}
      <Button
        full
        variant={correct ? "primary" : "danger"}
        onClick={onContinue}
      >
        {correct ? "Continue" : "Got it"}
      </Button>
    </motion.div>
  );
}

function CharacterAvatar({ name }: { name: string }) {
  const info = characterInfo(name);
  return (
    <div className="flex flex-col items-center gap-1.5">
      <div
        className="rounded-full border-4 p-0.5"
        style={{ borderColor: info.color }}
      >
        <SpeakerAvatar name={name} size={92} />
      </div>
      <span className="text-body-md font-extrabold text-ink">{info.name}</span>
    </div>
  );
}
