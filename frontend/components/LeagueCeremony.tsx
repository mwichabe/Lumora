"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import {
  AnimatePresence,
  motion,
  useReducedMotion,
  type Variants,
} from "framer-motion";
import { Crown, Flame, Gem, Sparkles, Target, Trophy, X } from "lucide-react";
import { Avatar } from "./Avatar";
import { Button } from "./Button";
import type { LeagueOutcome, LeagueResult } from "@/lib/types";

/**
 * The end-of-race ceremony: a full-screen, five-beat replay of the week that
 * just closed. It plays exactly once per season — the backend hands out a
 * result only while its `ceremonySeen` flag is false, and we set that as soon
 * as the curtain goes up.
 *
 * The beats are driven off one `scene` counter rather than nested delays, so
 * the whole sequence can be skipped (tap anywhere) or fast-forwarded without
 * leaving half-finished animations behind. Every beat is also legible as a
 * static frame, which is what someone with reduced motion actually gets.
 */

type Scene = 0 | 1 | 2 | 3 | 4;

// How long each beat holds before the next begins (ms).
const BEATS: Record<Scene, number> = { 0: 1400, 1: 2200, 2: 2600, 3: 2400, 4: 0 };

interface Tone {
  verb: string;
  headline: string;
  sub: string;
  emoji: string;
  accent: string;
  /** Which way the tier badge travels during the verdict beat. */
  travel: "up" | "down" | "hold";
  celebratory: boolean;
}

function toneFor(r: LeagueResult): Tone {
  if (r.flagged) {
    return {
      verb: "Withheld",
      headline: "Result withheld",
      sub: r.flagReason || "This week's activity was flagged for review.",
      emoji: "🛡️",
      accent: "#9090A0",
      travel: "hold",
      celebratory: false,
    };
  }
  const map: Record<LeagueOutcome, Tone> = {
    promoted: {
      verb: "Promoted",
      headline: `Welcome to ${r.to.name}`,
      sub: "You finished in the promotion zone. A tougher pod is waiting.",
      emoji: "🚀",
      accent: r.to.tint,
      travel: "up",
      celebratory: true,
    },
    held: {
      verb: "Held",
      headline: `You held ${r.from.name}`,
      sub: "No movement this week — the ladder is still yours to climb.",
      emoji: "🏅",
      accent: r.from.tint,
      travel: "hold",
      celebratory: false,
    },
    demoted: {
      verb: "Moved down",
      headline: `Down to ${r.to.name}`,
      sub: "A fresh pod, a clean slate. One lesson a day is usually enough to climb back.",
      emoji: "🌧️",
      accent: r.to.tint,
      travel: "down",
      celebratory: false,
    },
    qualified: {
      verb: "Qualified",
      headline: "Diamond Tournament",
      sub: "Top ten in Diamond. Three weeks, three rounds, one trophy.",
      emoji: "💠",
      accent: "#17A3DD",
      travel: "up",
      celebratory: true,
    },
    advanced: {
      verb: "Advanced",
      headline: "Through to the next round",
      sub: "You survived the cut. One more to go.",
      emoji: "⚔️",
      accent: "#17A3DD",
      travel: "up",
      celebratory: true,
    },
    champion: {
      verb: "Champion",
      headline: "Tournament champion",
      sub: "A podium finish in the final. The trophy is yours for good.",
      emoji: "🏆",
      accent: "#F5A623",
      travel: "up",
      celebratory: true,
    },
    eliminated: {
      verb: "Run over",
      headline: "Back to Diamond",
      sub: "The bracket ends here — and opens again next season.",
      emoji: "🛡️",
      accent: "#17A3DD",
      travel: "hold",
      celebratory: false,
    },
  };
  return map[r.result] ?? map.held;
}

export function LeagueCeremony({
  result,
  onClose,
}: {
  result: LeagueResult;
  onClose: () => void;
}) {
  const reduced = useReducedMotion();
  const [scene, setScene] = useState<Scene>(0);
  const tone = useMemo(() => toneFor(result), [result]);

  // Advance through the beats on a timer. Reduced motion jumps straight to the
  // summary rather than sitting through a sequence it can't see move.
  useEffect(() => {
    if (reduced) {
      setScene(4);
      return;
    }
    if (scene >= 4) return;
    const t = setTimeout(() => setScene((s) => (s + 1) as Scene), BEATS[scene]);
    return () => clearTimeout(t);
  }, [scene, reduced]);

  const skip = () => setScene(4);

  return (
    <motion.div
      className="fixed inset-0 z-[80] flex items-center justify-center overflow-y-auto px-5 py-8"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      onClick={scene < 4 ? skip : undefined}
      role="dialog"
      aria-modal="true"
      aria-label="League results"
    >
      <Backdrop accent={tone.accent} />
      {tone.celebratory && !reduced && scene >= 2 && (
        <Confetti accent={tone.accent} />
      )}

      <div className="relative z-10 w-full max-w-md text-center text-white">
        <button
          onClick={onClose}
          aria-label="Close results"
          className="absolute -top-2 right-0 rounded-full bg-white/10 p-2 text-white/70 transition hover:bg-white/20 hover:text-white"
        >
          <X size={18} />
        </button>

        <AnimatePresence mode="wait">
          {scene === 0 && <SceneOpen key="s0" result={result} />}
          {scene === 1 && <SceneRank key="s1" result={result} accent={tone.accent} />}
          {scene === 2 && <SceneVerdict key="s2" result={result} tone={tone} />}
          {scene === 3 && <SceneRewards key="s3" result={result} tone={tone} />}
          {scene === 4 && (
            <SceneSummary key="s4" result={result} tone={tone} onClose={onClose} />
          )}
        </AnimatePresence>

        {scene < 4 && (
          <motion.p
            className="mt-8 text-label-md text-white/40"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.8 }}
          >
            Tap anywhere to skip
          </motion.p>
        )}
      </div>
    </motion.div>
  );
}

// --- beat 1: the curtain -----------------------------------------------------

function SceneOpen({ result }: { result: LeagueResult }) {
  return (
    <motion.div
      variants={beat}
      initial="enter"
      animate="center"
      exit="exit"
      className="flex flex-col items-center"
    >
      <motion.div
        initial={{ scale: 0, rotate: -30 }}
        animate={{ scale: 1, rotate: 0 }}
        transition={{ type: "spring", stiffness: 220, damping: 14 }}
        className="relative"
      >
        <Shockwave tint={result.from.tint} />
        <TierBadge tier={result.from} size={104} />
      </motion.div>

      <motion.p
        className="mt-7 text-label-lg uppercase tracking-[0.3em] text-white/50"
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.35 }}
      >
        Season {result.seasonId}
      </motion.p>
      <motion.h2
        className="mt-2 text-display-lg font-extrabold"
        initial={{ opacity: 0, y: 14 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5, type: "spring", stiffness: 200, damping: 18 }}
      >
        The race is over
      </motion.h2>
      <motion.p
        className="mt-2 text-body-md text-white/60"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.75 }}
      >
        {result.from.name} League · {result.podSize} learners
      </motion.p>
    </motion.div>
  );
}

// --- beat 2: where you landed ------------------------------------------------

function SceneRank({ result, accent }: { result: LeagueResult; accent: string }) {
  // The rank counts *down* from last place to where you actually finished — a
  // climb reads as a climb, and the number lands on the beat.
  const rank = useCountTo(result.rank, 1200, result.podSize);
  const pct = 1 - (result.rank - 1) / Math.max(result.podSize - 1, 1);

  return (
    <motion.div
      variants={beat}
      initial="enter"
      animate="center"
      exit="exit"
      className="flex flex-col items-center"
    >
      <p className="text-label-lg uppercase tracking-[0.3em] text-white/50">
        You finished
      </p>

      <div className="relative mt-5 flex h-44 w-44 items-center justify-center">
        <ProgressRing progress={pct} tint={accent} />
        <div className="relative">
          <motion.span
            key={rank}
            initial={{ scale: 1.15, opacity: 0.6 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ duration: 0.12 }}
            className="text-[64px] font-extrabold leading-none"
          >
            #{rank}
          </motion.span>
        </div>
      </div>

      <motion.p
        className="mt-4 text-body-lg text-white/70"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.1 }}
      >
        of {result.podSize} in your pod
      </motion.p>

      <motion.div
        className="mt-6 flex items-end justify-center gap-3"
        initial="hidden"
        animate="show"
        variants={{ show: { transition: { staggerChildren: 0.12, delayChildren: 1.2 } } }}
      >
        {[2, 1, 3].map((place) => {
          const p = result.podium.find((x) => x.rank === place);
          if (!p) return null;
          const height = place === 1 ? 74 : place === 2 ? 56 : 42;
          return (
            <motion.div
              key={place}
              variants={{
                hidden: { opacity: 0, y: 18 },
                show: { opacity: 1, y: 0 },
              }}
              className="flex w-[84px] flex-col items-center"
            >
              <Avatar
                name={p.name}
                color={p.avatarColor}
                url={p.avatarUrl}
                size={place === 1 ? 44 : 36}
              />
              <p className="mt-1 w-full truncate text-label-md font-bold text-white/80">
                {p.isUser ? "You" : p.name}
              </p>
              <motion.div
                initial={{ height: 0 }}
                animate={{ height }}
                transition={{ delay: 1.4, type: "spring", stiffness: 160, damping: 20 }}
                className="mt-1 w-full rounded-t-lg bg-white/15"
                style={{
                  background:
                    place === 1
                      ? "linear-gradient(180deg,#F5A623,#F5A62333)"
                      : "rgba(255,255,255,0.14)",
                }}
              >
                <span className="block pt-2 text-label-md font-extrabold">
                  {place}
                </span>
              </motion.div>
            </motion.div>
          );
        })}
      </motion.div>
    </motion.div>
  );
}

// --- beat 3: the verdict -----------------------------------------------------

function SceneVerdict({ result, tone }: { result: LeagueResult; tone: Tone }) {
  const travel = tone.travel;
  const moving = travel !== "hold";
  const dy = travel === "up" ? -46 : 46;

  return (
    <motion.div
      variants={beat}
      initial="enter"
      animate="center"
      exit="exit"
      className="flex flex-col items-center"
    >
      {/* The badge physically travels between tiers: it dips, launches along a
          rail of light in the direction of travel, and lands as the new tier. */}
      <div className="relative flex h-[168px] w-full items-center justify-center">
        {travel !== "hold" && <TravelRail direction={travel} tint={tone.accent} />}

        <motion.div
          className="relative"
          initial={{ y: moving ? -dy * 0.6 : 0, scale: 0.9, opacity: 0 }}
          animate={{ y: 0, scale: 1, opacity: 1 }}
          transition={{ type: "spring", stiffness: 180, damping: 15, delay: 0.15 }}
        >
          <AnimatePresence mode="wait">
            <motion.div
              key={result.to.name}
              initial={{ rotateY: 90, opacity: 0 }}
              animate={{ rotateY: 0, opacity: 1 }}
              transition={{ delay: 0.55, duration: 0.45 }}
            >
              <Shockwave tint={tone.accent} delay={0.55} />
              <TierBadge tier={result.to} size={112} />
            </motion.div>
          </AnimatePresence>
        </motion.div>
      </div>

      <motion.p
        className="text-label-lg uppercase tracking-[0.3em]"
        style={{ color: tone.accent }}
        initial={{ opacity: 0, letterSpacing: "0.6em" }}
        animate={{ opacity: 1, letterSpacing: "0.3em" }}
        transition={{ delay: 0.75, duration: 0.5 }}
      >
        {tone.verb}
      </motion.p>
      <motion.h2
        className="mt-2 text-display-lg font-extrabold"
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.9, type: "spring", stiffness: 200, damping: 18 }}
      >
        {tone.headline}
      </motion.h2>
      <motion.p
        className="mt-3 max-w-xs text-body-md text-white/60"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 1.15 }}
      >
        {tone.sub}
      </motion.p>
    </motion.div>
  );
}

// --- beat 4: the chest -------------------------------------------------------

function SceneRewards({ result, tone }: { result: LeagueResult; tone: Tone }) {
  const gems = useCountTo(result.gems, 900);
  const chest =
    result.rank === 1 ? "Gold" : result.rank === 2 ? "Silver" : result.rank === 3 ? "Bronze" : "";

  return (
    <motion.div
      variants={beat}
      initial="enter"
      animate="center"
      exit="exit"
      className="flex flex-col items-center"
    >
      <div className="relative flex h-40 items-center justify-center">
        <Chest tint={tone.accent} empty={result.gems === 0} />
      </div>

      {result.gems > 0 ? (
        <>
          <motion.p
            className="text-label-lg uppercase tracking-[0.3em] text-white/50"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
          >
            {chest ? `${chest} chest` : "Reward"}
          </motion.p>
          <motion.div
            className="mt-2 flex items-center justify-center gap-2 text-display-xl font-extrabold"
            initial={{ scale: 0.6, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ delay: 0.7, type: "spring", stiffness: 240, damping: 14 }}
          >
            <Gem size={30} className="text-teal" />
            {gems}
          </motion.div>
        </>
      ) : (
        <motion.p
          className="mt-2 text-heading-lg font-extrabold"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          No chest this week
        </motion.p>
      )}

      {result.groupGoalHit && (
        <motion.div
          className="mt-5 flex items-center gap-2 rounded-full bg-teal/20 px-4 py-2 text-label-lg font-bold text-teal"
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 1, type: "spring", stiffness: 260, damping: 16 }}
        >
          <Target size={16} /> Pod goal smashed — bonus gems included
        </motion.div>
      )}
    </motion.div>
  );
}

// --- beat 5: the scorecard ---------------------------------------------------

function SceneSummary({
  result,
  tone,
  onClose,
}: {
  result: LeagueResult;
  tone: Tone;
  onClose: () => void;
}) {
  return (
    <motion.div
      variants={beat}
      initial="enter"
      animate="center"
      exit="exit"
      className="flex flex-col items-center"
    >
      <TierBadge tier={result.to} size={84} />
      <p
        className="mt-4 text-label-lg uppercase tracking-[0.3em]"
        style={{ color: tone.accent }}
      >
        {tone.verb}
      </p>
      <h2 className="mt-1 text-heading-xl font-extrabold">{tone.headline}</h2>
      <p className="mt-1 text-body-md text-white/60">
        #{result.rank} of {result.podSize} · Season {result.seasonId}
      </p>

      <div className="mt-6 grid w-full grid-cols-2 gap-3">
        <Stat label="League points" value={result.points} icon={<Sparkles size={14} />} />
        <Stat label="XP earned" value={result.rawXp} icon={<Flame size={14} />} />
        <Stat label="Avg accuracy" value={`${result.accuracy}%`} icon={<Target size={14} />} />
        <Stat label="Flawless runs" value={result.perfectRuns} icon={<Crown size={14} />} />
      </div>

      {result.gems > 0 && (
        <div className="mt-3 flex w-full items-center justify-center gap-2 rounded-2xl bg-white/10 py-3 text-heading-sm font-extrabold">
          <Gem size={18} className="text-teal" /> +{result.gems} gems
        </div>
      )}

      {result.trophies > 0 && (
        <div className="mt-3 flex items-center gap-2 text-label-lg font-bold text-amber">
          <Trophy size={16} /> {result.trophies} tournament{" "}
          {result.trophies === 1 ? "trophy" : "trophies"}
        </div>
      )}

      {result.flagged && (
        <p className="mt-4 rounded-2xl bg-white/10 px-4 py-3 text-body-sm text-white/70">
          Promotion and rewards were held back this week because the activity was
          flagged ({result.flagReason}). Your XP and course progress are untouched.
        </p>
      )}

      <Button className="mt-7" full onClick={onClose}>
        Start this week&apos;s race
      </Button>
    </motion.div>
  );
}

function Stat({
  label,
  value,
  icon,
}: {
  label: string;
  value: number | string;
  icon: React.ReactNode;
}) {
  return (
    <div className="rounded-2xl bg-white/10 px-3 py-3 text-center">
      <div className="flex items-center justify-center gap-1 text-label-md text-white/50">
        {icon} {label}
      </div>
      <div className="mt-1 text-heading-md font-extrabold">{value}</div>
    </div>
  );
}

// --- scenery -----------------------------------------------------------------

const beat: Variants = {
  enter: { opacity: 0, y: 24, scale: 0.96 },
  center: { opacity: 1, y: 0, scale: 1, transition: { duration: 0.4 } },
  exit: { opacity: 0, y: -20, scale: 0.97, transition: { duration: 0.28 } },
};

/** A slowly drifting spotlight over deep space, tinted by the outcome. */
function Backdrop({ accent }: { accent: string }) {
  return (
    <div className="absolute inset-0 bg-space">
      <motion.div
        className="absolute left-1/2 top-1/2 h-[130vmax] w-[130vmax] -translate-x-1/2 -translate-y-1/2 rounded-full"
        style={{
          background: `radial-gradient(circle, ${accent}44 0%, ${accent}11 35%, transparent 68%)`,
        }}
        animate={{ scale: [1, 1.12, 1], opacity: [0.75, 1, 0.75] }}
        transition={{ duration: 7, repeat: Infinity, ease: "easeInOut" }}
      />
      <div className="absolute inset-0 bg-space/40" />
    </div>
  );
}

/** An expanding ring, fired once when a badge lands. */
function Shockwave({ tint, delay = 0 }: { tint: string; delay?: number }) {
  return (
    <>
      {[0, 0.18].map((offset) => (
        <motion.span
          key={offset}
          className="pointer-events-none absolute left-1/2 top-1/2 h-24 w-24 -translate-x-1/2 -translate-y-1/2 rounded-full border-2"
          style={{ borderColor: tint }}
          initial={{ scale: 0.4, opacity: 0.8 }}
          animate={{ scale: 3.2, opacity: 0 }}
          transition={{ duration: 1.1, delay: delay + offset, ease: "easeOut" }}
        />
      ))}
    </>
  );
}

/** The beam a promoted (or demoted) badge rides between tiers. */
function TravelRail({
  direction,
  tint,
}: {
  direction: "up" | "down";
  tint: string;
}) {
  const up = direction === "up";
  return (
    <>
      <motion.div
        className="pointer-events-none absolute left-1/2 h-full w-[3px] -translate-x-1/2 rounded-full"
        style={{
          background: `linear-gradient(${up ? "0deg" : "180deg"}, transparent, ${tint}, transparent)`,
        }}
        initial={{ scaleY: 0, opacity: 0 }}
        animate={{ scaleY: 1, opacity: [0, 1, 0.25] }}
        transition={{ duration: 0.8, delay: 0.1 }}
      />
      {[...Array(8)].map((_, i) => (
        <motion.span
          key={i}
          className="pointer-events-none absolute left-1/2 h-1.5 w-1.5 rounded-full"
          style={{ background: tint, marginLeft: (i % 2 ? 1 : -1) * (8 + i * 3) }}
          initial={{ y: up ? 70 : -70, opacity: 0 }}
          animate={{ y: up ? -80 : 80, opacity: [0, 1, 0] }}
          transition={{
            duration: 0.9,
            delay: 0.15 + i * 0.06,
            ease: "easeOut",
          }}
        />
      ))}
    </>
  );
}

/** The tier crest: a rounded gem plate in the tier's colour. */
function TierBadge({
  tier,
  size = 96,
}: {
  tier: { name: string; tint: string };
  size?: number;
}) {
  return (
    <div
      className="relative flex items-center justify-center rounded-[28%] shadow-float"
      style={{
        width: size,
        height: size,
        background: `linear-gradient(150deg, ${tier.tint}, ${tier.tint}99)`,
        boxShadow: `0 12px 40px ${tier.tint}66`,
      }}
    >
      <div className="absolute inset-[10%] rounded-[26%] border-2 border-white/25" />
      <Crown size={size * 0.4} className="relative text-white drop-shadow" />
      <span className="absolute -bottom-7 w-max text-label-lg font-extrabold tracking-wide text-white/90">
        {tier.name}
      </span>
    </div>
  );
}

/** Sweeps a ring to show how far up the pod the user finished. */
function ProgressRing({ progress, tint }: { progress: number; tint: string }) {
  const r = 78;
  const c = 2 * Math.PI * r;
  return (
    <svg className="absolute inset-0 -rotate-90" viewBox="0 0 176 176">
      <circle cx="88" cy="88" r={r} fill="none" stroke="rgba(255,255,255,0.12)" strokeWidth="8" />
      <motion.circle
        cx="88"
        cy="88"
        r={r}
        fill="none"
        stroke={tint}
        strokeWidth="8"
        strokeLinecap="round"
        strokeDasharray={c}
        initial={{ strokeDashoffset: c }}
        animate={{ strokeDashoffset: c * (1 - Math.max(progress, 0.02)) }}
        transition={{ duration: 1.2, ease: "easeOut" }}
      />
    </svg>
  );
}

/** A chest that pops its lid and throws gems when there's something inside. */
function Chest({ tint, empty }: { tint: string; empty: boolean }) {
  return (
    <div className="relative">
      {!empty &&
        [...Array(10)].map((_, i) => {
          const angle = -140 + i * 16;
          const dist = 70 + (i % 3) * 22;
          return (
            <motion.span
              key={i}
              className="absolute left-1/2 top-1/2 h-2.5 w-2.5 rounded-sm"
              style={{ background: i % 2 ? "#00C2A8" : "#F5A623" }}
              initial={{ x: 0, y: 0, opacity: 0, scale: 0.4 }}
              animate={{
                x: Math.cos((angle * Math.PI) / 180) * dist,
                y: Math.sin((angle * Math.PI) / 180) * dist,
                opacity: [0, 1, 0],
                scale: [0.4, 1, 0.6],
                rotate: 240,
              }}
              transition={{ duration: 1.1, delay: 0.42 + i * 0.03, ease: "easeOut" }}
            />
          );
        })}

      <motion.div
        initial={{ scale: 0.6, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ type: "spring", stiffness: 240, damping: 15 }}
        className="relative"
      >
        {/* lid */}
        <motion.div
          className="mx-auto h-7 w-28 origin-bottom rounded-t-2xl"
          style={{ background: `linear-gradient(180deg, ${tint}, ${tint}aa)` }}
          initial={{ rotateX: 0 }}
          animate={empty ? { rotateX: 0 } : { rotateX: -105 }}
          transition={{ delay: 0.4, type: "spring", stiffness: 160, damping: 12 }}
        />
        {/* body */}
        <div
          className="mx-auto h-16 w-28 rounded-b-2xl"
          style={{ background: `linear-gradient(180deg, ${tint}cc, ${tint}66)` }}
        >
          <div className="mx-auto h-full w-5 bg-white/25" />
        </div>
      </motion.div>
    </div>
  );
}

/** Hand-rolled confetti — no extra dependency, and it respects the tier tint. */
function Confetti({ accent }: { accent: string }) {
  const pieces = useMemo(
    () =>
      [...Array(64)].map((_, i) => ({
        id: i,
        left: Math.random() * 100,
        delay: Math.random() * 1.2,
        duration: 2.6 + Math.random() * 2,
        drift: (Math.random() - 0.5) * 160,
        size: 6 + Math.random() * 7,
        rotate: Math.random() * 720 - 360,
        color: [accent, "#F5A623", "#00C2A8", "#FF5C5C", "#FFFFFF"][i % 5],
        round: i % 3 === 0,
      })),
    [accent]
  );

  return (
    <div className="pointer-events-none absolute inset-0 overflow-hidden">
      {pieces.map((p) => (
        <motion.span
          key={p.id}
          className={`absolute top-[-8%] ${p.round ? "rounded-full" : "rounded-[2px]"}`}
          style={{
            left: `${p.left}%`,
            width: p.size,
            height: p.round ? p.size : p.size * 0.5,
            background: p.color,
          }}
          initial={{ y: 0, opacity: 0, rotate: 0 }}
          animate={{
            y: "110vh",
            x: p.drift,
            rotate: p.rotate,
            opacity: [0, 1, 1, 0],
          }}
          transition={{
            duration: p.duration,
            delay: p.delay,
            ease: "linear",
            repeat: Infinity,
            repeatDelay: 0.6,
          }}
        />
      ))}
    </div>
  );
}

/**
 * Counts to `to`, optionally starting somewhere else — the rank counter runs
 * backwards from last place so finishing high feels like an ascent.
 */
function useCountTo(to: number, duration = 1000, from = 0) {
  const [value, setValue] = useState(from);
  const frame = useRef<number>();

  useEffect(() => {
    const start = performance.now();
    const tick = (now: number) => {
      const t = Math.min((now - start) / duration, 1);
      // easeOutCubic, so the number decelerates onto its final value.
      const eased = 1 - Math.pow(1 - t, 3);
      setValue(Math.round(from + (to - from) * eased));
      if (t < 1) frame.current = requestAnimationFrame(tick);
    };
    frame.current = requestAnimationFrame(tick);
    return () => {
      if (frame.current) cancelAnimationFrame(frame.current);
    };
  }, [to, duration, from]);

  return value;
}
