"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import {
  ChevronDown,
  Crown,
  Flag,
  Flame,
  Gem,
  Info,
  Lock,
  ShieldCheck,
  Sparkles,
  Swords,
  Target,
  Timer,
  TrendingDown,
  TrendingUp,
  Trophy,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { Avatar } from "@/components/Avatar";
import { Button } from "@/components/Button";
import { LeagueCeremony } from "@/components/LeagueCeremony";
import { api } from "@/lib/api";
import { languageMeta } from "@/lib/languages";
import type { LeagueResult, LeagueRow, LeagueStandings } from "@/lib/types";

/**
 * The weekly league.
 *
 * Three things distinguish this from a plain leaderboard, and the UI has to say
 * all three out loud or the scoring feels arbitrary:
 *
 *   1. The number racing isn't XP, it's weighted league points — so the row
 *      shows both, and "How scoring works" explains the weighting.
 *   2. Promotion and demotion are zones, not a vague notion of "doing well" —
 *      so the list is physically divided by labelled dividers.
 *   3. Nobody is forced to compete — casual mode is one tap away, right here,
 *      not buried in settings.
 */
export default function LeaguePage() {
  const [data, setData] = useState<LeagueStandings | null>(null);
  const [result, setResult] = useState<LeagueResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [explainOpen, setExplainOpen] = useState(false);
  const [busy, setBusy] = useState(false);

  const load = useCallback(() => {
    return api
      .league()
      .then(setData)
      .catch(() => setData(null));
  }, []);

  useEffect(() => {
    Promise.all([
      load(),
      // A settled season is handed over exactly once; the ceremony plays on top
      // of the new week's standings.
      api
        .leagueResult()
        .then((r) => setResult(r.result))
        .catch(() => setResult(null)),
    ]).finally(() => setLoading(false));
  }, [load]);

  const closeCeremony = async () => {
    if (result) {
      await api.markLeagueResultSeen(result.seasonId).catch(() => {});
    }
    setResult(null);
    load();
  };

  const toggleCasual = async (enabled: boolean) => {
    setBusy(true);
    try {
      await api.setLeagueCasual(enabled);
      await load();
    } finally {
      setBusy(false);
    }
  };

  const report = async (row: LeagueRow) => {
    if (row.reported || row.isUser) return;
    setData((d) =>
      d
        ? { ...d, rows: d.rows.map((r) => (r.id === row.id ? { ...r, reported: true } : r)) }
        : d
    );
    await api.reportLeagueMember(row.id, "suspicious activity").catch(() => {});
  };

  const tint = data?.tier.tint ?? "#6C3FC5";

  return (
    <AppShell tabs>
      <AnimatePresence>
        {result && <LeagueCeremony result={result} onClose={closeCeremony} />}
      </AnimatePresence>

      <div className="bg-cream pb-24 lg:pb-10">
        <Header data={data} loading={loading} tint={tint} />

        <div className="mx-auto max-w-2xl px-4">
          {data && !data.casual && data.stage && <TournamentBanner stage={data.stage} />}

          {data && !loading && <TierRail data={data} />}

          {data?.groupGoal && data.joined && (
            <GroupGoal goal={data.groupGoal} tint={tint} />
          )}

          {data?.me?.flagged && (
            <div className="mt-4 rounded-2xl border-2 border-coral/40 bg-coral-light px-4 py-3 text-body-sm text-ink">
              <strong>Your week is under review.</strong> {data.me.flagReason}. You
              keep every point of XP, but promotion and chest rewards are paused
              for this season.
            </div>
          )}

          <ScoringExplainer open={explainOpen} onToggle={() => setExplainOpen((o) => !o)} />

          {loading && <RowsSkeleton />}

          {!loading && data?.casual && (
            <CasualCard busy={busy} onLeave={() => toggleCasual(false)} />
          )}

          {!loading && data && !data.casual && !data.joined && <NotJoinedCard />}

          {!loading && data && !data.casual && data.joined && (
            <Standings data={data} onReport={report} />
          )}

          {!loading && data && !data.casual && (
            <button
              onClick={() => toggleCasual(true)}
              disabled={busy}
              className="mx-auto mt-8 block text-label-lg font-bold text-slatey underline decoration-dotted underline-offset-4 hover:text-ink disabled:opacity-40"
            >
              Switch to casual mode — learn without the leaderboard
            </button>
          )}
        </div>
      </div>
    </AppShell>
  );
}

// --- header ------------------------------------------------------------------

function Header({
  data,
  loading,
  tint,
}: {
  data: LeagueStandings | null;
  loading: boolean;
  tint: string;
}) {
  const remaining = useCountdown(data?.secondsRemaining ?? 0);

  return (
    <div
      className="relative overflow-hidden rounded-b-[32px] px-6 pb-8 pt-14 text-center text-white"
      style={{ background: `linear-gradient(160deg, ${tint}, ${tint}bb)` }}
    >
      <motion.div
        className="pointer-events-none absolute -left-16 -top-20 h-64 w-64 rounded-full bg-white/10"
        animate={{ scale: [1, 1.15, 1] }}
        transition={{ duration: 8, repeat: Infinity, ease: "easeInOut" }}
      />

      <motion.div
        initial={{ scale: 0.7, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ type: "spring", stiffness: 220, damping: 16 }}
        className="relative mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur"
      >
        <Crown size={32} />
      </motion.div>

      <h1 className="relative mt-3 text-heading-xl font-extrabold">
        {data ? `${data.tier.name} League` : "League"}
      </h1>

      {!loading && data && (
        <p className="relative mt-1 text-body-sm text-white/80">
          {data.casual
            ? "Casual mode — you're learning without a leaderboard"
            : data.joined
              ? `You're #${data.userRank} of ${data.rows.length} · top ${data.promoteTop} promote`
              : "Complete one lesson to enter this week's race"}
        </p>
      )}

      {!loading && data && !data.casual && (
        <div className="relative mt-4 inline-flex items-center gap-2 rounded-full bg-black/20 px-4 py-1.5 text-label-lg font-bold backdrop-blur">
          <Timer size={14} />
          {remaining} left
        </div>
      )}

      {!loading && data && (
        <div className="relative mt-4 flex items-center justify-center gap-2 text-label-md">
          {data.you.fairPlay && (
            <span className="flex items-center gap-1 rounded-full bg-white/15 px-3 py-1 font-bold">
              <ShieldCheck size={12} /> Fair play
            </span>
          )}
          {data.you.trophies > 0 && (
            <span className="flex items-center gap-1 rounded-full bg-white/15 px-3 py-1 font-bold">
              <Trophy size={12} /> {data.you.trophies}
            </span>
          )}
          <span className="flex items-center gap-1 rounded-full bg-white/15 px-3 py-1 font-bold">
            Best: {data.you.best}
          </span>
        </div>
      )}
    </div>
  );
}

/** The ten-tier ladder, so the climb has somewhere visible to go. */
function TierRail({ data }: { data: LeagueStandings }) {
  return (
    <div className="-mx-4 mt-5 overflow-x-auto px-4 pb-1">
      <div className="flex min-w-max items-end gap-2">
        {data.tiers.map((t) => {
          const current = t.index === data.tier.index;
          const reached = t.index <= data.you.bestTier;
          return (
            <motion.div
              key={t.name}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: Math.min(t.index * 0.03, 0.3) }}
              className={`flex w-[68px] flex-col items-center rounded-2xl px-2 py-2 ${
                current ? "bg-white shadow-card" : ""
              }`}
            >
              <div
                className="flex h-9 w-9 items-center justify-center rounded-xl"
                style={{
                  background: reached || current ? t.tint : "#EBEBEB",
                  opacity: current ? 1 : reached ? 0.75 : 1,
                  boxShadow: current ? `0 4px 14px ${t.tint}66` : undefined,
                }}
              >
                {reached || current ? (
                  <Crown size={16} className="text-white" />
                ) : (
                  <Lock size={14} className="text-gray-500" />
                )}
              </div>
              <span
                className={`mt-1 text-label-sm font-extrabold ${
                  current ? "text-ink" : "text-slatey"
                }`}
              >
                {t.name}
              </span>
              <span className="text-label-sm text-gray-500">
                {t.index === 9 ? "cup" : `top ${t.promoteTop}`}
              </span>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}

function TournamentBanner({ stage }: { stage: string }) {
  const label =
    stage === "quarterfinal"
      ? "Quarterfinal"
      : stage === "semifinal"
        ? "Semifinal"
        : "Final";
  return (
    <motion.div
      initial={{ opacity: 0, y: -8 }}
      animate={{ opacity: 1, y: 0 }}
      className="mt-4 flex items-center gap-3 rounded-2xl bg-gradient-to-r from-[#17A3DD] to-[#6C3FC5] px-4 py-3 text-white shadow-card"
    >
      <Swords size={20} />
      <div>
        <div className="text-heading-sm font-extrabold">
          Diamond Tournament · {label}
        </div>
        <div className="text-label-md text-white/80">
          {stage === "final"
            ? "Top 3 take the trophy. Everyone returns to Diamond next week."
            : "Top 10 of 30 advance to the next round."}
        </div>
      </div>
    </motion.div>
  );
}

/** The collaborative target — the part of the week that isn't zero-sum. */
function GroupGoal({
  goal,
  tint,
}: {
  goal: NonNullable<LeagueStandings["groupGoal"]>;
  tint: string;
}) {
  const pct = Math.min(100, Math.round((goal.current / Math.max(goal.target, 1)) * 100));
  return (
    <div className="mt-4 rounded-2xl bg-white p-4 shadow-card">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-heading-sm font-extrabold text-ink">
          <Target size={16} style={{ color: tint }} /> Pod goal
        </div>
        <div className="flex items-center gap-1 text-label-lg font-bold text-teal">
          <Gem size={13} /> +{goal.bonus} each
        </div>
      </div>
      <p className="mt-1 text-body-sm text-slatey">
        {goal.hit
          ? "Smashed it — every learner who scored gets bonus gems at reset."
          : "Your whole pod contributes. Hit the target and everyone who scored earns bonus gems, win or lose."}
      </p>
      <div className="mt-3 h-3 overflow-hidden rounded-full bg-gray-100">
        <motion.div
          className="h-full rounded-full"
          style={{ background: goal.hit ? "#00C2A8" : tint }}
          initial={{ width: 0 }}
          animate={{ width: `${pct}%` }}
          transition={{ duration: 0.9, ease: "easeOut" }}
        />
      </div>
      <div className="mt-1 flex justify-between text-label-md text-slatey">
        <span>
          {goal.current.toLocaleString()} / {goal.target.toLocaleString()} points
        </span>
        <span className="font-bold">{pct}%</span>
      </div>
    </div>
  );
}

function ScoringExplainer({
  open,
  onToggle,
}: {
  open: boolean;
  onToggle: () => void;
}) {
  return (
    <div className="mt-4 overflow-hidden rounded-2xl bg-purple-light">
      <button
        onClick={onToggle}
        className="flex w-full items-center justify-between px-4 py-3 text-left"
      >
        <span className="flex items-center gap-2 text-heading-sm font-extrabold text-purple">
          <Info size={16} /> How scoring works
        </span>
        <motion.span animate={{ rotate: open ? 180 : 0 }}>
          <ChevronDown size={18} className="text-purple" />
        </motion.span>
      </button>
      <AnimatePresence initial={false}>
        {open && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.25 }}
          >
            <ul className="space-y-2 px-4 pb-4 text-body-sm text-slatey">
              <li>
                <strong className="text-ink">Points aren&apos;t XP.</strong> Every
                activity is weighted by how hard the content was and how cleanly you
                cleared it — a flawless advanced lesson outscores three easy ones.
              </li>
              <li>
                <strong className="text-ink">Earlier is worth more.</strong> Monday
                points count 1.2x and Sunday 0.85x, so steady beats cramming.
              </li>
              <li>
                <strong className="text-ink">Grinding tapers off.</strong> After 500
                XP in a day, points halve; past 1,000 they quarter. Practice drills
                count for less and stop counting after five a day.
              </li>
              <li>
                <strong className="text-ink">Your pod is matched to you.</strong>{" "}
                Placement uses a hidden rating that persists across leagues, so
                joining late or tanking a week doesn&apos;t buy an easier group.
              </li>
              <li>
                <strong className="text-ink">Ties go to whoever got there first.</strong>
              </li>
            </ul>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

// --- standings ---------------------------------------------------------------

function Standings({
  data,
  onReport,
}: {
  data: LeagueStandings;
  onReport: (row: LeagueRow) => void;
}) {
  const { rows } = data;
  const promote = data.promoteTop ?? 0;
  const demote = data.demoteBottom ?? 0;
  const demoteFrom = demote > 0 ? rows.length - demote : -1;

  return (
    <div className="mt-5 space-y-2">
      {rows.map((row, i) => {
        const rank = i + 1;
        return (
          <div key={row.id}>
            {rank === promote + 1 && promote > 0 && (
              <ZoneDivider
                tone="promote"
                label={`Promotion zone — top ${promote} move up to ${
                  data.tiers[Math.min(data.tier.index + 1, data.tiers.length - 1)].name
                }`}
              />
            )}
            {demoteFrom > 0 && rank === demoteFrom + 1 && (
              <ZoneDivider
                tone="demote"
                label={`Demotion zone — bottom ${demote} move down`}
              />
            )}
            <Row row={row} index={i} tint={data.tier.tint} onReport={onReport} />
          </div>
        );
      })}
    </div>
  );
}

function ZoneDivider({
  tone,
  label,
}: {
  tone: "promote" | "demote";
  label: string;
}) {
  const promote = tone === "promote";
  const color = promote ? "#00C2A8" : "#FF5C5C";
  return (
    <motion.div
      initial={{ opacity: 0, scaleX: 0.85 }}
      animate={{ opacity: 1, scaleX: 1 }}
      className="my-3 flex items-center gap-2"
    >
      <span
        className="flex items-center gap-1 rounded-full px-3 py-1 text-label-md font-extrabold text-white"
        style={{ background: color }}
      >
        {promote ? <TrendingUp size={12} /> : <TrendingDown size={12} />}
        {promote ? "Promotion" : "Demotion"}
      </span>
      <span className="flex-1 text-label-md font-bold text-slatey">{label}</span>
      <span className="h-[2px] w-6 rounded-full" style={{ background: color }} />
    </motion.div>
  );
}

function Row({
  row,
  index,
  tint,
  onReport,
}: {
  row: LeagueRow;
  index: number;
  tint: string;
  onReport: (row: LeagueRow) => void;
}) {
  const [menu, setMenu] = useState(false);
  const zoneRing =
    row.zone === "promote"
      ? "ring-teal/40"
      : row.zone === "demote"
        ? "ring-coral/40"
        : "ring-transparent";

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: Math.min(index * 0.03, 0.4) }}
      className={`relative flex items-center gap-3 rounded-2xl px-4 py-3 shadow-card ring-2 ${
        row.isUser ? "bg-purple-light ring-purple" : `bg-white ${zoneRing}`
      }`}
    >
      <div className="w-7 text-center">
        <RankBadge rank={row.rank} />
      </div>

      <Avatar name={row.name} color={row.avatarColor} url={row.avatarUrl} size={44} />

      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-1.5 truncate font-extrabold text-ink">
          <span className="truncate">{row.name}</span>
          {row.isUser && <span className="text-label-sm text-purple">(you)</span>}
          {row.fairPlay && !row.isUser && (
            <ShieldCheck size={13} className="shrink-0 text-teal" aria-label="Fair play" />
          )}
          {row.flagged && (
            <Flag size={13} className="shrink-0 text-coral" aria-label="Under review" />
          )}
        </div>
        <div className="flex items-center gap-2 text-label-md text-slatey">
          <span className="flex items-center gap-1">
            <Flame size={12} className="text-coral" /> {row.streak}
          </span>
          {row.accuracy > 0 && (
            <span className="flex items-center gap-1">
              <Target size={12} className="text-teal" /> {row.accuracy}%
            </span>
          )}
          {row.perfectRuns > 0 && (
            <span className="flex items-center gap-1">
              <Sparkles size={12} className="text-amber" /> {row.perfectRuns}
            </span>
          )}
          <LangBadge code={row.language} />
        </div>
      </div>

      <div className="text-right">
        <div className="text-heading-sm font-extrabold" style={{ color: tint }}>
          {row.points.toLocaleString()}
        </div>
        <div className="text-label-sm text-slatey">{row.rawXp} XP</div>
      </div>

      {!row.isUser && (
        <button
          onClick={() => setMenu((m) => !m)}
          aria-label={`Report ${row.name}`}
          className="ml-1 rounded-full p-1 text-gray-300 transition hover:bg-gray-50 hover:text-slatey"
        >
          <Flag size={14} />
        </button>
      )}

      <AnimatePresence>
        {menu && !row.isUser && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: -4 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="absolute right-3 top-full z-20 mt-1 w-56 rounded-2xl bg-white p-3 text-left shadow-card-lg"
          >
            {row.reported ? (
              <p className="text-body-sm text-slatey">
                Reported. Thanks — reports are reviewed alongside automatic
                behaviour checks.
              </p>
            ) : (
              <>
                <p className="text-body-sm text-slatey">
                  Does this look automated? Reports never punish on their own —
                  they add weight to the anti-cheat checks.
                </p>
                <Button
                  variant="outline"
                  className="mt-2 h-9 w-full text-label-lg"
                  onClick={() => {
                    onReport(row);
                    setMenu(false);
                  }}
                >
                  Report activity
                </Button>
              </>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  );
}

// --- empty / opt-out states --------------------------------------------------

function NotJoinedCard() {
  return (
    <div className="mt-6 rounded-2xl bg-white p-6 text-center shadow-card">
      <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-purple-light">
        <Sparkles size={24} className="text-purple" />
      </div>
      <h2 className="mt-3 text-heading-lg font-extrabold text-ink">
        You&apos;re not racing yet
      </h2>
      <p className="mx-auto mt-2 max-w-sm text-body-md text-slatey">
        Score a single point this week and you&apos;ll be placed in a pod of 30
        learners matched to your level. Sit the week out and nothing happens —
        you simply pause, you don&apos;t drop.
      </p>
    </div>
  );
}

function CasualCard({ busy, onLeave }: { busy: boolean; onLeave: () => void }) {
  return (
    <div className="mt-6 rounded-2xl bg-white p-6 text-center shadow-card">
      <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-teal-light">
        <ShieldCheck size={24} className="text-teal" />
      </div>
      <h2 className="mt-3 text-heading-lg font-extrabold text-ink">Casual mode</h2>
      <p className="mx-auto mt-2 max-w-sm text-body-md text-slatey">
        You&apos;re learning without a leaderboard. Every lesson, every point of
        XP and your streak all still count — you&apos;re just not ranked against
        anyone, and there&apos;s nothing to lose on Sunday night.
      </p>
      <Button className="mt-5" onClick={onLeave} loading={busy}>
        Join the weekly race
      </Button>
    </div>
  );
}

function RowsSkeleton() {
  return (
    <div className="mt-5 space-y-2">
      {Array.from({ length: 8 }).map((_, i) => (
        <div key={i} className="h-16 animate-pulse rounded-2xl bg-gray-100" />
      ))}
    </div>
  );
}

// --- bits --------------------------------------------------------------------

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

/** Live "3d 04h 12m" countdown to the weekly reset. */
function useCountdown(initialSeconds: number) {
  const [seconds, setSeconds] = useState(initialSeconds);

  useEffect(() => setSeconds(initialSeconds), [initialSeconds]);

  useEffect(() => {
    const t = setInterval(() => setSeconds((s) => Math.max(s - 1, 0)), 1000);
    return () => clearInterval(t);
  }, []);

  return useMemo(() => {
    const d = Math.floor(seconds / 86400);
    const h = Math.floor((seconds % 86400) / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    if (d > 0) return `${d}d ${h}h`;
    if (h > 0) return `${h}h ${m}m`;
    return `${m}m ${s}s`;
  }, [seconds]);
}
