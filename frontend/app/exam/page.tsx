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
  Award,
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
import { SpeakerChip } from "@/components/Speaker";
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
  ExamResult,
  ExamPaper,
  PaperQuestion,
  PaymentStatus,
} from "@/lib/types";

type Phase =
  | "intro"
  | "pay"
  | "rules"
  | "listening"
  | "reading"
  | "writing"
  | "speaking"
  | "result";

type Termination = "tab" | "time" | "screen" | "camera" | null;

const LEVELS = [
  { code: "A1", name: "Beginner" },
  { code: "A2", name: "Elementary" },
  { code: "B1", name: "Intermediate" },
  { code: "B2", name: "Upper-Int." },
  { code: "C1", name: "Advanced" },
  { code: "C2", name: "Mastery" },
];
// The comprehensive A1→C2 exam is offered separately from the per-level ladder.
const FINAL_CODE = "FINAL";
const FALLBACK_DURATION = [600, 780, 960, 1140, 1320, 1500];
const FALLBACK_PASS = [50, 58, 65, 72, 80, 88];
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
  const [completed, setCompleted] = useState<string[]>([]);
  const [paper, setPaper] = useState<ExamPaper | null>(null);
  const [paperLoading, setPaperLoading] = useState(false);
  const [payStatus, setPayStatus] = useState<PaymentStatus | null>(null);
  const [unlocking, setUnlocking] = useState(false);

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
  // Blocks entry to the exam until proctoring requirements are met.
  const [setupError, setSetupError] = useState<string | null>(null);

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

  const fbIdx = Math.min(idx, 5);
  const duration = paper?.durationSeconds ?? FALLBACK_DURATION[fbIdx];
  const passMark = paper?.passMark ?? FALLBACK_PASS[fbIdx];
  const weights = paper?.weights ?? {
    listening: 30,
    reading: 30,
    writing: 20,
    speaking: 20,
  };

  const paymentsOn = !!payStatus?.paymentsEnabled;
  const priceFor = (lvl: string) => payStatus?.prices?.[lvl] ?? 0;
  const priceUsdFor = (lvl: string) => payStatus?.pricesUsd?.[lvl] ?? 0;
  const paidFor = (lvl: string) => !!payStatus?.paid?.[lvl];

  /* ---------- load completed levels + payment status ---------- */
  useEffect(() => {
    Promise.all([
      api.certificates().catch(() => ({ certificates: [] })),
      api.paymentStatus().catch(() => null),
    ])
      .then(([cs, ps]) => {
        setCompleted(
          cs.certificates.filter((c) => c.language === lang).map((c) => c.level)
        );
        setPayStatus(ps);
      })
      .finally(() => setLoading(false));
    return () => stopSpeaking();
  }, [lang]);

  // Start Paystack checkout for the currently-selected level.
  async function payForLevel() {
    setUnlocking(true);
    try {
      const r = await api.initializePayment(level);
      if (r.authorizationUrl) {
        window.location.href = r.authorizationUrl; // to Paystack's checkout
        return;
      }
      setUnlocking(false);
    } catch {
      setUnlocking(false);
    }
  }

  // Pick a level: load its paper, and route to payment first if it isn't paid.
  const selectLevel = useCallback(
    (lvl: string) => {
      setLevel(lvl);
      setScores({ listening: 0, reading: 0, writing: 0, speaking: 0 });
      setSetupError(null);
      setPaper(null);
      setPaperLoading(true);
      const mustPay = !!payStatus?.paymentsEnabled && !payStatus?.paid?.[lvl];
      setPhase(mustPay ? "pay" : "rules");
      api
        .examPaper(lvl)
        .then(setPaper)
        .catch(() => setPaper(null))
        .finally(() => setPaperLoading(false));
    },
    [payStatus]
  );

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
        // Surface the pass/fail notification right away.
        window.dispatchEvent(new Event("lumora:notifications"));
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

  /* ---------- begin (require camera + screen, then start) ---------- */
  const beginExam = useCallback(async () => {
    setStarting(true);
    setSetupError(null);
    endedRef.current = false;
    setTermination(null);

    const md =
      typeof navigator !== "undefined" ? navigator.mediaDevices : undefined;

    // 0) Environment must support proctoring at all (needs a secure context).
    if (!md || !md.getDisplayMedia || !md.getUserMedia) {
      stopProctoring();
      setScreenOn(false);
      setCamOn(false);
      setStarting(false);
      setSetupError(
        "Your browser can't be proctored here. Use Chrome or Edge over https (or localhost) so screen sharing and camera are available."
      );
      return;
    }

    // 1) Screen sharing is required.
    try {
      const scr = await md.getDisplayMedia({ video: true });
      if (!scr.getVideoTracks().length) throw new Error("no-video");
      screenStreamRef.current = scr;
      setScreenOn(true);
    } catch {
      stopProctoring();
      setScreenOn(false);
      setCamOn(false);
      setStarting(false);
      setSetupError(
        "Screen sharing is required to take this exam. When prompted, choose a screen to share and allow it — then try again."
      );
      return;
    }

    // 2) Camera is required.
    try {
      const cam = await md.getUserMedia({ video: true, audio: false });
      if (!cam.getVideoTracks().length) throw new Error("no-video");
      camStreamRef.current = cam;
      setCamOn(true);
    } catch {
      // Roll back the screen share too — we won't start without both.
      stopProctoring();
      setScreenOn(false);
      setCamOn(false);
      setStarting(false);
      setSetupError(
        "Camera access is required to take this exam. Allow your camera when prompted (and check no other app is using it) — then try again."
      );
      return;
    }

    // 3) If either feed is stopped mid-exam, end the exam immediately.
    screenStreamRef.current
      ?.getTracks()
      .forEach((t) =>
        t.addEventListener("ended", () => {
          if (activeRef.current && !endedRef.current) endExam("screen");
        })
      );
    camStreamRef.current
      ?.getTracks()
      .forEach((t) =>
        t.addEventListener("ended", () => {
          if (activeRef.current && !endedRef.current) endExam("camera");
        })
      );

    // Notify the user that the attempt has begun (fire-and-forget).
    api
      .startExam(level, lang)
      .then(() => window.dispatchEvent(new Event("lumora:notifications")))
      .catch(() => {});

    setSecondsLeft(duration);
    setStarting(false);
    setPhase("listening");
  }, [duration, stopProctoring, endExam, level, lang]);

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
    setSetupError(null);
    setResult(null);
    setSecondsLeft(0);
    setScores({ listening: 0, reading: 0, writing: 0, speaking: 0 });
    setPhase("intro");
    // Refresh payment status — the attempt just taken has been consumed, so a
    // retake will correctly require paying again.
    api.paymentStatus().then(setPayStatus).catch(() => {});
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

  const locale = LOCALE[lang] || "en-US";
  const notReady = !paperLoading && (!paper || !paper.ready);

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
      ) : phase === "intro" ? (
        <LevelSelect
          lang={lang}
          completed={completed}
          paymentsOn={paymentsOn}
          currency={payStatus?.currency ?? "KES"}
          priceFor={priceFor}
          priceUsdFor={priceUsdFor}
          paidFor={paidFor}
          onStart={selectLevel}
        />
      ) : phase === "pay" ? (
        <UnlockScreen
          lang={lang}
          level={level}
          isFinal={level === FINAL_CODE}
          price={priceFor(level)}
          priceUsd={priceUsdFor(level)}
          currency={payStatus?.currency ?? "KES"}
          unlocking={unlocking}
          onUnlock={payForLevel}
          onBack={() => setPhase("intro")}
        />
      ) : phase === "rules" ? (
        paperLoading ? (
          <Centered>
            <FoxMascot size={110} glow />
            <p className="mt-4 text-body-md text-slatey">
              Preparing your {level} paper…
            </p>
          </Centered>
        ) : notReady ? (
          <Centered>
            <FoxMascot size={110} glow />
            <p className="mt-4 text-heading-sm font-extrabold text-ink">
              Exam not ready for {languageName(lang)}
            </p>
            <p className="mt-1 max-w-xs text-body-md text-slatey">
              This language needs a full course (lessons, listening &amp;
              reading) before the {level} exam is available.
            </p>
            <div className="mt-6 w-full max-w-sm space-y-2">
              <Button full variant="outline" onClick={() => setPhase("intro")}>
                Choose another level
              </Button>
              <Button full variant="ghost" onClick={() => router.push("/learn")}>
                Go to lessons
              </Button>
            </div>
          </Centered>
        ) : (
          <RulesScreen
            lang={lang}
            level={level}
            durationSeconds={duration}
            passMark={passMark}
            weights={weights}
            starting={starting}
            error={setupError}
            onBack={() => {
              setSetupError(null);
              setPhase("intro");
            }}
            onBegin={beginExam}
          />
        )
      ) : phase === "listening" ? (
        <ListeningSection
          listening={paper!.listening!}
          onDone={(s) => advance("listening", s, "reading")}
        />
      ) : phase === "reading" ? (
        <ReadingSection
          reading={paper!.reading!}
          onDone={(s) => advance("reading", s, "writing")}
        />
      ) : phase === "writing" ? (
        <WritingSection
          lang={lang}
          minWords={paper!.writing.minWords}
          prompt={paper!.writing.prompt}
          onDone={(s) => advance("writing", s, "speaking")}
        />
      ) : phase === "speaking" ? (
        <SpeakingSection
          phrase={paper!.speaking.phrase}
          speaker={paper!.speaking.speaker || "Lumora"}
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

/* ---------- unlock (payment) ---------- */

function UnlockScreen({
  lang,
  level,
  isFinal,
  price,
  priceUsd,
  currency,
  unlocking,
  onUnlock,
  onBack,
}: {
  lang: string;
  level: string;
  isFinal: boolean;
  price: number;
  priceUsd: number;
  currency: string;
  unlocking: boolean;
  onUnlock: () => void;
  onBack: () => void;
}) {
  const perks = isFinal
    ? [
        "One comprehensive exam covering A1 → C2",
        "The longest, most demanding paper",
        "A prestige Mastery certificate",
        "Proctored — camera + screen share",
      ]
    : [
        `A proctored, weighted 4-skill ${level} exam`,
        "A verifiable, downloadable certificate",
        "Stamped & sealed by Lumora",
        "Retake anytime (pay per attempt)",
      ];
  return (
    <div className="flex flex-1 flex-col">
      <div className="text-center">
        <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-purple text-white">
          <GraduationCap size={28} />
        </span>
        <h2 className="mt-3 text-heading-xl font-extrabold text-ink">
          {isFinal
            ? `${languageName(lang)} Mastery Exam`
            : `Unlock the ${level} exam`}
        </h2>
        <p className="mt-1 text-body-md text-slatey">
          {isFinal
            ? "The ultimate test — everything from A1 to C2 in one sitting."
            : `Pay once per attempt. Higher levels cost more.`}
        </p>
      </div>

      <div className="mt-5 rounded-2xl border border-gray-100 bg-white p-5 shadow-card">
        <div className="flex items-baseline justify-center gap-1">
          <span className="text-body-md font-bold text-slatey">{currency}</span>
          <span className="text-display-lg font-extrabold text-ink">{price}</span>
          <span className="text-body-md text-slatey">· one attempt</span>
        </div>
        {priceUsd > 0 && (
          <p className="mt-0.5 text-center text-body-sm text-slatey">
            ≈ ${priceUsd.toFixed(2)} USD
          </p>
        )}
        <ul className="mt-4 space-y-2">
          {perks.map((p) => (
            <li key={p} className="flex items-center gap-2 text-body-md text-ink">
              <Check size={16} className="shrink-0 text-teal" />
              {p}
            </li>
          ))}
        </ul>
      </div>

      <p className="mt-3 flex items-center justify-center gap-1.5 text-label-md text-slatey">
        <ShieldCheck size={14} /> Secure checkout by Paystack
      </p>

      <div className="mt-auto space-y-2 pt-6">
        <Button full disabled={unlocking} onClick={onUnlock}>
          {unlocking ? "Redirecting…" : `Pay ${currency} ${price}`}
        </Button>
        <Button full variant="ghost" onClick={onBack}>
          Choose another level
        </Button>
      </div>
    </div>
  );
}

/* ---------- level select ---------- */

function LevelSelect({
  lang,
  completed,
  paymentsOn,
  currency,
  priceFor,
  priceUsdFor,
  paidFor,
  onStart,
}: {
  lang: string;
  completed: string[];
  paymentsOn: boolean;
  currency: string;
  priceFor: (level: string) => number;
  priceUsdFor: (level: string) => number;
  paidFor: (level: string) => boolean;
  onStart: (level: string) => void;
}) {
  const firstOpen = LEVELS.find((l) => !completed.includes(l.code))?.code || "A1";
  const [sel, setSel] = useState(firstOpen);
  const isDone = completed.includes(sel);
  const paid = paidFor(sel);
  const usd = priceUsdFor(sel);
  const ctaLabel =
    !paymentsOn || paid
      ? "Start"
      : `Pay ${currency} ${priceFor(sel)}${usd > 0 ? ` (≈ $${usd.toFixed(2)})` : ""} & start`;

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
          Choose a level — the higher the level, the harder the exam
          {paymentsOn ? " and the higher the price." : "."}
        </p>
      </div>

      <p className="mb-2 mt-5 text-label-lg font-bold uppercase tracking-wide text-gray-500">
        Select your level
      </p>
      <div className="grid grid-cols-2 gap-2 sm:grid-cols-3">
        {LEVELS.map((l, i) => {
          const done = completed.includes(l.code);
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
                {done && <Check size={16} className="text-teal" />}
              </div>
              <span className="block text-label-md text-slatey">{l.name}</span>
              {paymentsOn && (
                <span className="mt-1 block text-label-sm font-bold text-purple">
                  {paidFor(l.code)
                    ? "Paid ✓"
                    : `${currency} ${priceFor(l.code)}${
                        priceUsdFor(l.code) > 0
                          ? ` · $${priceUsdFor(l.code).toFixed(2)}`
                          : ""
                      }`}
                </span>
              )}
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

      {/* Comprehensive final exam */}
      <button
        onClick={() => setSel(FINAL_CODE)}
        className={`mt-3 flex items-center gap-3 rounded-2xl border-2 p-4 text-left transition ${
          sel === FINAL_CODE
            ? "border-purple bg-purple-light"
            : "border-gray-100 bg-white"
        }`}
      >
        <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-amber text-white">
          <Award size={22} />
        </span>
        <div className="min-w-0 flex-1">
          <p className="font-extrabold text-ink">Final Mastery Exam</p>
          <p className="text-body-sm text-slatey">
            Comprehensive A1 → C2 · the ultimate certificate
          </p>
        </div>
        {paymentsOn && (
          <span className="shrink-0 text-right text-label-lg font-bold text-purple">
            {paidFor(FINAL_CODE) ? (
              "Paid ✓"
            ) : (
              <>
                {currency} {priceFor(FINAL_CODE)}
                {priceUsdFor(FINAL_CODE) > 0 && (
                  <span className="block text-label-sm font-semibold text-slatey">
                    ≈ ${priceUsdFor(FINAL_CODE).toFixed(2)}
                  </span>
                )}
              </>
            )}
          </span>
        )}
      </button>

      <div className="mt-auto pt-6">
        {isDone && (
          <p className="mb-2 text-center text-body-sm font-semibold text-teal">
            ✓ You&apos;re already certified at {sel} — retaking issues a fresh
            certificate.
          </p>
        )}
        <Button full onClick={() => onStart(sel)}>
          {ctaLabel}
        </Button>
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
  error,
  onBack,
  onBegin,
}: {
  lang: string;
  level: string;
  durationSeconds: number;
  passMark: number;
  weights: { listening: number; reading: number; writing: number; speaking: number };
  starting: boolean;
  error: string | null;
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
      title: "Camera required",
      body: "You must allow your camera. It stays on for the whole exam — if you deny it or it turns off, the exam won't start (or ends immediately).",
    },
    {
      icon: <MonitorUp size={18} />,
      title: "Screen sharing required",
      body: "You must share your screen when the exam begins. If you deny it or stop sharing, the exam won't start (or ends immediately).",
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

      {error && (
        <div className="mt-4 flex items-start gap-2 rounded-2xl border border-coral/30 bg-coral/10 p-3.5 text-coral">
          <AlertTriangle size={18} className="mt-0.5 shrink-0" />
          <p className="text-body-sm font-semibold">{error}</p>
        </div>
      )}

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
          {starting
            ? "Preparing…"
            : error
            ? "Allow & try again"
            : `Begin ${level} exam`}
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
  questions: PaperQuestion[];
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
  listening,
  onDone,
}: {
  listening: NonNullable<ExamPaper["listening"]>;
  onDone: (score: number) => void;
}) {
  const lines = listening.lines || [];
  const [played, setPlayed] = useState(false);
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
        <div>
          <p className="mb-2 text-body-md font-extrabold text-ink">
            {listening.title}
          </p>
          <div className="mb-3 flex flex-wrap gap-2">
            {Array.from(new Set(lines.map((l) => l.character))).map((sp) => (
              <SpeakerChip key={sp} name={sp} />
            ))}
          </div>
          <button
            onClick={() => {
              setPlayed(true);
              speakSequence(
                lines.map((l) => ({ character: l.character, text: l.text }))
              );
            }}
            className="flex w-full items-center justify-center gap-2 rounded-full bg-purple py-3 font-extrabold text-white shadow-float"
          >
            <Volume2 size={20} /> {played ? "Play again" : "Play the recording"}
          </button>
          <p className="mt-2 text-center text-body-sm text-slatey">
            You can replay the recording while you answer.
          </p>
        </div>
      }
      questions={listening.questions || []}
      onDone={onDone}
    />
  );
}

function ReadingSection({
  reading,
  onDone,
}: {
  reading: NonNullable<ExamPaper["reading"]>;
  onDone: (score: number) => void;
}) {
  const paragraphs = reading.paragraphs || [];
  return (
    <McqSection
      header={
        <SectionHeader icon={<BookText size={18} />} name="Reading" tint="#17A3DD" />
      }
      intro={
        <div className="space-y-2 rounded-2xl border border-gray-100 bg-white p-4 shadow-card">
          <p className="text-body-md font-extrabold text-ink">{reading.title}</p>
          {paragraphs.map((p, i) => (
            <p
              key={i}
              className="text-body-lg leading-relaxed text-ink/90"
            >
              {p}
            </p>
          ))}
        </div>
      }
      questions={reading.questions || []}
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
        <SpeakerChip name={speaker} className="mb-4" />
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
    ) : termination === "screen" ? (
      <div className="mb-4 flex items-center gap-2 rounded-xl bg-coral/10 px-4 py-3 text-body-sm font-semibold text-coral">
        <MonitorUp size={18} />
        Your exam ended because screen sharing stopped. It was scored on what you
        completed.
      </div>
    ) : termination === "camera" ? (
      <div className="mb-4 flex items-center gap-2 rounded-xl bg-coral/10 px-4 py-3 text-body-sm font-semibold text-coral">
        <Camera size={18} />
        Your exam ended because the camera turned off. It was scored on what you
        completed.
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
        <div className="mt-6 w-full max-w-sm space-y-2">
          <Button full onClick={onRetake}>
            Retake exam
          </Button>
          <Link href="/profile">
            <Button full variant="outline">
              Back to profile
            </Button>
          </Link>
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
