"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { motion, useReducedMotion } from "framer-motion";
import { FoxMascot } from "@/components/FoxMascot";
import { useAuth } from "@/lib/auth";
import { getLastRoute } from "@/lib/api";

const APP_ROUTES = ["/home", "/learn", "/practice", "/leaderboard", "/profile", "/lesson"];

export default function SplashScreen() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const reduceMotion = useReducedMotion();

  useEffect(() => {
    if (loading) return;

    let target: string;
    if (user) {
      if (!user.targetLanguage) {
        target = "/onboarding/language";
      } else {
        // Restore the exact screen the user was last on (e.g. after they
        // trimmed the URL back to "/"), falling back to home.
        const last = getLastRoute();
        target =
          last && APP_ROUTES.some((r) => last.startsWith(r)) ? last : "/home";
      }
    } else {
      target = "/onboarding/welcome";
    }

    // Authenticated users skip the long intro — just restore them quickly.
    const delay = user ? 400 : 1900;
    const t = setTimeout(() => router.replace(target), delay);
    return () => clearTimeout(t);
  }, [loading, user, router]);

  return (
    <main
      className="relative flex min-h-[100dvh] w-full flex-col items-center justify-center overflow-hidden px-6 text-white"
      style={{
        background:
          "radial-gradient(120% 120% at 50% 0%, #7B4AD6 0%, #4A2A9E 38%, #1F1640 72%, #0F0F24 100%)",
      }}
    >
      {/* Starfield + aurora ambience — full-bleed, scales to any viewport */}
      <Backdrop reduceMotion={!!reduceMotion} />

      {/* Centred hero column — fluid sizing keeps it balanced on phones and 4K alike */}
      <div className="relative z-10 flex w-full max-w-2xl flex-col items-center text-center">
        <motion.div
          initial={reduceMotion ? false : { scale: 0.85, opacity: 0, y: 24 }}
          animate={{ scale: 1, opacity: 1, y: 0 }}
          transition={{ type: "spring", stiffness: 260, damping: 18 }}
          className="relative"
        >
          {/* halo behind the mascot */}
          <div
            className="pointer-events-none absolute left-1/2 top-1/2 -z-10 -translate-x-1/2 -translate-y-1/2 rounded-full bg-amber/30 blur-3xl"
            style={{ width: "120%", height: "120%" }}
          />
          <motion.div
            animate={reduceMotion ? undefined : { y: [0, -12, 0] }}
            transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
            className="w-[clamp(160px,34vw,300px)] [&_svg]:!h-full [&_svg]:!w-full [&>div]:!h-full [&>div]:!w-full"
            style={{ aspectRatio: "1 / 1" }}
          >
            <FoxMascot size={300} glow bounce={!reduceMotion} />
          </motion.div>
        </motion.div>

        <motion.h1
          initial={reduceMotion ? false : { opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.35 }}
          className="mt-2 bg-gradient-to-b from-white to-amber/90 bg-clip-text font-extrabold tracking-tight text-transparent"
          style={{ fontSize: "clamp(2.75rem, 9vw, 5rem)", letterSpacing: "-2px" }}
        >
          LUMORA
        </motion.h1>

        <motion.p
          initial={reduceMotion ? false : { opacity: 0 }}
          animate={{ opacity: 0.85 }}
          transition={{ delay: 0.55 }}
          className="mt-1 max-w-md px-4 italic text-white/80"
          style={{ fontSize: "clamp(0.95rem, 2.5vw, 1.25rem)" }}
        >
          Learn a language. Fall in love with it.
        </motion.p>

        {/* Trust strip — gives the first screen a little marketing weight on larger displays */}
        <motion.div
          initial={reduceMotion ? false : { opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.75 }}
          className="mt-7 hidden items-center gap-5 text-white/70 sm:flex"
          style={{ fontSize: "clamp(0.75rem, 1.6vw, 0.95rem)" }}
        >
          <Stat value="40+" label="languages" />
          <span className="h-4 w-px bg-white/20" />
          <Stat value="5 min" label="a day" />
          <span className="h-4 w-px bg-white/20" />
          <Stat value="Free" label="to start" />
        </motion.div>
      </div>

      {/* Loading indicator pinned to the bottom, responsive width */}
      <div className="absolute bottom-[max(2.5rem,env(safe-area-inset-bottom))] left-1/2 flex -translate-x-1/2 flex-col items-center gap-3">
        <div className="h-1 w-[min(60vw,14rem)] overflow-hidden rounded-full bg-white/20">
          <motion.div
            className="h-full rounded-full bg-gradient-to-r from-amber to-teal"
            initial={{ width: "0%" }}
            animate={{ width: "100%" }}
            transition={{ duration: 2, ease: "easeInOut" }}
          />
        </div>
        <span className="text-xs font-semibold tracking-widest text-white/50">
          LOADING
        </span>
      </div>
    </main>
  );
}

function Stat({ value, label }: { value: string; label: string }) {
  return (
    <span className="flex items-baseline gap-1.5">
      <strong className="font-extrabold text-white">{value}</strong>
      <span className="text-white/60">{label}</span>
    </span>
  );
}

/** Decorative, non-interactive background: drifting aurora blobs + twinkling stars. */
function Backdrop({ reduceMotion }: { reduceMotion: boolean }) {
  const blobs = [
    { className: "left-[-10%] top-[-8%] bg-teal/30", size: "38vmax", delay: 0 },
    { className: "right-[-12%] top-[10%] bg-purple/40", size: "42vmax", delay: 2 },
    { className: "bottom-[-15%] left-[20%] bg-amber/25", size: "34vmax", delay: 4 },
  ];

  return (
    <div className="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden>
      {blobs.map((b, i) => (
        <motion.div
          key={i}
          className={`absolute rounded-full blur-[80px] ${b.className}`}
          style={{ width: b.size, height: b.size }}
          animate={
            reduceMotion
              ? undefined
              : { x: [0, 30, -20, 0], y: [0, -25, 20, 0], scale: [1, 1.1, 0.95, 1] }
          }
          transition={{ duration: 16, repeat: Infinity, ease: "easeInOut", delay: b.delay }}
        />
      ))}

      {/* twinkling stars */}
      {STARS.map((s, i) => (
        <motion.span
          key={`s-${i}`}
          className="absolute rounded-full bg-white"
          style={{ left: s.left, top: s.top, width: s.r, height: s.r }}
          animate={reduceMotion ? undefined : { opacity: [0.2, 1, 0.2] }}
          transition={{ duration: s.dur, repeat: Infinity, ease: "easeInOut", delay: s.delay }}
        />
      ))}
    </div>
  );
}

// Deterministic star positions (no Math.random → no hydration mismatch).
const STARS = [
  { left: "12%", top: "18%", r: 2, dur: 3, delay: 0 },
  { left: "82%", top: "22%", r: 3, dur: 4, delay: 0.5 },
  { left: "68%", top: "12%", r: 2, dur: 3.5, delay: 1 },
  { left: "28%", top: "72%", r: 2, dur: 4.5, delay: 0.2 },
  { left: "90%", top: "62%", r: 2.5, dur: 3.2, delay: 1.4 },
  { left: "8%", top: "55%", r: 2, dur: 4, delay: 0.8 },
  { left: "45%", top: "8%", r: 1.5, dur: 3.8, delay: 1.2 },
  { left: "55%", top: "88%", r: 2, dur: 3.4, delay: 0.6 },
  { left: "20%", top: "38%", r: 1.5, dur: 4.2, delay: 1.6 },
  { left: "76%", top: "78%", r: 2, dur: 3.6, delay: 0.3 },
];
