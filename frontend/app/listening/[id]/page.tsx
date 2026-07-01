"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { X, Play, Pause, Volume2, Check } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { SpeakerAvatar } from "@/components/Speaker";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { speakAs, speakSequence, stopSpeaking } from "@/lib/voices";
import type { ListeningSession } from "@/lib/types";

type Phase = "match" | "listen" | "quiz" | "done";

export default function ListeningPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { setUser } = useAuth();

  const [session, setSession] = useState<ListeningSession | null>(null);
  const [phase, setPhase] = useState<Phase>("match");

  useEffect(() => {
    api
      .listeningSession(id)
      .then((d) => {
        setSession(d.session);
        // Skip the warm-up if a session has no matching pairs.
        setPhase(d.session.matches && d.session.matches.length ? "match" : "listen");
      })
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
              Listening · {session.unit}
            </p>
            <h1 className="text-heading-sm font-extrabold text-ink">
              {session.title}
            </h1>
          </div>
        </header>

        <div className="flex flex-1 flex-col px-5 pb-8 pt-2 lg:px-8">
          {phase === "match" && (
            <MatchPhase session={session} onDone={() => setPhase("listen")} />
          )}
          {phase === "listen" && (
            <ListenPhase session={session} onDone={() => setPhase("quiz")} />
          )}
          {phase === "quiz" && (
            <QuizPhase
              session={session}
              onDone={async () => {
                try {
                  const r = await api.completeListening(session.id);
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

function MatchPhase({
  session,
  onDone,
}: {
  session: ListeningSession;
  onDone: () => void;
}) {
  const pairs = session.matches || [];

  // Stable shuffle of the target-word column (seeded by id so it doesn't
  // re-order on every render).
  const shuffled = useRef(
    [...pairs].sort((a, b) => ((a.id * 7) % 5) - ((b.id * 7) % 5))
  ).current;

  const [selected, setSelected] = useState<number | null>(null); // left (English) id
  const [matched, setMatched] = useState<Set<number>>(new Set());
  const [wrong, setWrong] = useState<number | null>(null); // right id flashing red

  const allDone = matched.size === pairs.length && pairs.length > 0;

  function tapRight(rightId: number, word: string) {
    if (selected == null || matched.has(rightId)) return;
    const left = pairs.find((p) => p.id === selected);
    if (left && left.word === word) {
      const next = new Set(matched);
      next.add(rightId); // right id === pair id
      setMatched(next);
      setSelected(null);
      speakAs(undefined, word);
    } else {
      setWrong(rightId);
      setTimeout(() => setWrong(null), 500);
    }
  }

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-label-sm font-bold uppercase tracking-wide text-gray-500">
        Warm-up · match the words
      </p>
      <p className="mt-1 text-body-md text-slatey">
        These words appear in the conversation you&apos;re about to hear.
      </p>

      <div className="mt-5 grid flex-1 grid-cols-2 gap-3">
        {/* English column */}
        <div className="space-y-2">
          {pairs.map((p) => {
            const done = matched.has(p.id);
            const active = selected === p.id;
            return (
              <button
                key={p.id}
                disabled={done}
                onClick={() => setSelected(p.id)}
                className={`flex h-14 w-full items-center rounded-xl border-2 px-3 text-left text-body-md font-semibold transition ${
                  done
                    ? "border-teal bg-teal-light text-teal"
                    : active
                    ? "border-purple bg-purple-light"
                    : "border-gray-100 bg-white"
                }`}
              >
                {p.translation}
              </button>
            );
          })}
        </div>

        {/* Target-language column */}
        <div className="space-y-2">
          {shuffled.map((p) => {
            const done = matched.has(p.id);
            const isWrong = wrong === p.id;
            return (
              <button
                key={p.id}
                disabled={done}
                onClick={() => tapRight(p.id, p.word)}
                className={`flex h-14 w-full items-center justify-between gap-2 rounded-xl border-2 px-3 text-left text-body-md font-semibold transition ${
                  done
                    ? "border-teal bg-teal-light text-teal"
                    : isWrong
                    ? "border-coral bg-coral-light"
                    : "border-gray-100 bg-white"
                }`}
              >
                <span className="truncate">{p.word}</span>
                {done && <Check size={16} className="shrink-0" />}
              </button>
            );
          })}
        </div>
      </div>

      <div className="mt-4">
        <Button full disabled={!allDone} onClick={onDone}>
          {allDone ? "Listen to the conversation" : "Match all the words"}
        </Button>
      </div>
    </div>
  );
}

function ListenPhase({
  session,
  onDone,
}: {
  session: ListeningSession;
  onDone: () => void;
}) {
  const lines = session.lines || [];
  const [active, setActive] = useState(-1);
  const [playing, setPlaying] = useState(false);
  const [revealed, setRevealed] = useState(false);
  const playingRef = useRef(false);

  async function playAll() {
    if (playing) {
      playingRef.current = false;
      stopSpeaking();
      setPlaying(false);
      return;
    }
    setPlaying(true);
    playingRef.current = true;
    await speakSequence(
      lines.map((l) => ({ character: l.character, text: l.text })),
      (i) => setActive(i),
      () => playingRef.current
    );
    setPlaying(false);
    setActive(-1);
  }

  useEffect(() => () => {
    playingRef.current = false;
    stopSpeaking();
  }, []);

  return (
    <div className="flex flex-1 flex-col">
      <p className="text-body-md text-slatey">{session.description}</p>

      <button
        onClick={playAll}
        className="mt-4 flex items-center justify-center gap-2 rounded-full bg-purple py-3 font-extrabold text-white shadow-float"
      >
        {playing ? <Pause size={20} /> : <Play size={20} fill="#fff" strokeWidth={0} />}
        {playing ? "Pause" : "Play conversation"}
      </button>

      <div className="mt-5 space-y-3">
        {lines.map((l, i) => {
          const isActive = i === active;
          return (
            <motion.div
              key={l.id}
              animate={{ scale: isActive ? 1.0 : 1, opacity: isActive || active === -1 ? 1 : 0.55 }}
              className={`flex items-start gap-3 rounded-2xl border p-3 transition ${
                isActive ? "border-purple bg-purple-light" : "border-gray-100 bg-white"
              }`}
            >
              <SpeakerAvatar name={l.character} size={40} />
              <div className="min-w-0 flex-1">
                <p className="text-label-md font-bold text-slatey">{l.character}</p>
                <p className="text-body-lg font-semibold text-ink">{l.text}</p>
                {revealed && (
                  <p className="mt-0.5 text-body-sm text-slatey">{l.translation}</p>
                )}
              </div>
              <button
                onClick={() => speakAs(l.character, l.text)}
                aria-label="Play line"
                className="text-purple"
              >
                <Volume2 size={18} />
              </button>
            </motion.div>
          );
        })}
      </div>

      <button
        onClick={() => setRevealed((r) => !r)}
        className="mt-4 text-center text-body-sm font-semibold text-teal"
      >
        {revealed ? "Hide translations" : "Show translations"}
      </button>

      <div className="mt-auto pt-6">
        <Button full onClick={onDone}>
          I&apos;m ready — quiz me
        </Button>
      </div>
    </div>
  );
}

function QuizPhase({
  session,
  onDone,
}: {
  session: ListeningSession;
  onDone: () => void;
}) {
  const questions = session.questions || [];
  const [qi, setQi] = useState(0);
  const [answer, setAnswer] = useState("");
  const [feedback, setFeedback] = useState<null | "correct" | "incorrect">(null);

  const q = questions[qi];

  // If there are no questions, finish — but do it in an effect, never during
  // render (calling onDone() in render updates the parent mid-render = crash).
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
                {feedback === "correct"
                  ? "Correct!"
                  : `Answer: ${q.correctAnswer}`}
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
        Session complete!
      </h2>
      <p className="mt-1 text-body-md text-slatey">
        Your ear is getting sharper. ¡Bien hecho!
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
