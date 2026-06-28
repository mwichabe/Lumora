"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Crown, Flame } from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { Avatar } from "@/components/Avatar";
import { api } from "@/lib/api";
import { LeaderRow } from "@/lib/types";
import { languageMeta } from "@/lib/languages";

const LEAGUE_TINTS: Record<string, string> = {
  Bronze: "#CD7F32",
  Silver: "#9CA3AF",
  Gold: "#F5A623",
  Sapphire: "#00C2A8",
  Ruby: "#FF5C5C",
  Emerald: "#10B981",
  Amethyst: "#6C3FC5",
  Obsidian: "#1A1A2E",
};

export default function LeaderboardPage() {
  const [rows, setRows] = useState<LeaderRow[]>([]);
  const [league, setLeague] = useState("Bronze");
  const [userRank, setUserRank] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .leaderboard()
      .then((r) => {
        setRows(r.rows);
        setLeague(r.league);
        setUserRank(r.userRank);
      })
      .catch(() => setRows([]))
      .finally(() => setLoading(false));
  }, []);

  const tint = LEAGUE_TINTS[league] || "#6C3FC5";

  return (
    <AppShell tabs>
      <div className="bg-cream pb-24 lg:pb-10">
        {/* League header */}
        <div
          className="rounded-b-[32px] px-6 pb-8 pt-14 text-center text-white"
          style={{ background: `linear-gradient(160deg, ${tint}, ${tint}cc)` }}
        >
          <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur">
            <Crown size={32} />
          </div>
          <h1 className="mt-3 text-heading-xl font-extrabold">{league} League</h1>
          <p className="mt-1 text-body-sm text-white/80">
            {userRank > 0
              ? `You're ranked #${userRank} — earn XP to climb!`
              : "Earn XP to climb the league."}
          </p>
        </div>

        {/* Rows */}
        <div className="mt-5 space-y-2 px-4">
          {loading &&
            Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="h-16 animate-pulse rounded-2xl bg-gray-100" />
            ))}

          {!loading &&
            rows.map((row, i) => (
              <motion.div
                key={row.id}
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: Math.min(i * 0.04, 0.5) }}
                className={`flex items-center gap-3 rounded-2xl px-4 py-3 shadow-card ${
                  row.isUser ? "bg-purple-light ring-2 ring-purple" : "bg-white"
                }`}
              >
                <div className="w-7 text-center">
                  <RankBadge rank={row.rank} />
                </div>
                <Avatar
                  name={row.name}
                  color={row.avatarColor}
                  url={row.avatarUrl}
                  size={44}
                />
                <div className="min-w-0 flex-1">
                  <div className="truncate font-extrabold text-ink">
                    {row.name}
                    {row.isUser && (
                      <span className="ml-2 text-label-sm text-purple">(you)</span>
                    )}
                  </div>
                  <div className="flex items-center gap-2 text-label-md text-slatey">
                    <span className="flex items-center gap-1">
                      <Flame size={12} className="text-coral" /> {row.streak}
                    </span>
                    <LangBadge code={row.language} />
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-heading-sm font-extrabold text-purple">
                    {row.xp}
                  </div>
                  <div className="text-label-sm text-slatey">XP</div>
                </div>
              </motion.div>
            ))}

          {!loading && rows.length === 0 && (
            <p className="py-10 text-center text-body-md text-slatey">
              No league data yet. Complete a lesson to join the race!
            </p>
          )}
        </div>
      </div>
    </AppShell>
  );
}

function LangBadge({ code }: { code: string }) {
  const m = languageMeta(code);
  if (!m) {
    return (
      <span className="rounded-full bg-gray-100 px-2 py-0.5 text-label-sm font-bold text-gray-500">
        New learner
      </span>
    );
  }
  return (
    <span className="flex items-center gap-1 rounded-full bg-gray-50 px-2 py-0.5 text-label-sm font-bold text-slatey">
      <span>{m.flag}</span>
      {m.name}
    </span>
  );
}

function RankBadge({ rank }: { rank: number }) {
  if (rank === 1) return <span className="text-xl">🥇</span>;
  if (rank === 2) return <span className="text-xl">🥈</span>;
  if (rank === 3) return <span className="text-xl">🥉</span>;
  return <span className="text-body-md font-extrabold text-slatey">{rank}</span>;
}
