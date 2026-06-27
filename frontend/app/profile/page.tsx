"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Flame, Zap, BookOpen, Globe, LogOut, ChevronRight } from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { CharacterWithFriendship } from "@/lib/types";

export default function ProfilePage() {
  const { user, logout } = useAuth();
  const [characters, setCharacters] = useState<CharacterWithFriendship[]>([]);
  const [lessonsDone, setLessonsDone] = useState(0);

  useEffect(() => {
    api
      .characters()
      .then((r) => setCharacters(r.characters))
      .catch(() => setCharacters([]));
    // Lessons completed is derived from total XP as a friendly estimate.
    api
      .home()
      .then(() => {})
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (user) setLessonsDone(Math.max(0, Math.floor(user.xp / 10)));
  }, [user]);

  if (!user) return null;

  const fluencyPct = Math.min(100, Math.round((user.fluencyScore / 1000) * 100));

  return (
    <AppShell tabs wide>
      <div className="bg-cream pb-24 lg:pb-10">
        {/* Purple header */}
        <div className="relative rounded-b-[32px] bg-gradient-to-b from-purple to-purple-dark px-6 pb-8 pt-14 text-white">
          <div className="flex items-center gap-4">
            <div
              className="flex h-20 w-20 items-center justify-center rounded-full border-4 border-white/30 text-4xl font-extrabold"
              style={{ backgroundColor: user.avatarColor || "#6C3FC5" }}
            >
              {user.name?.charAt(0).toUpperCase() || "L"}
            </div>
            <div className="min-w-0">
              <h1 className="truncate text-heading-xl font-extrabold">
                {user.name}
              </h1>
              <p className="truncate text-body-sm text-purple-light">
                {user.email}
              </p>
              <span className="mt-1 inline-block rounded-full bg-amber px-3 py-0.5 text-label-md font-extrabold text-ink">
                {user.cefrLevel} · {user.levelName}
              </span>
            </div>
          </div>
        </div>

        {/* Fluency ring */}
        <div className="-mt-6 px-6">
          <div className="flex items-center gap-5 rounded-2xl bg-white p-5 shadow-card">
            <FluencyRing pct={fluencyPct} score={user.fluencyScore} />
            <div>
              <div className="text-heading-md font-extrabold text-ink">
                Fluency Score
              </div>
              <p className="text-body-sm text-slatey">
                {user.targetLanguage
                  ? `Your ${user.targetLanguage} mastery, out of 1000.`
                  : "Your mastery, out of 1000."}
              </p>
            </div>
          </div>
        </div>

        {/* Stats row */}
        <div className="mt-5 grid grid-cols-2 gap-3 px-6 lg:grid-cols-4">
          <StatTile icon={<Flame size={20} />} label="Day Streak" value={user.streak} tint="#FF5C5C" />
          <StatTile icon={<Zap size={20} />} label="Total XP" value={user.xp} tint="#F5A623" />
          <StatTile icon={<BookOpen size={20} />} label="Lessons" value={lessonsDone} tint="#6C3FC5" />
          <StatTile icon={<Globe size={20} />} label="Languages" value={user.targetLanguage ? 1 : 0} tint="#00C2A8" />
        </div>

        {/* Characters */}
        <div className="mt-7 px-6">
          <h2 className="mb-3 text-heading-md font-extrabold text-ink">
            Your Companions
          </h2>
          <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 lg:grid-cols-6">
            {characters.map((c) => (
              <motion.div
                key={c.id}
                whileTap={{ scale: 0.95 }}
                className="flex flex-col items-center rounded-2xl bg-white p-3 shadow-card"
              >
                <div
                  className="flex h-14 w-14 items-center justify-center rounded-full text-3xl"
                  style={{ backgroundColor: (c.color || "#EDE7F6") + "33" }}
                >
                  {c.emoji}
                </div>
                <div className="mt-2 text-center text-label-md font-extrabold text-ink">
                  {c.name}
                </div>
                <div className="mt-1 flex gap-0.5">
                  {Array.from({ length: 3 }).map((_, i) => (
                    <span
                      key={i}
                      className={`h-1.5 w-1.5 rounded-full ${
                        i < c.friendshipLevel ? "bg-amber" : "bg-gray-300"
                      }`}
                    />
                  ))}
                </div>
              </motion.div>
            ))}
            {characters.length === 0 &&
              Array.from({ length: 6 }).map((_, i) => (
                <div
                  key={i}
                  className="h-28 animate-pulse rounded-2xl bg-gray-100"
                />
              ))}
          </div>
        </div>

        {/* Settings */}
        <div className="mt-7 px-6">
          <div className="overflow-hidden rounded-2xl bg-white shadow-card">
            <SettingsRow label="Account Settings" />
            <SettingsRow label="Notifications" />
            <SettingsRow label="Help & Support" />
            <button
              onClick={logout}
              className="flex w-full items-center justify-between px-5 py-4 text-coral hover:bg-coral-light"
            >
              <span className="flex items-center gap-3 font-extrabold">
                <LogOut size={18} /> Sign Out
              </span>
              <ChevronRight size={18} />
            </button>
          </div>
        </div>
      </div>
    </AppShell>
  );
}

function FluencyRing({ pct, score }: { pct: number; score: number }) {
  const r = 42;
  const c = 2 * Math.PI * r;
  const offset = c - (pct / 100) * c;
  return (
    <div className="relative h-24 w-24 shrink-0">
      <svg width={96} height={96} className="-rotate-90">
        <circle cx={48} cy={48} r={r} stroke="#EDE7F6" strokeWidth={8} fill="none" />
        <motion.circle
          cx={48}
          cy={48}
          r={r}
          stroke="#6C3FC5"
          strokeWidth={8}
          fill="none"
          strokeLinecap="round"
          strokeDasharray={c}
          initial={{ strokeDashoffset: c }}
          animate={{ strokeDashoffset: offset }}
          transition={{ duration: 1, ease: "easeOut" }}
        />
      </svg>
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <span className="text-heading-md font-extrabold text-purple">{score}</span>
        <span className="text-label-sm text-slatey">/ 1000</span>
      </div>
    </div>
  );
}

function StatTile({
  icon,
  label,
  value,
  tint,
}: {
  icon: React.ReactNode;
  label: string;
  value: number;
  tint: string;
}) {
  return (
    <div className="flex items-center gap-3 rounded-2xl bg-white p-4 shadow-card">
      <div
        className="flex h-10 w-10 items-center justify-center rounded-xl"
        style={{ backgroundColor: tint + "1A", color: tint }}
      >
        {icon}
      </div>
      <div>
        <div className="text-heading-md font-extrabold text-ink">{value}</div>
        <div className="text-label-md text-slatey">{label}</div>
      </div>
    </div>
  );
}

function SettingsRow({ label }: { label: string }) {
  return (
    <button className="flex w-full items-center justify-between border-b border-gray-100 px-5 py-4 text-ink hover:bg-gray-50">
      <span className="font-semibold">{label}</span>
      <ChevronRight size={18} className="text-gray-300" />
    </button>
  );
}
