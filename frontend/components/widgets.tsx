"use client";

import { motion } from "framer-motion";
import { Heart, Flame, Gem } from "lucide-react";

// XP progress bar — amber gradient fill, animated on change (spec 5.5).
export function XPBar({ value, max }: { value: number; max: number }) {
  const pct = Math.min(100, max > 0 ? (value / max) * 100 : 0);
  return (
    <div className="h-3 w-full overflow-hidden rounded-full bg-gray-100">
      <motion.div
        className="h-full rounded-full"
        style={{ background: "linear-gradient(90deg,#F5A623,#FF9A00)" }}
        initial={{ width: 0 }}
        animate={{ width: `${pct}%` }}
        transition={{ type: "spring", stiffness: 120, damping: 20 }}
      />
    </div>
  );
}

// Row of 5 hearts; full hearts coral, empty hearts outlined (spec 5.5).
export function HeartIndicator({ hearts }: { hearts: number }) {
  return (
    <div className="flex items-center gap-1">
      {Array.from({ length: 5 }).map((_, i) => (
        <Heart
          key={i}
          size={20}
          className={i < hearts ? "text-coral" : "text-coral/30"}
          fill={i < hearts ? "#FF5C5C" : "transparent"}
          strokeWidth={2}
        />
      ))}
    </div>
  );
}

// Streak flame with day count.
export function StreakFlame({
  streak,
  light = false,
}: {
  streak: number;
  light?: boolean;
}) {
  return (
    <div className="flex items-center gap-1">
      <Flame
        size={22}
        className={streak > 0 ? "text-amber" : light ? "text-white/50" : "text-gray-300"}
        fill={streak > 0 ? "#F5A623" : "transparent"}
        strokeWidth={2}
      />
      <span className={`font-extrabold ${light ? "text-white" : "text-ink"}`}>{streak}</span>
    </div>
  );
}

// Gem counter pill (teal bg).
export function GemCounter({ gems }: { gems: number }) {
  return (
    <div className="flex items-center gap-1 rounded-full bg-teal px-3 py-1 text-white">
      <Gem size={16} strokeWidth={2.5} />
      <span className="text-label-md font-bold">{gems}</span>
    </div>
  );
}

// Character speech bubble (spec 5.6).
export function SpeechBubble({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`relative max-w-[260px] rounded-[28px] border border-gray-100 bg-white px-4 py-3 text-body-md text-ink shadow-card ${className}`}
    >
      {children}
      <span className="absolute -bottom-2 left-8 h-4 w-4 rotate-45 border-b border-r border-gray-100 bg-white" />
    </div>
  );
}
