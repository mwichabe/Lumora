"use client";

import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Languages, Loader2, RefreshCw } from "lucide-react";
import type { MessageTranslation } from "@/lib/types";

/**
 * Shows a foreign-language message with its English translation.
 *
 * Which text leads is the whole design question here, and in a
 * language-learning app the answer isn't obvious. The translation leads,
 * because the person reading it is usually trying to follow a conversation
 * rather than study — but the original is always one tap away and clearly
 * labelled, never hidden or replaced. Nothing renders at all for English
 * messages; the component isn't even mounted.
 *
 * Three states, all of which really happen:
 *   translated — show English, offer the original
 *   pending    — the background pass is still running; show the original and say so
 *   unavailable— translation is off or failed; show the original and offer a retry
 */
export function TranslatedText({
  body,
  translation,
  onRetry,
  tone = "light",
  children,
}: {
  body: string;
  translation?: MessageTranslation;
  /** Retry hook; omit to hide the retry affordance entirely. */
  onRetry?: () => Promise<void>;
  /** "dark" for use inside a coloured message bubble. */
  tone?: "light" | "dark";
  /** Renders the message text — lets callers keep their own @mention linking. */
  children: (text: string) => React.ReactNode;
}) {
  const [showOriginal, setShowOriginal] = useState(false);
  const [retrying, setRetrying] = useState(false);

  // No translation block means the message is English (or too short to judge)
  // — render it exactly as any other message, with no chrome.
  if (!translation) return <>{children(body)}</>;

  const hasTranslation = translation.text !== "";
  const showing = hasTranslation && !showOriginal ? translation.text : body;

  const muted = tone === "dark" ? "text-white/70" : "text-slatey";
  const accent = tone === "dark" ? "text-white/90" : "text-purple";

  const retry = async () => {
    if (!onRetry || retrying) return;
    setRetrying(true);
    try {
      await onRetry();
    } finally {
      setRetrying(false);
    }
  };

  return (
    <>
      <AnimatePresence mode="wait" initial={false}>
        <motion.span
          key={showing === body ? "original" : "translated"}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.12 }}
          className="block"
        >
          {children(showing)}
        </motion.span>
      </AnimatePresence>

      <span
        className={`mt-1 flex flex-wrap items-center gap-x-2 gap-y-0.5 text-label-sm ${muted}`}
      >
        <span className="flex items-center gap-1">
          <Languages size={11} className="shrink-0" />
          {hasTranslation && !showOriginal
            ? `Translated from ${translation.langName}`
            : translation.langName}
        </span>

        {hasTranslation && (
          <button
            onClick={() => setShowOriginal((o) => !o)}
            className={`font-bold underline decoration-dotted underline-offset-2 ${accent}`}
          >
            {showOriginal ? "Show translation" : "Show original"}
          </button>
        )}

        {/* The background pass is still running — the message is already
            readable in its own language, so this is information, not a wait. */}
        {!hasTranslation && translation.pending && (
          <span className="flex items-center gap-1">
            <Loader2 size={10} className="animate-spin" /> translating…
          </span>
        )}

        {/* Translation failed or is switched off. Offer a retry rather than
            leaving a dead label. */}
        {!hasTranslation && !translation.pending && onRetry && (
          <button
            onClick={retry}
            disabled={retrying}
            className={`flex items-center gap-1 font-bold underline decoration-dotted underline-offset-2 disabled:opacity-50 ${accent}`}
          >
            {retrying ? (
              <Loader2 size={10} className="animate-spin" />
            ) : (
              <RefreshCw size={10} />
            )}
            Translate
          </button>
        )}
      </span>
    </>
  );
}
