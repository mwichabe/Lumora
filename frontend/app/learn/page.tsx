"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Lock, Check, Play, ChevronRight, Headphones } from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { SkillIcon } from "@/components/SkillIcon";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type { Skill, ListeningSession } from "@/lib/types";

const LANGUAGE_NAMES: Record<string, string> = {
  es: "Spanish",
  fr: "French",
  ja: "Japanese",
  zh: "Mandarin",
  ar: "Arabic",
  sw: "Swahili",
  pt: "Portuguese",
  de: "German",
  it: "Italian",
  ko: "Korean",
  hi: "Hindi",
};

export default function LearnPage() {
  return (
    <AppShell wide>
      <CourseScreen />
    </AppShell>
  );
}

function CourseScreen() {
  const { user } = useAuth();
  const router = useRouter();
  const [skills, setSkills] = useState<Skill[] | null>(null);
  const [sessions, setSessions] = useState<ListeningSession[]>([]);
  const [error, setError] = useState(false);

  const load = useCallback(() => {
    setError(false);
    api
      .skills()
      .then((d) => setSkills(d.skills))
      .catch(() => setError(true));
    api
      .listeningSessions()
      .then((d) => setSessions(d.sessions))
      .catch(() => setSessions([]));
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const languageName = user?.targetLanguage
    ? LANGUAGE_NAMES[user.targetLanguage] || "your language"
    : "your language";

  const totalLessons =
    skills?.reduce((n, s) => n + (s.lessonCount || 0), 0) ?? 0;
  const doneLessons =
    skills?.reduce((n, s) => n + (s.completedCount || 0), 0) ?? 0;
  const pct = totalLessons ? Math.round((doneLessons / totalLessons) * 100) : 0;

  // Preserve backend order while grouping into units.
  const units: { name: string; skills: Skill[] }[] = [];
  (skills || []).forEach((s) => {
    const name = s.unit || "Course";
    let u = units.find((x) => x.name === name);
    if (!u) {
      u = { name, skills: [] };
      units.push(u);
    }
    u.skills.push(s);
  });

  return (
    <div className="flex flex-col bg-cream">
      {/* Header */}
      <header className="bg-purple px-5 pb-6 pt-12 text-white lg:px-8 lg:pt-8">
        <p className="text-label-lg font-bold uppercase tracking-wide text-white/60">
          {languageName} course
        </p>
        <h1 className="mt-0.5 text-heading-xl font-extrabold">Your path</h1>

        <div className="mt-4 flex items-center gap-3">
          <div className="h-2 flex-1 overflow-hidden rounded-full bg-white/20">
            <motion.div
              className="h-full rounded-full bg-amber"
              initial={{ width: 0 }}
              animate={{ width: `${pct}%` }}
              transition={{ type: "spring", stiffness: 120, damping: 20 }}
            />
          </div>
          <span className="text-body-sm font-bold">
            {doneLessons}/{totalLessons || "—"} lessons
          </span>
        </div>
        <p className="mt-1 text-body-sm text-white/70">
          {user?.xp ?? 0} XP · {user?.levelName || "Spark"}
        </p>
      </header>

      <div className="flex-1 px-4 py-5 lg:px-8 lg:py-7">
        {error ? (
          <ErrorCard onRetry={load} />
        ) : skills === null ? (
          <LoadingState />
        ) : skills.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-8">
            {units.map((u) => (
              <section key={u.name}>
                <h2 className="mb-3 px-1 text-label-lg font-bold uppercase tracking-wide text-gray-500">
                  {u.name}
                </h2>
                <div className="grid gap-3 lg:grid-cols-2">
                  {u.skills.map((s) => (
                    <SkillCard
                      key={s.id}
                      skill={s}
                      onOpen={(id) => router.push(`/lesson/${id}`)}
                      userXp={user?.xp ?? 0}
                    />
                  ))}
                </div>

                {/* Unit listening sessions */}
                {sessions
                  .filter((ls) => ls.unit === u.name)
                  .map((ls) => (
                    <ListeningCard key={ls.id} session={ls} />
                  ))}
              </section>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function ListeningCard({ session }: { session: ListeningSession }) {
  return (
    <Link
      href={`/listening/${session.id}`}
      className="mt-3 flex items-center gap-3 rounded-2xl border border-purple/15 bg-purple-light/50 p-4 transition hover:shadow-card"
    >
      <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-purple text-white">
        <Headphones size={22} />
      </span>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <p className="truncate text-heading-sm font-extrabold text-ink">
            {session.title}
          </p>
          <span className="rounded-full bg-white px-2 py-0.5 text-label-sm font-bold text-purple">
            Listening
          </span>
        </div>
        <p className="truncate text-body-sm text-slatey">
          {session.description}
        </p>
      </div>
      <span className="shrink-0 rounded-full bg-amber/20 px-2 py-0.5 text-label-md font-bold text-amber">
        +{session.xpReward} XP
      </span>
    </Link>
  );
}

function SkillCard({
  skill,
  onOpen,
  userXp,
}: {
  skill: Skill;
  onOpen: (lessonId: number) => void;
  userXp: number;
}) {
  const lessons = skill.lessons || [];
  const pct = skill.lessonCount
    ? Math.round((skill.completedCount / skill.lessonCount) * 100)
    : 0;

  if (!skill.unlocked) {
    const remaining = Math.max(0, skill.requiredXp - userXp);
    return (
      <div className="rounded-2xl border border-gray-100 bg-white p-4 opacity-90 shadow-card">
        <div className="flex items-center gap-3">
          <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gray-100 text-gray-500">
            <Lock size={20} />
          </span>
          <div className="min-w-0 flex-1">
            <p className="truncate text-heading-sm font-extrabold text-ink">
              {skill.title}
            </p>
            <p className="truncate text-body-sm text-slatey">
              {skill.description}
            </p>
          </div>
        </div>
        <p className="mt-3 rounded-lg bg-gray-50 px-3 py-2 text-body-sm font-semibold text-slatey">
          Earn {remaining} more XP to unlock
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col rounded-2xl bg-white p-4 shadow-card">
      <div className="flex items-center gap-3">
        <span
          className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl text-white"
          style={{ backgroundColor: skill.color }}
        >
          <SkillIcon name={skill.icon} size={24} />
        </span>
        <div className="min-w-0 flex-1">
          <p className="truncate text-heading-sm font-extrabold text-ink">
            {skill.title}
          </p>
          <p className="truncate text-body-sm text-slatey">
            {skill.description}
          </p>
        </div>
        {skill.completed && (
          <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-teal text-white">
            <Check size={16} />
          </span>
        )}
      </div>

      {/* progress */}
      <div className="mt-3 flex items-center gap-2">
        <div className="h-1.5 flex-1 overflow-hidden rounded-full bg-gray-100">
          <div
            className="h-full rounded-full"
            style={{ width: `${pct}%`, backgroundColor: skill.color }}
          />
        </div>
        <span className="text-label-md font-bold text-slatey">
          {skill.completedCount}/{skill.lessonCount}
        </span>
      </div>

      {/* lessons */}
      <ul className="mt-3 divide-y divide-gray-100 border-t border-gray-100">
        {lessons.map((l, i) => {
          const done = i < skill.completedCount;
          const current = i === skill.completedCount;
          return (
            <li key={l.id}>
              <button
                onClick={() => onOpen(l.id)}
                className="flex w-full items-center gap-3 py-2.5 text-left transition hover:opacity-80"
              >
                <span
                  className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full ${
                    done
                      ? "bg-teal text-white"
                      : current
                      ? "text-white"
                      : "border-2 border-gray-200 text-gray-300"
                  }`}
                  style={current && !done ? { backgroundColor: skill.color } : undefined}
                >
                  {done ? (
                    <Check size={16} />
                  ) : current ? (
                    <Play size={14} fill="currentColor" strokeWidth={0} />
                  ) : (
                    <span className="text-label-md font-bold">{i + 1}</span>
                  )}
                </span>
                <span className="min-w-0 flex-1">
                  <span className="block truncate text-body-md font-semibold text-ink">
                    {l.title}
                  </span>
                  <span className="block text-label-md text-slatey">
                    +{l.xpReward} XP
                    {current ? " · Start now" : done ? " · Completed" : ""}
                  </span>
                </span>
                <ChevronRight size={18} className="shrink-0 text-gray-300" />
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

function LoadingState() {
  return (
    <div className="space-y-8">
      {Array.from({ length: 2 }).map((_, u) => (
        <div key={u}>
          <div className="mb-3 h-3 w-24 animate-pulse rounded-full bg-gray-200" />
          <div className="grid gap-3 lg:grid-cols-2">
            {Array.from({ length: 2 }).map((_, i) => (
              <div
                key={i}
                className="h-40 animate-pulse rounded-2xl bg-gray-100"
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="rounded-2xl bg-white p-8 text-center shadow-card">
      <p className="text-heading-sm font-extrabold text-ink">
        No lessons yet
      </p>
      <p className="mt-1 text-body-md text-slatey">
        Your course is being prepared. Please check back soon.
      </p>
    </div>
  );
}

function ErrorCard({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-2xl bg-white p-6 text-center shadow-card-lg sm:flex-row sm:text-left">
      <div className="flex-1">
        <p className="text-heading-sm font-extrabold text-ink">
          Couldn&apos;t load your course
        </p>
        <p className="mt-1 text-body-md text-slatey">
          Check your connection and try again.
        </p>
      </div>
      <Button variant="primary" className="h-12 px-6" onClick={onRetry}>
        Retry
      </Button>
    </div>
  );
}
