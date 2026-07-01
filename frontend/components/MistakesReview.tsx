"use client";

import { Volume2, RotateCcw } from "lucide-react";
import { Button } from "./Button";
import { speakAs } from "@/lib/voices";

export interface ReviewItem {
  prompt?: string;
  question: string;
  correctAnswer: string;
  /** Target-language text to pronounce, when known (omit to hide the play button). */
  playText?: string;
  speaker?: string;
}

/**
 * An end-of-session recap of the items the learner got wrong, so they can study
 * them before finishing. Renders as page content (the caller provides the shell).
 */
export function MistakesReview({
  items,
  onDone,
  finishLabel = "Finish",
}: {
  items: ReviewItem[];
  onDone: () => void;
  finishLabel?: string;
}) {
  return (
    <div className="flex flex-1 flex-col">
      <div className="flex items-center gap-2">
        <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-light text-amber">
          <RotateCcw size={18} />
        </span>
        <div>
          <h2 className="text-heading-md font-extrabold text-ink">
            Review your mistakes
          </h2>
          <p className="text-body-sm text-slatey">
            Go over these {items.length === 1 ? "one" : items.length} before you
            finish.
          </p>
        </div>
      </div>

      <div className="mt-4 flex-1 space-y-2.5 overflow-y-auto">
        {items.map((it, i) => (
          <div
            key={i}
            className="rounded-2xl border border-gray-100 bg-white p-4 shadow-card"
          >
            {it.prompt && (
              <p className="text-label-sm font-bold uppercase tracking-wide text-gray-400">
                {it.prompt}
              </p>
            )}
            <div className="mt-0.5 flex items-start justify-between gap-3">
              <p className="text-body-lg font-extrabold text-ink">
                {it.question}
              </p>
              {it.playText && (
                <button
                  onClick={() => speakAs(it.speaker, it.playText!)}
                  aria-label="Hear it"
                  className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-purple text-white"
                >
                  <Volume2 size={16} />
                </button>
              )}
            </div>
            <div className="mt-2 rounded-lg bg-teal-light px-3 py-2">
              <p className="text-label-sm font-bold uppercase tracking-wide text-teal">
                Correct answer
              </p>
              <p className="text-body-md font-extrabold text-ink">
                {it.correctAnswer}
              </p>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-4 pt-2">
        <Button full onClick={onDone}>
          {finishLabel}
        </Button>
      </div>
    </div>
  );
}
