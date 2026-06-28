"use client";

import { useEffect, useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { X, Volume2, Check, Mic } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { SpeechBubble, HeartIndicator } from "@/components/widgets";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import {
  speakAs,
  stopSpeaking,
  recognizeSpeech,
  scorePronunciation,
  speechRecognitionSupported,
} from "@/lib/voices";
import type { Lesson, Exercise } from "@/lib/types";

type Feedback = null | "correct" | "incorrect";

export default function LessonPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { user } = useAuth();

  const [lesson, setLesson] = useState<Lesson | null>(null);
  const [phase, setPhase] = useState<"vocab" | "practice">("vocab");
  const [idx, setIdx] = useState(0);
  const [hearts, setHearts] = useState(user?.hearts ?? 5);
  const [answer, setAnswer] = useState("");
  const [feedback, setFeedback] = useState<Feedback>(null);
  const [correctCount, setCorrectCount] = useState(0);
  const [gradedCount, setGradedCount] = useState(0);
  const [submitting, setSubmitting] = useState(false);

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

  // Stop any speech when leaving the lesson.
  useEffect(() => () => stopSpeaking(), []);

  const exercises = lesson?.exercises || [];
  const ex: Exercise | undefined = exercises[idx];
  const total = exercises.length;
  const progress = total ? (idx / total) * 100 : 0;

  const isNarrative = ex?.type === "character";
  const isSpeak = ex?.type === "speak";
  // Translate/fill now arrive with generated options, so render them as choices
  // too. Only fall back to typing if no options were provided.
  const hasOptions = !!(ex?.options && ex.options.length > 0);
  const needsChoice =
    !!ex &&
    !isNarrative &&
    !isSpeak &&
    (["multiple_choice", "listen", "match"].includes(ex.type) || hasOptions);
  const needsTyping =
    !!ex && ["translate", "fill"].includes(ex.type) && !hasOptions;

  const canCheck = useMemo(() => {
    if (!ex || feedback) return false;
    if (isNarrative || isSpeak) return true;
    return answer.trim().length > 0;
  }, [ex, feedback, answer, isNarrative, isSpeak]);

  function normalise(s: string) {
    return s.trim().toLowerCase().replace(/[.,!¡¿?]/g, "");
  }

  function check() {
    if (!ex) return;
    if (isNarrative || isSpeak) {
      advance();
      return;
    }
    const ok = normalise(answer) === normalise(ex.correctAnswer);
    setGradedCount((c) => c + 1);
    if (ok) {
      setCorrectCount((c) => c + 1);
      setFeedback("correct");
    } else {
      setHearts((h) => Math.max(0, h - 1));
      setFeedback("incorrect");
    }
  }

  function advance() {
    setFeedback(null);
    setAnswer("");
    if (idx + 1 < total) {
      setIdx((i) => i + 1);
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
        <HeartIndicator hearts={hearts} />
      </header>

      <div className="flex flex-1 flex-col px-5 pt-4 lg:px-8">
        {phase === "vocab" ? (
          <VocabPhase
            vocab={lesson.vocab || []}
            onDone={() => setPhase("practice")}
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
                    autoFocus
                    value={answer}
                    onChange={(e) => setAnswer(e.target.value)}
                    disabled={!!feedback}
                    placeholder="Type your answer…"
                    className="h-[52px] w-full rounded-lg border border-gray-100 bg-white px-4 text-body-lg outline-none focus:border-purple"
                  />
                )}
                {isSpeak && <SpeakControl phrase={ex.question} disabled={!!feedback} />}
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
                onContinue={advance}
              />
            ) : (
              <Button
                key="check"
                full
                disabled={!canCheck}
                loading={submitting}
                onClick={check}
              >
                {isNarrative ? "Continue" : isSpeak ? "I said it!" : "Check"}
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
              <p className="mt-3 text-label-md text-gray-500">
                Voiced by {item.speaker}
              </p>
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
      <div className="flex flex-col items-center gap-4 py-4">
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
      const said = await recognizeSpeech("es-ES");
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
  onContinue,
}: {
  feedback: Feedback;
  correctAnswer: string;
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
          {!correct && (
            <p className="text-body-sm text-ink">
              Answer: <span className="font-bold">{correctAnswer}</span>
            </p>
          )}
        </div>
      </div>
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
  const map: Record<string, string> = {
    Lumora: "🦊",
    "Professor Finch": "🦅",
    Cora: "🐙",
    Blaze: "🔥",
    Mira: "🐆",
    Riko: "🐼",
    Zephyr: "🌬️",
    Nana: "🐢",
    Pip: "🦔",
  };
  if (name === "Lumora") return <FoxMascot size={120} glow bounce />;
  return (
    <div className="flex h-24 w-24 items-center justify-center rounded-full border-4 border-purple bg-white text-5xl">
      {map[name] || "🦊"}
    </div>
  );
}
