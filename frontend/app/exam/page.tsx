"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  X,
  Volume2,
  Mic,
  Headphones,
  BookText,
  PenLine,
  GraduationCap,
  Check,
  Clock,
  Camera,
  MonitorUp,
  ShieldCheck,
  AlertTriangle,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { languageName } from "@/lib/languages";
import {
  speakAs,
  speakSequence,
  stopSpeaking,
  recognizeSpeech,
  scorePronunciation,
  speechRecognitionSupported,
} from "@/lib/voices";
import type {
  ListeningSession,
  ReadingSession,
  VocabItem,
  ExamResult,
  ExamMeta,
} from "@/lib/types";

type Phase =
  | "intro"
  | "rules"
  | "listening"
  | "reading"
  | "writing"
  | "speaking"
  | "result";

type Termination = "tab" | "time" | null;

const LEVELS = [
  { code: "A1", name: "Beginner" },
  { code: "A2", name: "Elementary" },
  { code: "B1", name: "Intermediate" },
  { code: "B2", name: "Upper-Int." },
  { code: "C1", name: "Advanced" },
  { code: "C2", name: "Mastery" },
];
const MIN_WORDS = [8, 12, 18, 26, 36, 48];
const FALLBACK_DURATION = [420, 540, 660, 780, 900, 1020];
const FALLBACK_PASS = [50, 58, 65, 72, 80, 88];
const WRITING_PROMPTS = [
  "introducing yourself — your name, where you're from, and what you like.",
  "describing your daily routine and your plans for the weekend.",
  "writing a short email to a friend about a recent trip you took.",
  "giving your opinion on living in a city versus the countryside, with reasons.",
  "discussing the advantages and disadvantages of working from home.",
  "writing a structured argument on how technology is reshaping society.",
];
const LOCALE: Record<string, string> = {
  es: "es-ES",
  fr: "fr-FR",
  de: "de-DE",
  it: "it-IT",
  pt: "pt-PT",
};

const SECTION_PHASES: Phase[] = ["listening", "reading", "writing", "speaking"];
const levelIdx = (code: string) =>
  Math.max(0, LEVELS.findIndex((l) => l.code === code));
const fmtTime = (s: number) =>
  `${Math.floor(Math.max(0, s) / 60)}:${String(Math.max(0, s) % 60).padStart(2, "0")}`;

/** Pick the content best matching the chosen level (hardest session in that band). */
function pickForLevel<T extends { unit?: string; orderIndex: number }>(
  arr: T[],
  level: string
): T | undefined {
  if (!arr.length) return undefined;
  const code = level.toUpperCase();
  const exact = arr.filter((s) => (s.unit || "").toUpperCase().startsWith(code));
  if (exact.length) return exact[exact.length - 1];
  const sorted = [...arr].sort((a, b) => a.orderIndex - b.orderIndex);
  const i = Math.min(
    sorted.length - 1,
    Math.round((levelIdx(level) / 5) * (sorted.length - 1))
  );
  return sorted[i];
}

export default function ExamPage() {
  return (
    <AppShell tabs>
      <ExamRunner />
    </AppShell>
  );
}

function ExamRunner() {
  const router = useRouter();
  const { user } = useAuth();
  const lang = user?.targetLanguage || "";

  const [loading, setLoading] = useState(true);
  const [allListening, setAllListening] = useState<ListeningSession[]>([]);
  const [allReading, setAllReading] = useState<ReadingSession[]>([]);
  const [vocab, setVocab] = useState<VocabItem[]>([]);
  const [meta, setMeta] = useState<ExamMeta | null>(null);
  const [completed, setCompleted] = useState<string[]>([]);

  const [level, setLevel] = useState("A1");
  const [phase, setPhase] = useState<Phase>("intro");
  const [starting, setStarting] = useState(false);
  const [scores, setScores] = useState({
    listening: 0,
    reading: 0,
    writing: 0,
    speaking: 0,
  });
  const [result, setResult] = useState<ExamResult | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Proctoring + timer state
  const [secondsLeft, setSecondsLeft] = useState(0);
  const [termination, setTermination] = useState<Termination>(null);
  const [camOn, setCamOn] = useState(false);
  const [screenOn, setScreenOn] = useState(false);

  const camStreamRef = useRef<MediaStream | null>(null);
  const screenStreamRef = useRef<MediaStream | null>(null);
  const endedRef = useRef(false);
  const activeRef = useRef(false);
  const scoresRef = useRef(scores);
  useEffect(() => {
    scoresRef.current = scores;
  }, [scores]);

  const idx = levelIdx(level);
  const isActive = SECTION_PHASES.includes(phase);
  useEffect(() => {
    activeRef.current = isActive;
  }, [isActive]);

  const duration =
    meta?.levels[level]?.durationSeconds ?? FALLBACK_DURATION[idx];
  const passMark = meta?.levels[level]?.passMark ?? FALLBACK_PASS[idx];
  const weights = meta?.weights ?? {
    listening: 30,
    reading: 30,
    writing: 20,
    speaking: 20,
  };

  /* ---------- load content ---------- */
  useEffect(() => {
    Promise.all([
      api.listeningSessions().catch(() => ({ sessions: [] })),
      api.readingSessions().catch(() => ({ sessions: [] })),
      api.practice().catch(() => ({ vocab: [], mistakes: [] })),
      api.certificates().catch(() => ({ certificates: [] })),
      api.examMeta().catch(() => null),
    ])
      .then(([ls, rs, p, cs, m]) => {
        setAllListening(ls.sessions || []);
        setAllReading(rs.sessions || []);
        setVocab((p.vocab || []).filter((v) => v.word));
        setMeta(m);
        setCompleted(
          cs.certificates.filter((c) => c.language === lang).map((c) => c.level)
        );
      })
      .finally(() => setLoading(false));
    return () => stopSpeaking();
  }, [lang]);

  /* ---------- proctoring streams ---------- */
  const stopProctoring = useCallback(() => {
    [camStreamRef, screenStreamRef].forEach((r) => {
      r.current?.getTracks().forEach((t) => t.stop());
      r.current = null;
    });
  }, []);

  useEffect(() => stopProctoring, [stopProctoring]);

  const submit = useCallback(
    async (final: typeof scores) => {
      setSubmitting(true);
      setPhase("result");
      try {
        const r = await api.submitExam({ language: lang, level, ...final });
        setResult(r);
      } catch {
        setResult(null);
      } finally {
        setSubmitting(false);
      }
    },
    [lang, level]
  );

  // Ends the exam immediately and scores whatever has been completed so far.
  const endExam = useCallback(
    (reason: Termination) => {
      if (endedRef.current) return;
      endedRef.current = true;
      activeRef.current = false;
      setTermination(reason);
      stopProctoring();
      stopSpeaking();
      submit(scoresRef.current);
    },
    [stopProctoring, submit]
  );

  /* ---------- tab-switch guard ---------- */
  useEffect(() => {
    function onVisibility() {
      if (document.hidden && activeRef.current && !endedRef.current) {
        endExam("tab");
      }
    }
    document.addEventListener("visibilitychange", onVisibility);
    return () =>
      document.removeEventListener("visibilitychange", onVisibility);
  }, [endExam]);

  /* ---------- countdown ---------- */
  useEffect(() => {
    if (!isActive) return;
    const id = setInterval(
      () => setSecondsLeft((s) => Math.max(0, s - 1)),
      1000
    );
    return () => clearInterval(id);
  }, [isActive]);

  useEffect(() => {
    if (isActive && secondsLeft <= 0 && !endedRef.current) endExam("time");
  }, [secondsLeft, isActive, endExam]);

  /* ---------- begin (request camera + screen, then start) ---------- */
  const beginExam = useCallback(async () => {
    setStarting(true);
    endedRef.current = false;
    setTermination(null);
    // Screen share first — it has the strictest user-gesture requirement.
    try {
      const scr = await navigator.mediaDevices.getDisplayMedia({ video: true });
      screenStreamRef.current = scr;
      setScreenOn(true);
    } catch {
      setScreenOn(false);
    }
    try {
      const cam = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: false,
      });
      camStreamRef.current = cam;
      setCamOn(true);
    } catch {
      setCamOn(false);
    }
    setSecondsLeft(duration);
    setStarting(false);
    setPhase("listening");
  }, [duration]);

  function advance(
    section: keyof typeof scores,
    score: number,
    nextPhase: Phase
  ) {
    if (endedRef.current) return;
    const updated = { ...scores, [section]: score };
    setScores(updated);
    if (nextPhase === "result") {
      endedRef.current = true;
      activeRef.current = false;
      stopProctoring();
      submit(updated);
    } else {
      setPhase(nextPhase);
    }
  }

  // Reset back to the level picker so the user can take another exam without a
  // full page navigation (a <Link href="/exam"> is a no-op when already here).
  function retake() {
    stopProctoring();
    stopSpeaking();
    endedRef.current = false;
    activeRef.current = false;
    setTermination(null);
    setResult(null);
    setSecondsLeft(0);
    setScores({ listening: 0, reading: 0, writing: 0, speaking: 0 });
    setPhase("intro");
  }

  function handleClose() {
    if (activeRef.current && !endedRef.current) {
      if (
        !window.confirm(
          "Leaving now will end your exam and score only what you've completed. Continue?"
        )
      )
        return;
    }
    stopProctoring();
    stopSpeaking();
    router.push("/profile");
  }

  const lsSession = useMemo(
    () => pickForLevel(allListening, level),
    [allListening, level]
  );
  const rsSession = useMemo(
    () => pickForLevel(allReading, level),
    [allReading, level]
  );
  const speakPick = useMemo(
    () => vocab[Math.floor(vocab.length / 2)] || vocab[0],
    [vocab]
  );
  const ready = !!lsSession && !!rsSession;
  const speakPhrase =
    idx >= 2
      ? speakPick?.example || speakPick?.word || ""
      : speakPick?.word || "";
  const locale = LOCALE[lang] || "en-US";

  return (
    <Shell
      onClose={handleClose}
      timer={isActive ? secondsLeft : null}
      camOn={isActive ? camOn : null}
      screenOn={isActive ? screenOn : null}
    >
      {loading ? (
        <Centered>
          <FoxMascot size={110} glow />
        </Centered>
      ) : !ready ? (
        <Centered>
          <FoxMascot size={110} glow />
          <p className="mt-4 text-heading-sm font-extrabold text-ink">
            Exam not ready for {languageName(lang)}
          </p>
          <p className="mt-1 max-w-xs text-body-md text-slatey">
            This language needs a full course (lessons, listening &amp; reading)
            before the exam is available.
          </p>
          <div className="mt-6 w-full max-w-sm">
            <Button full variant="outline" onClick={() => router.push("/learn")}>
              Go to lessons
            </Button>
          </div>
        </Centered>
      ) : phase === "intro" ? (
        <LevelSelect
          lang={lang}
          completed={completed}
          onStart={(lvl) => {
            setLevel(lvl);
            setScores({ listening: 0, reading: 0, writing: 0, speaking: 0 });
            setPhase("rules");
          }}
        />
      ) : phase === "rules" ? (
        <RulesScreen
          lang={lang}
          level={level}
          durationSeconds={duration}
          passMark={passMark}
          weights={weights}
          starting={starting}
          onBack={() => setPhase("intro")}
          onBegin={beginExam}
        />
      ) : phase === "listening" ? (
        <ListeningSection
          session={lsSession!}
          onDone={(s) => advance("listening", s, "reading")}
        />
      ) : phase === "reading" ? (
        <ReadingSection
          session={rsSession!}
          onDone={(s) => advance("reading", s, "writing")}
        />
      ) : phase === "writing" ? (
        <WritingSection
          lang={lang}
          minWords={MIN_WORDS[idx]}
          prompt={WRITING_PROMPTS[idx]}
          onDone={(s) => advance("writing", s, "speaking")}
        />
      ) : phase === "speaking" ? (
        <SpeakingSection
          phrase={speakPhrase}
          speaker={speakPick?.speaker || "Lumora"}
          locale={locale}
          onDone={(s) => advance("speaking", s, "result")}
        />
      ) : (
        <ResultView
          result={result}
          submitting={submitting}
          lang={lang}
          termination={termination}
          onRetake={retake}
        />
      )}

      {/* Proctoring camera preview — visible only during the exam */}
      {isActive && camOn && (
        <div className="fixed bottom-24 right-4 z-40 overflow-hidden rounded-xl border-2 border-white shadow-float lg:bottom-6">
          <video
            ref={(el) => {
              if (el && camStreamRef.current && el.srcObject !== camStreamRef.current)
                el.srcObject = camStreamRef.current;
            }}
            autoPlay
            muted
            playsInline
            className="h-24 w-32 -scale-x-100 bg-black object-cover"
          />
          <span className="absolute left-1.5 top-1.5 flex items-center gap-1 rounded-full bg-black/50 px-1.5 py-0.5 text-[10px] font-bold text-white">
            <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-coral" />
            REC
          </span>
        </div>
      )}
    </Shell>
  );
}

/* ---------- shell ---------- */

function Shell({
  onClose,
  timer,
  camOn,
  screenOn,
  children,
}: {
  onClose: () => void;
  timer: number | null;
  camOn: boolean | null;
  screenOn: boolean | null;
  children: React.ReactNode;
}) {
  const low = timer != null && timer <= 60;
  return (
    <div className="flex min-h-full flex-col bg-cream">
      <header className="flex items-center gap-3 border-b border-gray-100 bg-white px-4 py-3.5 lg:rounded-t-3xl lg:px-6">
        <button onClick={onClose} aria-label="Exit exam">
          <X className="text-gray-500" />
        </button>
        <h1 className="text-heading-sm font-extrabold text-ink">
          Proficiency Exam
        </h1>
        <div className="ml-auto flex items-center gap-2">
          {camOn != null && (
            <span
              title={camOn ? "Camera on" : "Camera off"}
              className={`flex h-7 w-7 items-center justify-center rounded-full ${
                camOn ? "bg-teal-light text-teal" : "bg-gray-100 text-gray-400"
              }`}
            >
              <Camera size={15} />
            </span>
          )}
          {screenOn != null && (
            <span
              title={screenOn ? "Screen shared" : "Screen not shared"}
              className={`flex h-7 w-7 items-center justify-center rounded-full ${
                screenOn ? "bg-teal-light text-teal" : "bg-gray-100 text-gray-400"
              }`}
            >
              <MonitorUp size={15} />
            </span>
          )}
          {timer != null && (
            <span
              className={`flex items-center gap-1.5 rounded-full px-3 py-1 text-body-sm font-extrabold tabular-nums ${
                low ? "bg-coral text-white" : "bg-purple-light text-purple"
              }`}
            >
              <Clock size={15} />
              {fmtTime(timer)}
            </span>
          )}
        </div>
      </header>
      <div className="flex flex-1 flex-col px-5 py-5 lg:px-8">{children}</div>
    </div>
  );
}

function Centered({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex flex-1 flex-col items-center justify-center text-center">
      {children}
    </div>
  );
}

function SectionHeader({
  icon,
  name,
  tint,
}: {
  icon: React.ReactNode;
  name: string;
  tint: string;
}) {
  return (
    <div className="mb-4 flex items-center gap-2">
      <span
        className="flex h-9 w-9 items-center justify-center rounded-lg"
        style={{ backgroundColor: tint + "1A", color: tint }}
      >
        {icon}
      </span>
      <span className="text-label-lg font-bold uppercase tracking-wide text-gray-500">
        {name}
      </span>
    </div>
  );
}

/* ---------- level select ---------- */

function LevelSelect({
  lang,
  completed,
  onStart,
}: {
  lang: string;
  completed: string[];
  onStart: (level: string) => void;
}) {
  const firstOpen = LEVELS.find((l) => !completed.includes(l.code))?.code || "A1";
  const [sel, setSel] = useState(firstOpen);
  const done = completed.includes(sel);

  return (
    <div className="flex flex-1 flex-col">
      <div className="text-center">
        <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-purple text-white">
          <GraduationCap size={28} />
        </span>
        <h2 className="mt-3 text-heading-xl font-extrabold text-ink">
          {languageName(lang)} Certification
        </h2>
        <p className="mt-1 text-body-md text-slatey">
          Choose a level to attempt — the higher the level, the harder the exam.
        </p>
      </div>

      <p className="mb-2 mt-5 text-label-lg font-bold uppercase tracking-wide text-gray-500">
        Select your level
      </p>
      <div className="grid grid-cols-2 gap-2 sm:grid-cols-3">
        {LEVELS.map((l, i) => {
          const isDone = completed.includes(l.code);
          const active = sel === l.code;
          return (
            <button
              key={l.code}
              onClick={() => setSel(l.code)}
              className={`relative rounded-xl border-2 p-3 text-left transition ${
                active
                  ? "border-purple bg-purple-light"
                  : "border-gray-100 bg-white"
              }`}
            >
              <div className="flex items-center justify-between">
                <span className="text-heading-sm font-extrabold text-ink">
                  {l.code}
                </span>
                {isDone && <Check size={16} className="text-teal" />}
              </div>
              <span className="block text-label-md text-slatey">{l.name}</span>
              <div className="mt-1.5 flex gap-0.5">
                {Array.from({ length: 6 }).map((_, d) => (
                  <span
                    key={d}
                    className={`h-1.5 w-1.5 rounded-full ${
                      d <= i ? "bg-amber" : "bg-gray-200"
                    }`}
                  />
                ))}
              </div>
            </button>
          );
        })}
      </div>

      <p className="mt-4 rounded-full bg-amber-light px-3 py-1 text-center text-label-md font-bold text-amber">
        Normally a paid exam · free during testing
      </p>

      <div className="mt-auto pt-6">
        {done ? (
          <div className="rounded-xl bg-teal-light p-3 text-center text-body-sm font-semibold text-teal">
            ✓ You&apos;ve already taken the {sel} exam — pick another level.
          </div>
        ) : (
          <Button full onClick={() => onStart(sel)}>
            Continue to {sel} rules
          </Button>
        )}
      </div>
    </div>
  );
}

/* ---------- rules ---------- */

function RulesScreen({
  lang,
  level,
  durationSeconds,
  passMark,
  weights,
  starting,
  onBack,
  onBegin,
}: {
  lang: string;
  level: string;
  durationSeconds: number;
  passMark: number;
  weights: { listening: number; reading: number; writing: number; speaking: number };
  starting: boolean;
  onBack: () => void;
  onBegin: () => void;
}) {
  const [agree, setAgree] = useState(false);
  const mins = Math.round(durationSeconds / 60);

  const rules: { icon: React.ReactNode; title: string; body: string }[] = [
    {
      icon: <Clock size={18} />,
      title: `${mins}-minute timer`,
      body: `You have ${fmtTime(durationSeconds)} to complete all four sections. When the timer reaches zero the exam is submitted automatically.`,
    },
    {
      icon: <AlertTriangle size={18} />,
      title: "Stay on this tab",
      body: "If you switch to another tab or app, or minimise the window, the exam ends immediately and is scored on what you've completed so far.",
    },
    {
      icon: <Camera size={18} />,
      title: "Camera stays on",
      body: "You'll be asked to allow your camera. It stays on for the duration of the exam for integrity.",
    },
    {
      icon: <MonitorUp size={18} />,
      title: "Screen sharing",
      body: "You'll be asked to share your screen when the exam begins. This is part of the proctoring.",
    },
    {
      icon: <ShieldCheck size={18} />,
      title: `Weighted sections`,
      body: `Listening ${weights.listening}% · Reading ${weights.reading}% · Writing ${weights.writing}% · Speaking ${weights.speaking}%. Sections run back-to-back — you can't return to a previous one.`,
    },
    {
      icon: <GraduationCap size={18} />,
      title: `Pass mark: ${passMark}%`,
      body: `Score at least ${passMark}% (weighted) to earn your ${level} certificate. Each level can only be passed once — you can't retake a level you've certified.`,
    },
  ];

  return (
    <div className="flex flex-1 flex-col">
      <div className="text-center">
        <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-purple text-white">
          <ShieldCheck size={28} />
        </span>
        <h2 className="mt-3 text-heading-xl font-extrabold text-ink">
          Before you start — {level} exam
        </h2>
        <p className="mt-1 text-body-md text-slatey">
          This is a proctored {languageName(lang)} exam. Please read the rules.
        </p>
      </div>

      <div className="mt-5 space-y-2.5">
        {rules.map((r) => (
          <div
            key={r.title}
            className="flex gap-3 rounded-2xl border border-gray-100 bg-white p-3.5 shadow-card"
          >
            <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-purple-light text-purple">
              {r.icon}
            </span>
            <div>
              <p className="text-body-md font-extrabold text-ink">{r.title}</p>
              <p className="mt-0.5 text-body-sm text-slatey">{r.body}</p>
            </div>
          </div>
        ))}
      </div>

      <label className="mt-5 flex cursor-pointer items-start gap-3 rounded-2xl bg-purple-light p-4">
        <input
          type="checkbox"
          checked={agree}
          onChange={(e) => setAgree(e.target.checked)}
          className="mt-0.5 h-5 w-5 accent-purple"
        />
        <span className="text-body-sm font-semibold text-ink">
          I understand the rules and agree to be proctored. I know that leaving
          this tab or running out of time will end my exam.
        </span>
      </label>

      <div className="mt-auto flex gap-3 pt-5">
        <Button variant="outline" onClick={onBack} className="flex-1">
          Back
        </Button>
        <Button
          full
          disabled={!agree || starting}
          onClick={onBegin}
          className="flex-[2]"
        >
          {starting ? "Preparing…" : `Begin ${level} exam`}
        </Button>
      </div>
    </div>
  );
}

/* ---------- MCQ-based sections (listening + reading) ---------- */

function McqSection({
  header,
  intro,
  questions,
  onDone,
}: {
  header: React.ReactNode;
  intro: React.ReactNode;
  questions: { question: string; options: string[] | null; correctAnswer: string }[];
  onDone: (score: number) => void;
}) {
  const [answers, setAnswers] = useState<Record<number, string>>({});
  const allAnswered = questions.every((_, i) => answers[i]);

  function finish() {
    const correct = questions.filter(
      (q, i) => answers[i] === q.correctAnswer
    ).length;
    onDone(questions.length ? Math.round((correct / questions.length) * 100) : 0);
  }

  return (
    <div className="flex flex-1 flex-col">
      {header}
      {intro}
      <div className="mt-4 space-y-4">
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
                    className={`flex h-12 w-full items-center rounded-md border-2 px-4 text-left text-body-md font-semibold transition ${
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
        <Button full disabled={!allAnswered} onClick={finish}>
          Submit section
        </Button>
      </div>
    </div>
  );
}

function ListeningSection({
  session,
  onDone,
}: {
  session: ListeningSession;
  onDone: (score: number) => void;
}) {
  const lines = session.lines || [];
  return (
    <McqSection
      header={
        <SectionHeader
          icon={<Headphones size={18} />}
          name="Listening"
          tint="#00C2A8"
        />
      }
      intro={
        <button
          onClick={() =>
            speakSequence(
              lines.map((l) => ({ character: l.character, text: l.text }))
            )
          }
          className="flex items-center justify-center gap-2 rounded-full bg-purple py-3 font-extrabold text-white shadow-float"
        >
          <Volume2 size={20} /> Play the conversation
        </button>
      }
      questions={session.questions || []}
      onDone={onDone}
    />
  );
}

function ReadingSection({
  session,
  onDone,
}: {
  session: ReadingSession;
  onDone: (score: number) => void;
}) {
  const lines = session.lines || [];
  return (
    <McqSection
      header={
        <SectionHeader icon={<BookText size={18} />} name="Reading" tint="#17A3DD" />
      }
      intro={
        <div className="rounded-2xl border border-gray-100 bg-white p-4 shadow-card">
          {lines.map((l) => (
            <p
              key={l.id}
              className="text-body-lg font-semibold leading-relaxed text-ink"
            >
              {l.text}
            </p>
          ))}
        </div>
      }
      questions={session.questions || []}
      onDone={onDone}
    />
  );
}

/* ---------- writing ---------- */

function WritingSection({
  lang,
  minWords,
  prompt,
  onDone,
}: {
  lang: string;
  minWords: number;
  prompt: string;
  onDone: (score: number) => void;
}) {
  const [text, setText] = useState("");
  const words = useMemo(
    () => text.trim().split(/\s+/).filter(Boolean).length,
    [text]
  );

  function finish() {
    const score = Math.min(100, Math.round((words / minWords) * 100));
    onDone(score);
  }

  return (
    <div className="flex flex-1 flex-col">
      <SectionHeader icon={<PenLine size={18} />} name="Writing" tint="#F5A623" />
      <p className="text-body-md text-ink">
        Write in <strong>{languageName(lang)}</strong>, {prompt}
      </p>
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        rows={7}
        placeholder="Type your answer here…"
        className="mt-4 w-full rounded-xl border border-gray-100 bg-white p-4 text-body-lg outline-none transition focus:border-purple"
      />
      <p className="mt-2 text-body-sm text-slatey">
        {words} words{" "}
        {words < minWords ? `· aim for at least ${minWords}` : "· nice!"}
      </p>
      <div className="mt-auto pt-6">
        <Button full disabled={words < 3} onClick={finish}>
          Submit section
        </Button>
      </div>
    </div>
  );
}

/* ---------- speaking ---------- */

function SpeakingSection({
  phrase,
  speaker,
  locale,
  onDone,
}: {
  phrase: string;
  speaker: string;
  locale: string;
  onDone: (score: number) => void;
}) {
  const [listening, setListening] = useState(false);
  const [score, setScore] = useState<number | null>(null);
  const [supported] = useState(() => speechRecognitionSupported());

  async function record() {
    if (listening) return;
    setListening(true);
    setScore(null);
    try {
      const said = await recognizeSpeech(locale);
      setScore(scorePronunciation(phrase, said));
    } catch {
      setScore(null);
    } finally {
      setListening(false);
    }
  }

  return (
    <div className="flex flex-1 flex-col">
      <SectionHeader icon={<Mic size={18} />} name="Speaking" tint="#FF5C5C" />
      <div className="flex flex-1 flex-col items-center justify-center text-center">
        <p className="text-body-md text-slatey">Read this aloud clearly:</p>
        <p className="mt-2 text-heading-lg font-extrabold text-purple">
          &ldquo;{phrase}&rdquo;
        </p>
        <div className="mt-5 flex items-center gap-3">
          <button
            onClick={() => speakAs(speaker, phrase)}
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
            ? "Speech scoring needs Chrome/Edge — you can skip."
            : listening
            ? "Listening…"
            : "Tap the mic and read the sentence"}
        </p>
        {score != null && (
          <p className="mt-4 text-display-lg font-extrabold text-teal">{score}%</p>
        )}
      </div>
      <div className="mt-auto pt-6">
        {score == null ? (
          <Button full variant="outline" onClick={() => onDone(0)}>
            Skip &amp; finish
          </Button>
        ) : (
          <Button full onClick={() => onDone(score)}>
            Finish exam
          </Button>
        )}
      </div>
    </div>
  );
}

/* ---------- result ---------- */

function ResultView({
  result,
  submitting,
  lang,
  termination,
  onRetake,
}: {
  result: ExamResult | null;
  submitting: boolean;
  lang: string;
  termination: Termination;
  onRetake: () => void;
}) {
  if (submitting || !result) {
    return (
      <Centered>
        <FoxMascot size={110} glow />
        <p className="mt-4 text-body-md text-slatey">Marking your exam…</p>
      </Centered>
    );
  }

  const banner =
    termination === "tab" ? (
      <div className="mb-4 flex items-center gap-2 rounded-xl bg-coral/10 px-4 py-3 text-body-sm font-semibold text-coral">
        <AlertTriangle size={18} />
        Your exam ended because you left the tab. It was scored on what you
        completed.
      </div>
    ) : termination === "time" ? (
      <div className="mb-4 flex items-center gap-2 rounded-xl bg-amber-light px-4 py-3 text-body-sm font-semibold text-amber">
        <Clock size={18} />
        Time&apos;s up — your exam was submitted automatically.
      </div>
    ) : null;

  if (!result.passed) {
    return (
      <Centered>
        {banner}
        <span className="text-5xl">📚</span>
        <h2 className="mt-3 text-heading-xl font-extrabold text-ink">
          Almost there!
        </h2>
        <p className="mt-1 text-body-md text-slatey">
          You scored {result.overall}% (weighted). You need {result.passMark}% to
          certify — keep practising and try again.
        </p>
        <SectionScores result={result} />
        <div className="mt-6 w-full max-w-sm">
          <Button full onClick={onRetake}>
            Retake exam
          </Button>
        </div>
      </Centered>
    );
  }

  return (
    <Centered>
      {banner}
      <FoxMascot size={120} glow bounce />
      <h2 className="mt-3 text-heading-xl font-extrabold text-ink">
        {result.alreadyTaken ? "Already certified ✓" : "You passed! 🎉"}
      </h2>
      <p className="mt-1 text-body-md text-slatey">
        {result.alreadyTaken
          ? `You've already earned the ${result.level} certificate for ${languageName(lang)}.`
          : `${languageName(lang)} · Level `}
        {!result.alreadyTaken && (
          <>
            <strong className="text-purple">{result.level}</strong> ·{" "}
            {result.overall}%
          </>
        )}
      </p>
      <SectionScores result={result} />
      <div className="mt-6 w-full max-w-sm space-y-2">
        {result.certificate && (
          <Link href={`/certificates/${result.certificate.id}`}>
            <Button full>View your certificate</Button>
          </Link>
        )}
        <Link href="/profile">
          <Button full variant="outline">
            Back to profile
          </Button>
        </Link>
      </div>
    </Centered>
  );
}

function SectionScores({ result }: { result: ExamResult }) {
  const s = result.sections;
  if (!s) return null;
  const w = result.weights;
  const items = [
    ["Listening", s.listening, w?.listening],
    ["Reading", s.reading, w?.reading],
    ["Writing", s.writing, w?.writing],
    ["Speaking", s.speaking, w?.speaking],
  ] as const;
  return (
    <div className="mt-5 grid w-full max-w-sm grid-cols-2 gap-2">
      {items.map(([name, val, weight]) => (
        <div key={name} className="rounded-xl bg-white p-3 text-left shadow-card">
          <p className="text-label-md text-slatey">
            {name}
            {weight != null && (
              <span className="ml-1 text-label-sm text-gray-400">
                · {weight}%
              </span>
            )}
          </p>
          <p className="flex items-center gap-1 text-heading-sm font-extrabold text-ink">
            {val}% {val >= 60 && <Check size={14} className="text-teal" />}
          </p>
        </div>
      ))}
    </div>
  );
}
