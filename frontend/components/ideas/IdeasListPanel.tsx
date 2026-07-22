"use client";

import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import {
  AlertTriangle,
  Archive,
  CircleCheck,
  Flame,
  Hammer,
  Lightbulb,
  MessageSquare,
  Plus,
  Search,
  Star,
  TrendingUp,
  UserRound,
  Zap,
  type LucideIcon,
} from "lucide-react";
import { Avatar } from "@/components/Avatar";
import { StatusBadge, TagChip, VoteControl, timeAgo } from "./IdeaBits";
import type { Idea, IdeaBoard } from "@/lib/types";

/** The left panel: the board itself — filter, sort, and pick an idea. */

// Icons rather than emoji: emoji render differently on every platform, don't
// inherit the chip's colour, and can't be sized against the type scale.
const FILTERS: { key: string; label: string; icon?: LucideIcon }[] = [
  { key: "", label: "All" },
  { key: "starred", label: "Pinned", icon: Star },
  { key: "under_review", label: "Active", icon: Lightbulb },
  { key: "in_progress", label: "Building", icon: Hammer },
  { key: "completed", label: "Done", icon: CircleCheck },
  { key: "mine", label: "Mine", icon: UserRound },
  { key: "archived", label: "Archived", icon: Archive },
];

const SORTS: { key: string; label: string; icon: typeof Flame }[] = [
  { key: "hot", label: "Hot", icon: Flame },
  { key: "top", label: "Top", icon: TrendingUp },
  { key: "new", label: "Newest", icon: Zap },
  { key: "controversial", label: "Controversial", icon: AlertTriangle },
];

export function IdeasListPanel({
  board,
  loading,
  selectedId,
  filter,
  sort,
  tag,
  query,
  onSelect,
  onFilter,
  onSort,
  onTag,
  onQuery,
  onVote,
  onStar,
  onNew,
}: {
  board: IdeaBoard | null;
  loading: boolean;
  selectedId: number | null;
  filter: string;
  sort: string;
  tag: string;
  query: string;
  onSelect: (id: number) => void;
  onFilter: (f: string) => void;
  onSort: (s: string) => void;
  onTag: (t: string) => void;
  onQuery: (q: string) => void;
  onVote: (id: number, value: number) => void;
  onStar: (id: number) => void;
  onNew: () => void;
}) {
  const [sortOpen, setSortOpen] = useState(false);
  const activeSort = SORTS.find((s) => s.key === sort) ?? SORTS[0];

  return (
    <div className="flex h-full flex-col border-r border-gray-100 bg-white">
      {/* Header — search and New share one row. The panel had its own "Ideas"
          title, which on a phone sat directly under the toolbar's, costing a
          whole row of vertical space to say the same word twice. */}
      <div className="flex shrink-0 items-center gap-2 border-b border-gray-100 px-4 py-3">
        <div className="relative min-w-0 flex-1">
          <Search
            size={15}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-300"
          />
          <input
            value={query}
            onChange={(e) => onQuery(e.target.value)}
            placeholder="Search ideas…"
            className="w-full rounded-full bg-gray-50 py-2 pl-9 pr-3 text-body-md text-ink outline-none ring-purple/30 placeholder:text-gray-300 focus:ring-2"
          />
        </div>
        <button
          onClick={onNew}
          className="flex shrink-0 items-center gap-1 rounded-full bg-purple px-3 py-2 text-label-lg font-extrabold text-white shadow-float transition hover:bg-purple-dark"
        >
          <Plus size={15} /> New
        </button>
      </div>

      {/* Filters */}
      <div className="shrink-0 space-y-2 border-b border-gray-100 px-4 py-3">
        <div className="flex flex-wrap gap-1.5">
          {FILTERS.map((f) => (
            <button
              key={f.key}
              onClick={() => onFilter(f.key)}
              className={`flex items-center gap-1 rounded-full px-2.5 py-1 text-label-md font-bold transition ${
                filter === f.key
                  ? "bg-purple text-white"
                  : "bg-gray-50 text-slatey hover:bg-gray-100"
              }`}
            >
              {f.icon && <f.icon size={12} />}
              {f.label}
              {board?.counts?.[countKeyFor(f.key)] !== undefined && (
                <span className="opacity-60">
                  {board.counts[countKeyFor(f.key)]}
                </span>
              )}
            </button>
          ))}
        </div>

        <div className="flex items-center justify-between">
          <div className="relative">
            <button
              onClick={() => setSortOpen((o) => !o)}
              className="flex items-center gap-1.5 text-label-lg font-bold text-slatey hover:text-ink"
            >
              <activeSort.icon size={14} /> Sort: {activeSort.label}
            </button>
            <AnimatePresence>
              {sortOpen && (
                <motion.div
                  initial={{ opacity: 0, y: -4 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0 }}
                  className="absolute left-0 top-full z-20 mt-1 w-52 rounded-2xl bg-white p-1.5 shadow-card-lg"
                >
                  {SORTS.map((s) => (
                    <button
                      key={s.key}
                      onClick={() => {
                        onSort(s.key);
                        setSortOpen(false);
                      }}
                      className={`flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-body-sm font-bold transition ${
                        s.key === sort
                          ? "bg-purple-light text-purple"
                          : "text-slatey hover:bg-gray-50"
                      }`}
                    >
                      <s.icon size={14} /> {s.label}
                      {s.key === "controversial" && (
                        <span className="ml-auto text-label-sm font-normal text-gray-500">
                          evenly split
                        </span>
                      )}
                    </button>
                  ))}
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {tag && (
            <button
              onClick={() => onTag("")}
              className="text-label-md font-bold text-purple hover:underline"
            >
              Clear #{tag}
            </button>
          )}
        </div>

        {board && board.tags.length > 0 && (
          <div className="flex flex-wrap gap-1.5 pt-1">
            {board.tags.slice(0, 10).map((t) => (
              <TagChip
                key={t.tag}
                tag={t.tag}
                active={t.tag === tag}
                onClick={() => onTag(t.tag === tag ? "" : t.tag)}
              />
            ))}
          </div>
        )}
      </div>

      {/* A board that only grows is a board nobody reads. */}
      {board?.crowded && (
        <div className="mx-4 mt-3 flex items-start gap-2 rounded-lg bg-amber-light px-3 py-2 text-label-md text-ink">
          <AlertTriangle size={14} className="mt-0.5 shrink-0 text-amber" />
          <span>
            {board.openIdeas} ideas are open. Archiving the ones that have run
            their course (with a reason) keeps this list readable.
          </span>
        </div>
      )}

      {/* The list */}
      <div className="min-h-0 flex-1 overflow-y-auto p-3">
        {loading && (
          <div className="space-y-2">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="h-24 animate-pulse rounded-2xl bg-gray-50" />
            ))}
          </div>
        )}

        {!loading && board?.ideas.length === 0 && (
          <div className="px-4 py-12 text-center">
            <p className="text-body-md font-bold text-ink">Nothing here yet</p>
            <p className="mt-1 text-body-sm text-slatey">
              {query || tag
                ? "No ideas match that filter."
                : "Post the first one — a title is all it takes."}
            </p>
          </div>
        )}

        <div className="space-y-2">
          {board?.ideas.map((idea, i) => (
            <IdeaCard
              key={idea.id}
              idea={idea}
              index={i}
              selected={idea.id === selectedId}
              onSelect={() => onSelect(idea.id)}
              onVote={(v) => onVote(idea.id, v)}
              onStar={() => onStar(idea.id)}
            />
          ))}
        </div>
      </div>
    </div>
  );
}

function countKeyFor(filter: string): string {
  switch (filter) {
    case "":
      return "all";
    case "under_review":
      return "underReview";
    case "in_progress":
      return "inProgress";
    default:
      return filter;
  }
}

function IdeaCard({
  idea,
  index,
  selected,
  onSelect,
  onVote,
  onStar,
}: {
  idea: Idea;
  index: number;
  selected: boolean;
  onSelect: () => void;
  onVote: (v: number) => void;
  onStar: () => void;
}) {
  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: Math.min(index * 0.02, 0.3) }}
      onClick={onSelect}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === "Enter" && onSelect()}
      className={`cursor-pointer rounded-2xl border-2 p-3 transition ${
        selected
          ? "border-purple bg-purple-light"
          : "border-transparent bg-white hover:bg-gray-50"
      } ${idea.archived ? "opacity-60" : ""}`}
    >
      <div className="flex gap-3">
        <VoteControl idea={idea} onVote={onVote} disabled={idea.archived} />

        <div className="min-w-0 flex-1">
          <div className="flex items-start gap-2">
            <p className="min-w-0 flex-1 font-extrabold leading-snug text-ink">
              <span className="text-slatey">#{idea.id}</span> {idea.title}
            </p>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onStar();
              }}
              aria-label={idea.starred ? "Unpin" : "Pin for later"}
              className={`shrink-0 rounded-md p-1 transition ${
                idea.starred
                  ? "text-amber"
                  : "text-gray-300 hover:bg-gray-100 hover:text-slatey"
              }`}
            >
              <Star size={14} fill={idea.starred ? "currentColor" : "none"} />
            </button>
          </div>

          {idea.description && (
            <p className="mt-1 line-clamp-2 text-body-sm text-slatey">
              {idea.description}
            </p>
          )}

          <div className="mt-2 flex flex-wrap items-center gap-1.5">
            <StatusBadge status={idea.status} />
            {idea.tags.slice(0, 3).map((t) => (
              <TagChip key={t} tag={t} />
            ))}
          </div>

          <div className="mt-2 flex items-center gap-2 text-label-md text-slatey">
            <Avatar
              name={idea.owner.name}
              color={idea.owner.avatarColor}
              url={idea.owner.avatarUrl}
              size={18}
            />
            <span className="truncate">{idea.owner.name}</span>
            {idea.messageCount > 0 && (
              <span className="flex shrink-0 items-center gap-1">
                <MessageSquare size={11} /> {idea.messageCount}
              </span>
            )}
            <span className="ml-auto shrink-0">{timeAgo(idea.lastActivity)}</span>
          </div>
        </div>
      </div>
    </motion.div>
  );
}
