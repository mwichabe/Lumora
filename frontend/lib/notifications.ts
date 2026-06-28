"use client";

import { useCallback, useEffect, useState } from "react";
import { api } from "./api";
import type { AppNotification } from "./types";

/**
 * Notifications are delivered by the backend (welcome message, tips, feature and
 * language announcements pushed by the campaign scheduler). The frontend polls
 * for new ones and broadcasts a window event so the bell, sidebar and screen
 * stay in sync.
 */
const EVENT = "lumora:notifications";
const POLL_MS = 60_000;

function ping() {
  if (typeof window !== "undefined") window.dispatchEvent(new Event(EVENT));
}

export function useNotifications() {
  const [items, setItems] = useState<AppNotification[]>([]);
  const [unread, setUnread] = useState(0);

  const reload = useCallback(async () => {
    try {
      const d = await api.notifications();
      setItems(d.notifications);
      setUnread(d.unread);
    } catch {
      /* offline / unauthenticated — leave as-is */
    }
  }, []);

  useEffect(() => {
    reload();
    const id = setInterval(reload, POLL_MS);
    const onPing = () => reload();
    window.addEventListener(EVENT, onPing);
    return () => {
      clearInterval(id);
      window.removeEventListener(EVENT, onPing);
    };
  }, [reload]);

  const markAllRead = useCallback(async () => {
    try {
      await api.markNotificationsRead();
      setItems((prev) => prev.map((n) => ({ ...n, read: true })));
      setUnread(0);
      ping();
    } catch {
      /* ignore */
    }
  }, []);

  return { items, unread, reload, markAllRead };
}

/** Lightweight unread-count hook for the bell / sidebar badges. */
export function useUnreadCount(): number {
  const { unread } = useNotifications();
  return unread;
}
