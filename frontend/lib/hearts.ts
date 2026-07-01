"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { api } from "./api";
import type { HeartsStatus } from "./types";

/**
 * Hearts are server-authoritative: they're spent on wrong answers, regenerate
 * over time, and can be purchased. This hook keeps a local copy in sync, runs a
 * live countdown to the next regenerated heart, and broadcasts changes so every
 * mounted badge/screen stays consistent.
 */
const EVENT = "lumora:hearts";
const POLL_MS = 60_000;

function ping() {
  if (typeof window !== "undefined") window.dispatchEvent(new Event(EVENT));
}

export function useHearts() {
  const [status, setStatus] = useState<HeartsStatus | null>(null);
  const [secondsToNext, setSecondsToNext] = useState(0);
  const tick = useRef<ReturnType<typeof setInterval> | null>(null);

  const apply = useCallback((s: HeartsStatus) => {
    setStatus(s);
    setSecondsToNext(s.secondsToNext);
  }, []);

  const reload = useCallback(async () => {
    try {
      apply(await api.heartsStatus());
    } catch {
      /* offline / unauthenticated */
    }
  }, [apply]);

  // Spend a heart (on a wrong answer) — returns the fresh status.
  const lose = useCallback(async (): Promise<HeartsStatus | null> => {
    try {
      const s = await api.loseHeart();
      apply(s);
      ping();
      return s;
    } catch {
      return null;
    }
  }, [apply]);

  // Start Paystack checkout to refill hearts.
  const buy = useCallback(async () => {
    try {
      const r = await api.buyHearts();
      if (r.authorizationUrl) window.location.href = r.authorizationUrl;
    } catch {
      /* ignore */
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

  // Local 1s countdown; when it hits zero, refetch to collect the new heart.
  useEffect(() => {
    if (tick.current) clearInterval(tick.current);
    if (!status || status.full || secondsToNext <= 0) return;
    tick.current = setInterval(() => {
      setSecondsToNext((s) => {
        if (s <= 1) {
          reload();
          return 0;
        }
        return s - 1;
      });
    }, 1000);
    return () => {
      if (tick.current) clearInterval(tick.current);
    };
  }, [status, secondsToNext, reload]);

  return {
    status,
    hearts: status?.hearts ?? 0,
    max: status?.max ?? 5,
    full: status?.full ?? true,
    secondsToNext,
    reload,
    lose,
    buy,
  };
}

/** mm:ss for a hearts countdown. */
export function fmtCountdown(s: number): string {
  const m = Math.floor(s / 60);
  const sec = s % 60;
  return `${m}:${String(sec).padStart(2, "0")}`;
}
