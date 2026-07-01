"use client";

import { motion } from "framer-motion";
import { Heart, Clock, X } from "lucide-react";
import { Button } from "./Button";
import { fmtCountdown } from "@/lib/hearts";
import type { HeartsStatus } from "@/lib/types";

/**
 * Shown when the learner runs out of hearts: wait for regeneration, or refill
 * instantly with a purchase.
 */
export function OutOfHeartsModal({
  status,
  secondsToNext,
  note,
  closeLabel = "I'll wait",
  onClose,
  onBuy,
  buying,
}: {
  status: HeartsStatus | null;
  secondsToNext: number;
  note?: string;
  closeLabel?: string;
  onClose: () => void;
  onBuy: () => void;
  buying: boolean;
}) {
  if (!status) return null;
  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center bg-black/50 p-0 sm:items-center sm:p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ y: 40, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        onClick={(e) => e.stopPropagation()}
        className="w-full max-w-sm overflow-hidden rounded-t-3xl bg-white shadow-card-lg sm:rounded-3xl"
      >
        <div className="relative flex flex-col items-center bg-coral/10 px-6 pb-5 pt-8">
          <button
            onClick={onClose}
            aria-label="Close"
            className="absolute right-4 top-4 flex h-8 w-8 items-center justify-center rounded-full bg-white/70 text-gray-500 transition hover:bg-white"
          >
            <X size={18} />
          </button>
          <span className="flex h-16 w-16 items-center justify-center rounded-full bg-coral/15 text-coral">
            <Heart size={34} fill="currentColor" />
          </span>
          <h2 className="mt-3 text-heading-lg font-extrabold text-ink">
            You&apos;re out of hearts
          </h2>
          <p className="mt-1 text-center text-body-sm text-slatey">
            {note || "Hearts refill over time — or top up now to keep learning."}
          </p>
        </div>

        <div className="px-6 py-5">
          {/* countdown */}
          <div className="flex items-center justify-center gap-2 rounded-2xl bg-gray-50 py-3 text-ink">
            <Clock size={18} className="text-slatey" />
            {secondsToNext > 0 ? (
              <span className="text-body-md font-bold">
                Next heart in {fmtCountdown(secondsToNext)}
              </span>
            ) : (
              <span className="text-body-md font-bold">A heart is on its way…</span>
            )}
          </div>
          <p className="mt-1 text-center text-label-md text-slatey">
            One heart every {status.regenMinutes} minutes
          </p>

          <div className="mt-4 space-y-2">
            {status.paymentsEnabled && (
              <Button full disabled={buying} onClick={onBuy}>
                {buying
                  ? "Redirecting…"
                  : `Refill ${status.max} hearts — KES ${status.refillPriceKes}` +
                    (status.refillPriceUsd > 0
                      ? ` (≈ $${status.refillPriceUsd.toFixed(2)})`
                      : "")}
              </Button>
            )}
            <Button full variant="outline" onClick={onClose}>
              {closeLabel}
            </Button>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
