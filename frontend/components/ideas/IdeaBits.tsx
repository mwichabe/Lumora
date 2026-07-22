"use client";

import { motion } from "framer-motion";
import { ChevronDown, ChevronUp, ThumbsDown, ThumbsUp } from "lucide-react";
import type { Idea, IdeaStatus } from "@/lib/types";

/** Small pieces shared across the three panels of the ideas workspace. */

export const STATUS_META: Record<IdeaStatus, { label: string; tint: string; bg: string }> = {
  draft: { label: "Draft", tint: "#9090A0", bg: "#F5F5F5" },
  under_review: { label: "Under review", tint: "#F5A623", bg: "#FFF8E7" },
  approved: { label: "Approved", tint: "#6C3FC5", bg: "#EDE7F6" },
  in_progress: { label: "In progress", tint: "#17A3DD", bg: "#E5F4FB" },
  completed: { label: "Completed", tint: "#00C2A8", bg: "#E0FAF7" },
  archived: { label: "Archived", tint: "#4A4A6A", bg: "#EBEBEB" },
};

export function StatusBadge({
  status,
  size = "sm",
}: {
  status: IdeaStatus;
  size?: "sm" | "md";
}) {
  const meta = STATUS_META[status] ?? STATUS_META.draft;
  return (
    <span
      className={`inline-flex shrink-0 items-center rounded-full font-extrabold ${
        size === "md" ? "px-3 py-1 text-label-lg" : "px-2 py-0.5 text-label-sm"
      }`}
      style={{ background: meta.bg, color: meta.tint }}
    >
      {meta.label}
    </span>
  );
}

/**
 * One-tap voting. Deliberately no confirmation and no comment box — asking
 * people to justify support is how you end up with a board nobody votes on.
 */
export function VoteControl({
  idea,
  onVote,
  vertical = true,
  disabled,
}: {
  idea: Idea;
  onVote: (value: number) => void;
  vertical?: boolean;
  disabled?: boolean;
}) {
  const cast = (v: number) => {
    if (disabled) return;
    // Tapping your existing vote clears it.
    onVote(idea.myVote === v ? 0 : v);
  };

  return (
    <div
      className={`flex items-center ${vertical ? "flex-col" : "gap-1"} shrink-0`}
      onClick={(e) => e.stopPropagation()}
    >
      <button
        onClick={() => cast(1)}
        disabled={disabled}
        aria-label="Upvote"
        aria-pressed={idea.myVote === 1}
        className={`rounded-md p-1 transition disabled:opacity-40 ${
          idea.myVote === 1
            ? "text-teal"
            : "text-gray-300 hover:bg-gray-50 hover:text-slatey"
        }`}
      >
        <ChevronUp size={vertical ? 20 : 16} />
      </button>

      <motion.span
        key={idea.score}
        initial={{ scale: 1.3 }}
        animate={{ scale: 1 }}
        className={`min-w-[24px] text-center font-extrabold ${
          vertical ? "text-body-md" : "text-label-lg"
        } ${idea.score > 0 ? "text-ink" : idea.score < 0 ? "text-coral" : "text-slatey"}`}
      >
        {idea.score}
      </motion.span>

      <button
        onClick={() => cast(-1)}
        disabled={disabled}
        aria-label="Downvote"
        aria-pressed={idea.myVote === -1}
        className={`rounded-md p-1 transition disabled:opacity-40 ${
          idea.myVote === -1
            ? "text-coral"
            : "text-gray-300 hover:bg-gray-50 hover:text-slatey"
        }`}
      >
        <ChevronDown size={vertical ? 20 : 16} />
      </button>
    </div>
  );
}

/**
 * The support bar in the details panel: support against opposition, at a
 * glance. Votes alone are a poor signal — 30 up with 28 down is a very
 * different thing from 30 up with none, and only this shows the difference.
 */
export function VoteBar({ idea }: { idea: Idea }) {
  const total = idea.upvotes + idea.downvotes;
  const upPct = total === 0 ? 0 : Math.round((idea.upvotes / total) * 100);

  return (
    <div>
      <div className="flex items-center justify-between text-label-md font-bold">
        <span className="flex items-center gap-1 text-teal">
          <ThumbsUp size={12} /> {idea.upvotes}
        </span>
        <span className="text-slatey">
          {total === 0 ? "No votes yet" : `${upPct}% support`}
        </span>
        <span className="flex items-center gap-1 text-coral">
          {idea.downvotes} <ThumbsDown size={12} />
        </span>
      </div>
      <div className="mt-1.5 flex h-2.5 overflow-hidden rounded-full bg-gray-100">
        <motion.div
          className="h-full bg-teal"
          initial={{ width: 0 }}
          animate={{ width: `${upPct}%` }}
          transition={{ duration: 0.5, ease: "easeOut" }}
        />
        <motion.div
          className="h-full bg-coral"
          initial={{ width: 0 }}
          animate={{ width: `${total === 0 ? 0 : 100 - upPct}%` }}
          transition={{ duration: 0.5, ease: "easeOut" }}
        />
      </div>
      {total > 4 && idea.downvotes > 0 && upPct < 70 && (
        <p className="mt-1.5 text-label-md text-slatey">
          Divisive — worth reading the thread before deciding.
        </p>
      )}
    </div>
  );
}

export function TagChip({
  tag,
  active,
  onClick,
}: {
  tag: string;
  active?: boolean;
  onClick?: () => void;
}) {
  const Comp = onClick ? "button" : "span";
  return (
    <Comp
      onClick={onClick}
      className={`rounded-full px-2 py-0.5 text-label-sm font-bold transition ${
        active
          ? "bg-purple text-white"
          : "bg-gray-50 text-slatey" + (onClick ? " hover:bg-gray-100" : "")
      }`}
    >
      #{tag}
    </Comp>
  );
}

/** Compact relative time — "3m", "2h", "5d", then a date. */
export function timeAgo(iso: string): string {
  const then = new Date(iso).getTime();
  if (!then) return "";
  const secs = Math.max(0, Math.floor((Date.now() - then) / 1000));
  if (secs < 60) return "just now";
  const mins = Math.floor(secs / 60);
  if (mins < 60) return `${mins}m`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d`;
  return new Date(iso).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
  });
}

/**
 * Renders message text with @mentions and @idea#N links highlighted. Clicking
 * an idea reference jumps to that thread, which is what makes cross-linking
 * worth doing at all.
 */
export function RichText({
  text,
  onOpenIdea,
}: {
  text: string;
  onOpenIdea?: (id: number) => void;
}) {
  // Idea references are matched first so "@idea#12" isn't half-consumed by the
  // plain @mention pattern.
  const pattern = /(@idea#\d+|@[A-Za-z][A-Za-z0-9_.-]{1,30})/g;
  const parts = text.split(pattern);

  return (
    <span className="whitespace-pre-wrap break-words">
      {parts.map((part, i) => {
        const ideaRef = /^@idea#(\d+)$/.exec(part);
        if (ideaRef) {
          const id = Number(ideaRef[1]);
          return (
            <button
              key={i}
              onClick={() => onOpenIdea?.(id)}
              className="rounded bg-teal-light px-1 font-bold text-teal hover:underline"
            >
              #{id}
            </button>
          );
        }
        if (/^@[A-Za-z]/.test(part)) {
          return (
            <span key={i} className="rounded bg-purple-light px-1 font-bold text-purple">
              {part}
            </span>
          );
        }
        return <span key={i}>{part}</span>;
      })}
    </span>
  );
}
