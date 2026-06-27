"use client";

import { useEffect, useState } from "react";
import type { User } from "./types";

/**
 * Notifications are generated client-side from the user's current state (there
 * is no backend feed yet). Read-state is persisted in localStorage and a custom
 * window event keeps the bell badge and the screen in sync within a session.
 */
export interface AppNotification {
  id: string;
  emoji: string;
  tint: string;
  title: string;
  body: string;
  time: string;
  read: boolean;
}

const READ_KEY = "lumora_notifications_read";
const EVENT = "lumora:notifications";

function readSet(): Set<string> {
  if (typeof window === "undefined") return new Set();
  try {
    return new Set(JSON.parse(window.localStorage.getItem(READ_KEY) || "[]"));
  } catch {
    return new Set();
  }
}

function writeSet(s: Set<string>) {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(READ_KEY, JSON.stringify([...s]));
  window.dispatchEvent(new Event(EVENT));
}

function templates(user: User | null): Omit<AppNotification, "read">[] {
  const streak = user?.streak ?? 0;
  const toGoal = Math.max(0, (user?.dailyGoalXp ?? 50) - (user?.xpToday ?? 0));
  const league = user?.league || "Bronze";
  const name = user?.name?.split(" ")[0] || "friend";

  return [
    {
      id: "streak",
      emoji: "🔥",
      tint: "#FF5C5C",
      title: streak > 0 ? `${streak}-day streak!` : "Start your streak today",
      body:
        streak > 0
          ? "Practise today to keep your flame burning bright."
          : "Finish one lesson to light your first flame.",
      time: "2h ago",
    },
    {
      id: "goal",
      emoji: "🎯",
      tint: "#F5A623",
      title: "Daily goal",
      body:
        toGoal > 0
          ? `You're just ${toGoal} XP away from today's goal, ${name}!`
          : "You smashed today's goal. Incredible work! 🎉",
      time: "5h ago",
    },
    {
      id: "league",
      emoji: "🏆",
      tint: "#6C3FC5",
      title: `${league} League`,
      body: "The weekly race is heating up — climb into the top 3 to advance!",
      time: "1d ago",
    },
    {
      id: "quest",
      emoji: "🦔",
      tint: "#00C2A8",
      title: "Pip has new quests",
      body: "Three fresh daily quests are waiting for you. Let's GO!",
      time: "1d ago",
    },
    {
      id: "feature",
      emoji: "✨",
      tint: "#7B4AD6",
      title: "Coming soon: AI Conversations",
      body: "Chat with Blaze to practise speaking in real scenarios. Stay tuned!",
      time: "3d ago",
    },
  ];
}

export function buildNotifications(user: User | null): AppNotification[] {
  const read = readSet();
  return templates(user).map((t) => ({ ...t, read: read.has(t.id) }));
}

export function markRead(id: string) {
  const s = readSet();
  s.add(id);
  writeSet(s);
}

export function markAllRead(ids: string[]) {
  const s = readSet();
  ids.forEach((id) => s.add(id));
  writeSet(s);
}

/** Reactive unread count for the bell badge. */
export function useUnreadCount(user: User | null): number {
  const [count, setCount] = useState(0);

  useEffect(() => {
    const update = () =>
      setCount(buildNotifications(user).filter((n) => !n.read).length);
    update();
    window.addEventListener(EVENT, update);
    window.addEventListener("storage", update);
    return () => {
      window.removeEventListener(EVENT, update);
      window.removeEventListener("storage", update);
    };
  }, [user]);

  return count;
}
