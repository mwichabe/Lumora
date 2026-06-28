"use client";

import { motion } from "framer-motion";
import { Lock, Check } from "lucide-react";
import { SkillIcon } from "@/components/SkillIcon";
import type { Skill } from "@/lib/types";

/**
 * A visual "journey" of the course — skill nodes on a winding path, grouped by
 * unit, coloured by status (completed / current / locked).
 */
export function RoadmapView({
  skills,
  userXp,
  onOpen,
}: {
  skills: Skill[];
  userXp: number;
  onOpen: (lessonId: number) => void;
}) {
  // Group by unit, preserving backend order.
  const units: { name: string; skills: Skill[] }[] = [];
  skills.forEach((s) => {
    const name = s.unit || "Course";
    let u = units.find((x) => x.name === name);
    if (!u) {
      u = { name, skills: [] };
      units.push(u);
    }
    u.skills.push(s);
  });

  let nodeIndex = -1;

  return (
    <div className="mx-auto max-w-md">
      {units.map((u) => (
        <section key={u.name} className="mb-4">
          <div className="my-3 flex items-center gap-3">
            <span className="h-px flex-1 bg-gray-200" />
            <span className="rounded-full bg-white px-3 py-1 text-label-md font-bold uppercase tracking-wide text-purple shadow-card">
              {u.name}
            </span>
            <span className="h-px flex-1 bg-gray-200" />
          </div>

          <div className="relative flex flex-col items-center gap-6">
            <div className="pointer-events-none absolute bottom-6 left-1/2 top-6 w-px -translate-x-1/2 border-l-2 border-dashed border-gray-300" />
            {u.skills.map((s) => {
              nodeIndex++;
              const offset = nodeIndex % 2 === 0 ? "-56px" : "56px";
              const locked = !s.unlocked;
              const completed = s.completed;
              return (
                <div
                  key={s.id}
                  className="relative flex flex-col items-center"
                  style={{ transform: `translateX(${offset})` }}
                >
                  <motion.button
                    whileTap={{ scale: locked ? 1 : 0.94 }}
                    onClick={() => {
                      if (!locked && s.lessons?.length) onOpen(s.lessons[0].id);
                    }}
                    aria-label={s.title}
                    className="relative flex h-20 w-20 items-center justify-center rounded-full shadow-skill"
                    style={{
                      background: locked ? "#E5E5EC" : s.color,
                      color: locked ? "#9090A0" : "#fff",
                    }}
                  >
                    {locked ? (
                      <Lock size={26} />
                    ) : (
                      <SkillIcon name={s.icon} size={28} className="text-white" />
                    )}
                    {!locked && !completed && (
                      <span className="absolute inset-0 animate-pulseRing rounded-full ring-4 ring-amber/70" />
                    )}
                    {completed && (
                      <span className="absolute -right-1 -top-1 flex h-7 w-7 items-center justify-center rounded-full border-2 border-cream bg-teal text-white">
                        <Check size={14} />
                      </span>
                    )}
                  </motion.button>
                  <p className="mt-2 max-w-[140px] text-center text-body-sm font-bold text-ink">
                    {s.title}
                  </p>
                  {locked ? (
                    <p className="text-label-md text-amber">
                      {Math.max(0, s.requiredXp - userXp)} XP to unlock
                    </p>
                  ) : (
                    <p className="text-label-md text-slatey">
                      {s.completedCount}/{s.lessonCount} lessons
                    </p>
                  )}
                </div>
              );
            })}
          </div>
        </section>
      ))}

      <div className="mt-6 flex flex-col items-center">
        <span className="text-3xl">🏁</span>
        <p className="text-body-sm font-bold text-slatey">Fluency</p>
      </div>
    </div>
  );
}
