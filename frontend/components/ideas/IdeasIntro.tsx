"use client";

import { useEffect, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import {
  ArrowRight,
  CheckSquare,
  Lightbulb,
  MessagesSquare,
  ThumbsUp,
  UserX,
  X,
} from "lucide-react";
import { Button } from "@/components/Button";

/**
 * First-run guide for the ideas board.
 *
 * Shown once, then never again — the flag lives in localStorage rather than on
 * the account because it describes this browser's experience, and getting it
 * wrong in that direction (a returning user on a new device sees it once more)
 * is far cheaper than the opposite (someone who never saw it doesn't get it).
 *
 * It explains the four things that aren't guessable from looking at the board:
 * that every idea owns its own thread, that voting is one tap and needs no
 * justification, that status is a workflow rather than a label, and that a
 * popular idea escalates on its own.
 */

const STORAGE_KEY = "lumora_ideas_intro_seen";

const STEPS = [
  {
    icon: Lightbulb,
    tint: "#F5A623",
    title: "Post it rough",
    body: "A title is all you need — the thread is where an idea gets sharpened. If something similar already exists, you'll be shown it before you post rather than after.",
  },
  {
    icon: ThumbsUp,
    tint: "#00C2A8",
    title: "Vote in one tap",
    body: "No comment required to show support, and you can vote it down too. Both directions matter: an idea that's 30 up and 28 down is worth reading, and 'Controversial' sorting surfaces exactly those.",
  },
  {
    icon: MessagesSquare,
    tint: "#6C3FC5",
    title: "Every idea owns its thread",
    body: "Discussion stays attached to the idea it's about, so nothing gets lost in one long stream. Reply to branch a side-question, share a sketch, or record a voice memo.",
  },
  {
    icon: CheckSquare,
    tint: "#17A3DD",
    title: "Status is the progress",
    body: "Draft → Under review → Approved → In progress → Completed. Pass 20 votes and an idea moves to review on its own — nobody has to remember to escalate it.",
  },
];

/**
 * [visible, dismiss, replay] — replay lets the guide be reopened from the
 * toolbar, so dismissing it isn't a one-way door.
 */
export function useIdeasIntro(): [boolean, () => void, () => void] {
  const [show, setShow] = useState(false);

  useEffect(() => {
    try {
      setShow(window.localStorage.getItem(STORAGE_KEY) !== "1");
    } catch {
      // Private browsing with storage disabled: skip rather than nag on every
      // page load, which is the failure mode that would actually annoy people.
      setShow(false);
    }
  }, []);

  const dismiss = () => {
    setShow(false);
    try {
      window.localStorage.setItem(STORAGE_KEY, "1");
    } catch {
      /* nothing to do */
    }
  };

  return [show, dismiss, () => setShow(true)];
}

export function IdeasIntro({ onClose }: { onClose: () => void }) {
  const [step, setStep] = useState(0);
  const last = step === STEPS.length - 1;
  const current = STEPS[step];

  // Arrow keys and Escape, so it's navigable without reaching for the mouse.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
      if (e.key === "ArrowRight") setStep((s) => Math.min(s + 1, STEPS.length - 1));
      if (e.key === "ArrowLeft") setStep((s) => Math.max(s - 1, 0));
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [onClose]);

  return (
    <div className="fixed inset-0 z-[70] flex items-end justify-center bg-black/50 p-0 backdrop-blur-[2px] sm:items-center sm:p-4">
      <motion.div
        initial={{ opacity: 0, y: 40, scale: 0.98 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ type: "spring", stiffness: 300, damping: 30 }}
        className="w-full max-w-md overflow-hidden rounded-t-[28px] bg-white shadow-card-lg sm:rounded-[28px]"
        role="dialog"
        aria-modal="true"
        aria-label="How ideas work"
      >
        {/* Hero */}
        <div
          className="relative px-6 pb-6 pt-7 text-white transition-colors duration-500"
          style={{
            background: `linear-gradient(150deg, ${current.tint}, ${current.tint}bb)`,
          }}
        >
          <button
            onClick={onClose}
            aria-label="Skip"
            className="absolute right-4 top-4 rounded-full bg-white/20 p-1.5 text-white transition hover:bg-white/30"
          >
            <X size={16} />
          </button>

          <p className="text-label-md font-bold uppercase tracking-[0.2em] text-white/70">
            How Ideas work
          </p>

          <AnimatePresence mode="wait">
            <motion.div
              key={step}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -12 }}
              transition={{ duration: 0.25 }}
              className="mt-4 flex items-center gap-3"
            >
              <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-white/20 backdrop-blur">
                <current.icon size={24} />
              </span>
              <h2 className="text-heading-lg font-extrabold leading-tight">
                {current.title}
              </h2>
            </motion.div>
          </AnimatePresence>
        </div>

        {/* Body */}
        <div className="px-6 pb-[max(1.25rem,env(safe-area-inset-bottom))] pt-5">
          <AnimatePresence mode="wait">
            <motion.p
              key={step}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="min-h-[88px] text-body-md text-slatey"
            >
              {current.body}
            </motion.p>
          </AnimatePresence>

          {last && (
            <motion.p
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-1 flex items-start gap-2 rounded-2xl bg-gray-50 px-3 py-2.5 text-label-lg text-slatey"
            >
              <UserX size={15} className="mt-0.5 shrink-0 text-ink" />
              <span>
                One more: <strong className="text-ink">Silent brainstorm</strong>{" "}
                opens a timed window where every post is anonymous, so ideas are
                judged on content rather than on who said them.
              </span>
            </motion.p>
          )}

          {/* Progress */}
          <div className="mt-5 flex items-center justify-center gap-1.5">
            {STEPS.map((s, i) => (
              <button
                key={s.title}
                onClick={() => setStep(i)}
                aria-label={`Step ${i + 1}: ${s.title}`}
                className="h-1.5 rounded-full transition-all"
                style={{
                  width: i === step ? 22 : 7,
                  background: i === step ? current.tint : "#EBEBEB",
                }}
              />
            ))}
          </div>

          <div className="mt-5 flex gap-3">
            {!last ? (
              <>
                <Button variant="ghost" className="flex-1" onClick={onClose}>
                  Skip
                </Button>
                <Button className="flex-1" onClick={() => setStep((s) => s + 1)}>
                  Next <ArrowRight size={17} className="ml-1.5" />
                </Button>
              </>
            ) : (
              <Button full onClick={onClose}>
                Start posting ideas
              </Button>
            )}
          </div>
        </div>
      </motion.div>
    </div>
  );
}
