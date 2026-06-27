"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";

interface LessonResult {
  xp: number;
  accuracy: number;
  firstClear: boolean;
}

// A small burst of confetti pieces rendered with framer-motion.
function Confetti() {
  const pieces = useMemo(() => {
    const colors = ["#6C3FC5", "#F5A623", "#00C2A8", "#FF5C5C", "#3A1F8A"];
    return Array.from({ length: 36 }).map((_, i) => ({
      id: i,
      left: Math.random() * 100,
      delay: Math.random() * 0.5,
      duration: 1.8 + Math.random() * 1.4,
      color: colors[i % colors.length],
      size: 7 + Math.random() * 8,
      rotate: Math.random() * 360,
    }));
  }, []);

  return (
    <div className="pointer-events-none absolute inset-0 overflow-hidden">
      {pieces.map((p) => (
        <motion.div
          key={p.id}
          initial={{ y: -40, opacity: 0, rotate: 0 }}
          animate={{ y: "110vh", opacity: [0, 1, 1, 0], rotate: p.rotate }}
          transition={{ duration: p.duration, delay: p.delay, ease: "easeIn" }}
          style={{
            position: "absolute",
            left: `${p.left}%`,
            width: p.size,
            height: p.size * 0.6,
            backgroundColor: p.color,
            borderRadius: 2,
          }}
        />
      ))}
    </div>
  );
}

export default function LessonCompletePage() {
  const router = useRouter();
  const [result, setResult] = useState<LessonResult | null>(null);

  useEffect(() => {
    const raw = sessionStorage.getItem("lumora_lesson_result");
    if (raw) {
      try {
        setResult(JSON.parse(raw));
      } catch {
        setResult({ xp: 10, accuracy: 100, firstClear: true });
      }
    } else {
      setResult({ xp: 10, accuracy: 100, firstClear: true });
    }
  }, []);

  const xp = result?.xp ?? 0;
  const accuracy = result?.accuracy ?? 0;

  function onContinue() {
    sessionStorage.removeItem("lumora_lesson_result");
    router.replace("/home");
  }

  return (
    <div className="relative flex min-h-[100dvh] w-full flex-col items-center overflow-hidden bg-gradient-to-b from-purple-dark via-purple to-purple-dark px-6 pb-10 pt-20 text-white">
      <Confetti />

      <div className="z-10 flex w-full max-w-md flex-1 flex-col items-center justify-between">
      <div className="flex flex-1 flex-col items-center justify-center text-center">
        <motion.div
          initial={{ scale: 0.4, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ type: "spring", stiffness: 200, damping: 14 }}
        >
          <FoxMascot size={150} glow bounce />
        </motion.div>

        <motion.h1
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.25 }}
          className="mt-6 text-heading-xl font-extrabold"
        >
          Lesson Complete!
        </motion.h1>
        <motion.p
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.35 }}
          className="mt-2 text-body-lg text-purple-light"
        >
          {result?.firstClear
            ? "Brilliant work — Lumora is proud of you!"
            : "Nicely done — practice makes fluent!"}
        </motion.p>

        <motion.div
          initial={{ y: 24, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="mt-8 flex w-full gap-4"
        >
          <StatCard label="XP Earned" value={`+${xp}`} accent="#F5A623" />
          <StatCard label="Accuracy" value={`${accuracy}%`} accent="#00C2A8" />
        </motion.div>
      </div>

      <motion.div
        initial={{ y: 30, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.7 }}
        className="z-10 w-full"
      >
        <Button full variant="secondary" onClick={onContinue}>
          Continue
        </Button>
      </motion.div>
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
  accent,
}: {
  label: string;
  value: string;
  accent: string;
}) {
  return (
    <div className="flex-1 rounded-2xl bg-white/10 px-4 py-5 backdrop-blur">
      <div className="text-display-lg font-extrabold" style={{ color: accent }}>
        {value}
      </div>
      <div className="mt-1 text-label-sm uppercase tracking-wide text-purple-light">
        {label}
      </div>
    </div>
  );
}
