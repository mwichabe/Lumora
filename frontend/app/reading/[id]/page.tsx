"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { X, Volume2, Check } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { speakAs, stopSpeaking } from "@/lib/voices";
import type { ReadingSession } from "@/lib/types";

type Phase = "read" | "quiz" | "done";

export default function ReadingPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { setUser } = useAuth();

  const [session, setSession] = useState<ReadingSession | null>(null);
  const [phase, setPhase] = useState<Phase>("read");

  useEffect(() => {
    api
      .readingSession(id)
      .then((d) => setSession(d.session))
      .catch(() => {});
    return () => stopSpeaking();
  }, [id]);

  if (!session) {
    return (
      <div className="flex min-h-[100dvh] items-center justify-center bg-cream">
        <FoxMascot size={120} glow />
      </div>
    );
  }

  return (
    <div className="flex min-h-[100dvh] w-full justify-center bg-cream lg:bg-[#eceaf3]">
      <div className="flex min-h-[100dvh] w-full max-w-2xl flex-col bg-cream lg:my-8 lg:min-h-[calc(100dvh-4rem)] lg:overflow-hidden lg:rounded-3xl lg:shadow-card-lg">
        <header className="flex items-center gap-3 px-4 pb-3 pt-12 lg:px-8 lg:pt-6">
          <button
            onClick={() => {
              stopSpeaking();
              router.push("/learn");
            }}
            aria-label="Close"
          >
            <X className="text-gray-500" />
          </button>
          <div>
            <p className="text-label-md font-bold uppercase tracking-wide text-gray-500">
              Reading · {session.unit}
            </p>
            <h1 className="text-heading-sm font-extrabold text-ink">
              {session.title}
            </h1>
          </div>
        </header>

        <div className="flex flex-1 flex-col px-5 pb-8 pt-2 lg:px-8">
          {phase === "read" && (
            <ReadPhase session={session} onDone={() => setPhase("quiz")} />
          )}
          {phase === "quiz" && (
            <QuizPhase
              session={session}
              onDone={async () => {
                try {
                  const r = await api.completeReading(session.id);
                  setUser(r.user);
                } catch {
                  /* ignore */
                }
                setPhase("done");
              }}
            />
          )}
          {phase === "done" && (
            <DonePhase
              xp={session.xpReward}
              onContinue={() => router.push("/learn")}
            />
          )}
        </div>
      </div>
    </div>
  );
}

function ReadPhase({
  session,
  onDone,
}: {
  session: ReadingSession;
  onDone: () => void;
}) {
  const lines = session.lines || [];
  const [revealed, setRevealed] = useState(false);

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-body-md text-slatey">{session.description}</p>

      <div className="mt-4 space-y-3">
        {lines.map((l) => (
          <div
            key={l.id}
            className="rounded-2xl border border-gray-100 bg-white p-4 shadow-card"
          >
            <div className="flex items-start justify-between gap-3">
              <p className="text-body-lg font-semibold leading-relaxed text-ink">
                {l.text}
              </p>
              <button
                onClick={() => speakAs("Mira", l.text)}
                aria-label="Hear sentence"
                className="shrink-0 text-purple"
              >
                <Volume2 size={18} />
              </button>
            </div>
            {revealed && (
              <p className="mt-1 text-body-sm text-slatey">{l.translation}</p>
            )}
          </div>
        ))}
      </div>

      <button
        onClick={() => setRevealed((r) => !r)}
        className="mt-4 text-center text-body-sm font-semibold text-teal"
      >
        {revealed ? "Hide translations" : "Show translations"}
      </button>

      <div className="mt-auto pt-6">
        <Button full onClick={onDone}>
          I&apos;ve read it — quiz me
        </Button>
      </div>
    </div>
  );
}

function QuizPhase({
  session,
  onDone,
}: {
  session: ReadingSession;
  onDone: () => void;
}) {
  const questions = session.questions || [];
  const [qi, setQi] = useState(0);
  const [answer, setAnswer] = useState("");
  const [feedback, setFeedback] = useState<null | "correct" | "incorrect">(null);

  const q = questions[qi];
  useEffect(() => {
    if (!q) onDone();
  }, [q, onDone]);
  if (!q) return null;

  function check() {
    if (!answer) return;
    setFeedback(answer === q.correctAnswer ? "correct" : "incorrect");
  }
  function next() {
    setFeedback(null);
    setAnswer("");
    if (qi + 1 < questions.length) setQi((n) => n + 1);
    else onDone();
  }

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-label-sm font-bold uppercase tracking-wide text-gray-500">
        Comprehension · {qi + 1}/{questions.length}
      </p>
      <p className="mt-3 text-heading-lg font-bold leading-snug text-ink">
        {q.question}
      </p>

      <div className="mt-5 space-y-2">
        {(q.options || []).map((opt) => {
          const selected = answer === opt;
          const isCorrect = opt === q.correctAnswer;
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

      <div className="mt-auto pb-2 pt-6">
        <AnimatePresence mode="wait">
          {feedback ? (
            <motion.div
              key="fb"
              initial={{ y: 20, opacity: 0 }}
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
                {feedback === "correct" ? "Correct!" : `Answer: ${q.correctAnswer}`}
              </p>
              <Button
                full
                variant={feedback === "correct" ? "primary" : "danger"}
                onClick={next}
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

function DonePhase({ xp, onContinue }: { xp: number; onContinue: () => void }) {
  return (
    <div className="flex flex-1 flex-col items-center justify-center text-center">
      <FoxMascot size={130} glow bounce />
      <h2 className="mt-4 text-heading-xl font-extrabold text-ink">
        Reading complete!
      </h2>
      <p className="mt-1 text-body-md text-slatey">
        You read real Spanish. ¡Muy bien!
      </p>
      <span className="mt-5 flex h-9 items-center gap-1.5 rounded-full bg-amber px-4 font-extrabold text-ink">
        <Check size={16} /> +{xp} XP
      </span>
      <div className="mt-8 w-full max-w-sm">
        <Button full onClick={onContinue}>
          Back to course
        </Button>
      </div>
    </div>
  );
}
