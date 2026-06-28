"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  Check,
  Play,
  ArrowRight,
  BookOpen,
  Zap,
  Flame,
  Map as MapIcon,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { FoxMascot } from "@/components/FoxMascot";
import { XPBar, StreakFlame, GemCounter } from "@/components/widgets";
import { NotificationBell } from "@/components/NotificationBell";
import { ChatBell } from "@/components/ChatBell";
import { SkillIcon } from "@/components/SkillIcon";
import { Avatar } from "@/components/Avatar";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import type { HomeData } from "@/lib/types";

function greeting() {
  const h = new Date().getHours();
  if (h < 12) return "Good morning";
  if (h < 18) return "Good afternoon";
  return "Late night study sesh?";
}

const STEPS = [
  { icon: BookOpen, title: "Take a lesson", desc: "Short, playful exercises" },
  { icon: Zap, title: "Earn XP", desc: "Hit your daily goal" },
  { icon: Flame, title: "Build a streak", desc: "Practise a little daily" },
];

export default function HomePage() {
  return (
    <AppShell wide>
      <HomeContent />
    </AppShell>
  );
}

function HomeContent() {
  const { user, setUser } = useAuth();
  const [data, setData] = useState<HomeData | null>(null);
  const [error, setError] = useState(false);

  const load = useCallback(() => {
    setError(false);
    api
      .home()
      .then((d) => {
        setData(d);
        setUser(d.user);
      })
      .catch(() => setError(true));
  }, [setUser]);

  useEffect(() => {
    load();
  }, [load]);

  const u = data?.user || user!;
  const quests = data?.quests || [];
  const loading = !data && !error;
  const isNewUser = (u?.xp ?? 0) === 0;

  return (
    <div className="flex flex-col">
      {/* Header */}
      <header className="bg-purple px-5 pb-5 pt-12 text-white lg:px-8 lg:pt-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="rounded-full border-2 border-white/60">
              <Avatar name={u.name} color={u.avatarColor} url={u.avatarUrl} size={40} />
            </span>
            <p className="text-heading-sm font-bold">
              {greeting()}, {u.name || "friend"}!
            </p>
          </div>
          <div className="flex items-center gap-2">
            <GemCounter gems={u.gems} />
            <StreakFlame streak={u.streak} light />
            <ChatBell />
            <NotificationBell />
          </div>
        </div>

        {u.streak > 0 && (
          <div className="mt-4 flex items-center gap-2 rounded-xl bg-white/15 px-4 py-2.5">
            <span className="text-xl">🔥</span>
            <div>
              <p className="text-heading-sm font-bold">{u.streak} day streak!</p>
              <p className="text-body-sm text-white/80">Keep it going!</p>
            </div>
          </div>
        )}
      </header>

      <div className="flex-1 px-5 py-5 lg:px-8 lg:py-7">
        <div className="grid gap-6 lg:grid-cols-5">
          {/* Primary call-to-action — the clearest "start here" on the screen */}
          <div className="lg:col-span-5">
            {error ? (
              <ErrorCard onRetry={load} />
            ) : loading ? (
              <HeroSkeleton />
            ) : (
              <StartHero data={data!} isNewUser={isNewUser} />
            )}
          </div>

          {/* First-timer guide */}
          {!loading && !error && isNewUser && (
            <section className="lg:col-span-5">
              <p className="mb-2 text-body-sm font-semibold uppercase tracking-wide text-gray-500">
                How Lumora works
              </p>
              <div className="grid gap-3 sm:grid-cols-3">
                {STEPS.map((s, i) => {
                  const Icon = s.icon;
                  return (
                    <div
                      key={s.title}
                      className="flex items-center gap-3 rounded-xl bg-white p-4 shadow-card"
                    >
                      <div className="relative flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-purple-light text-purple">
                        <Icon size={20} />
                        <span className="absolute -left-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-purple text-label-sm font-extrabold text-white">
                          {i + 1}
                        </span>
                      </div>
                      <div className="min-w-0">
                        <p className="text-body-md font-extrabold text-ink">
                          {s.title}
                        </p>
                        <p className="text-body-sm text-slatey">{s.desc}</p>
                      </div>
                    </div>
                  );
                })}
              </div>
            </section>
          )}

          {/* Left column */}
          <div className="space-y-6 lg:col-span-3">
            {/* Daily goal */}
            <section>
              <p className="mb-2 text-body-sm font-semibold uppercase tracking-wide text-gray-500">
                Your daily goal
              </p>
              <div className="rounded-xl bg-white p-5 shadow-card">
                <XPBar value={u.xpToday} max={u.dailyGoalXp} />
                <div className="mt-2 flex items-center justify-between text-body-sm">
                  <span className="font-bold text-ink">
                    {u.xpToday}/{u.dailyGoalXp} XP today
                  </span>
                  <span className="text-slatey">
                    {u.xpToday >= u.dailyGoalXp
                      ? "Goal reached! 🎉"
                      : `${Math.max(0, u.dailyGoalXp - u.xpToday)} XP to go`}
                  </span>
                </div>
              </div>
            </section>

            {/* Learning path entry point */}
            <section>
              <p className="mb-2 text-body-sm font-semibold uppercase tracking-wide text-gray-500">
                Your learning path
              </p>
              <Link
                href="/learn"
                className="flex items-center gap-3 rounded-xl bg-white p-4 shadow-card transition hover:shadow-card-lg"
              >
                <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-teal-light text-teal">
                  <MapIcon size={24} />
                </span>
                <div className="flex-1">
                  <p className="text-heading-sm font-bold text-ink">
                    Explore the galaxy map
                  </p>
                  <p className="text-body-sm text-slatey">
                    See every skill and what to learn next.
                  </p>
                </div>
                <ArrowRight size={20} className="text-gray-300" />
              </Link>
            </section>
          </div>

          {/* Right column — Daily quests */}
          <div className="space-y-6 lg:col-span-2">
            <section>
              <p className="mb-2 text-body-sm font-semibold uppercase tracking-wide text-gray-500">
                Daily quests
              </p>
              <div className="mb-3 flex items-center gap-2">
                <span className="text-2xl">🦔</span>
                <p className="text-body-sm italic text-slatey">
                  Pip: &quot;Your quests await!! Let&apos;s GO!!&quot;
                </p>
              </div>

              {loading ? (
                <div className="space-y-2">
                  {Array.from({ length: 3 }).map((_, i) => (
                    <div
                      key={i}
                      className="h-16 animate-pulse rounded-lg bg-gray-100"
                    />
                  ))}
                </div>
              ) : quests.length === 0 ? (
                <div className="rounded-lg bg-white p-4 text-body-sm text-slatey shadow-quest">
                  Complete a lesson to unlock today&apos;s quests.
                </div>
              ) : (
                <div className="space-y-2">
                  {quests.map((q, i) => (
                    <motion.div
                      key={q.id}
                      initial={{ opacity: 0, x: -8 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: i * 0.05 }}
                      className="flex items-center gap-3 rounded-lg bg-white p-3 shadow-quest"
                    >
                      <span className="text-xl">{q.quest?.icon}</span>
                      <div className="flex-1">
                        <p className="text-body-md font-semibold">
                          {q.quest?.title}
                        </p>
                        <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-100">
                          <div
                            className="h-full rounded-full bg-amber"
                            style={{
                              width: `${Math.min(
                                100,
                                (q.progress / (q.quest?.target || 1)) * 100
                              )}%`,
                            }}
                          />
                        </div>
                      </div>
                      <span className="rounded-full bg-amber/15 px-2 py-0.5 text-label-md font-bold text-amber">
                        +{q.quest?.xpReward}
                      </span>
                      <span
                        className={`flex h-6 w-6 items-center justify-center rounded-full ${
                          q.completed
                            ? "bg-teal text-white"
                            : "border-2 border-gray-100"
                        }`}
                      >
                        {q.completed && <Check size={14} />}
                      </span>
                    </motion.div>
                  ))}
                </div>
              )}
            </section>
          </div>
        </div>
      </div>
    </div>
  );
}

/** The hero card: the single, obvious place to begin learning. */
function StartHero({ data, isNewUser }: { data: HomeData; isNewUser: boolean }) {
  const lesson = data.nextLesson;
  const skill = data.nextSkill;

  // Nothing left to do right now — guide the learner to the map.
  if (!lesson) {
    return (
      <div className="flex flex-col items-center gap-3 rounded-2xl bg-white p-6 text-center shadow-card-lg sm:flex-row sm:text-left">
        <FoxMascot size={96} glow />
        <div className="flex-1">
          <p className="text-heading-lg font-extrabold text-ink">
            You&apos;re all caught up! 🎉
          </p>
          <p className="mt-1 text-body-md text-slatey">
            Great work. Unlock new skills on your learning path.
          </p>
        </div>
        <Link href="/learn">
          <Button variant="primary" className="h-12 px-6">
            Explore the map
          </Button>
        </Link>
      </div>
    );
  }

  const label = isNewUser ? "Start here" : "Continue where you left off";
  const cta = isNewUser ? "Start your first lesson" : "Continue lesson";

  return (
    <div className="rounded-2xl bg-purple p-5 text-white shadow-float lg:p-6">
      <p className="text-label-lg font-bold uppercase tracking-wide text-white/70">
        {label}
      </p>
      <div className="mt-3 flex items-center gap-4">
        <span className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-white/15">
          <SkillIcon name={skill?.icon} size={28} className="text-white" />
        </span>
        <div className="min-w-0 flex-1">
          <p className="truncate text-heading-lg font-extrabold">
            {lesson.title}
          </p>
          <p className="truncate text-body-md text-white/80">
            {skill?.title || "Your next lesson"}
          </p>
          <div className="mt-1.5 flex flex-wrap items-center gap-2 text-label-md font-bold">
            <span className="rounded-full bg-white/15 px-2.5 py-0.5">
              Lesson {lesson.orderIndex}
            </span>
            <span className="rounded-full bg-amber px-2.5 py-0.5 text-ink">
              +{lesson.xpReward} XP
            </span>
            <span className="rounded-full bg-white/15 px-2.5 py-0.5">
              ~5 min
            </span>
          </div>
        </div>
      </div>

      <Link
        href={`/lesson/${lesson.id}`}
        className="mt-5 flex h-12 w-full items-center justify-center gap-2 rounded-full bg-white text-heading-sm font-extrabold text-purple transition hover:bg-white/90"
      >
        <Play size={18} fill="#6C3FC5" strokeWidth={0} />
        {cta}
      </Link>
    </div>
  );
}

function HeroSkeleton() {
  return (
    <div className="rounded-2xl bg-white p-5 shadow-card-lg lg:p-6">
      <div className="h-3 w-32 animate-pulse rounded-full bg-gray-100" />
      <div className="mt-3 flex items-center gap-4">
        <div className="h-16 w-16 animate-pulse rounded-2xl bg-gray-100" />
        <div className="flex-1 space-y-2">
          <div className="h-4 w-2/3 animate-pulse rounded-full bg-gray-100" />
          <div className="h-3 w-1/3 animate-pulse rounded-full bg-gray-100" />
        </div>
      </div>
      <div className="mt-5 h-12 w-full animate-pulse rounded-full bg-gray-100" />
    </div>
  );
}

function ErrorCard({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-2xl bg-white p-6 text-center shadow-card-lg sm:flex-row sm:text-left">
      <div className="flex-1">
        <p className="text-heading-md font-extrabold text-ink">
          Couldn&apos;t load your dashboard
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
