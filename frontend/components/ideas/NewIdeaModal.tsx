"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { AlertCircle, ArrowRight, Lightbulb } from "lucide-react";
import { Button } from "@/components/Button";
import { StatusBadge } from "./IdeaBits";
import { api } from "@/lib/api";
import type { SimilarIdea } from "@/lib/types";

/**
 * Posting an idea. Only the title is required — every extra mandatory field is
 * one more reason not to bother, and friction is what kills a board.
 *
 * The duplicate check runs while you type and only suggests; it never blocks.
 * A false "this already exists" is worse than a real duplicate, because it
 * teaches people to ignore the prompt.
 */
export function NewIdeaModal({
  onClose,
  onCreated,
  onOpenIdea,
}: {
  onClose: () => void;
  onCreated: (id: number) => void;
  onOpenIdea: (id: number) => void;
}) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [tags, setTags] = useState("");
  const [similar, setSimilar] = useState<SimilarIdea[]>([]);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  // Debounced so it isn't a request per keystroke.
  useEffect(() => {
    const text = `${title} ${description}`.trim();
    if (text.length < 8) {
      setSimilar([]);
      return;
    }
    const t = setTimeout(() => {
      api
        .similarIdeas(text)
        .then((r) => setSimilar(r.similar))
        .catch(() => setSimilar([]));
    }, 400);
    return () => clearTimeout(t);
  }, [title, description]);

  const submit = async () => {
    if (!title.trim() || busy) return;
    setBusy(true);
    setError("");
    try {
      const { idea } = await api.createIdea({
        title: title.trim(),
        description: description.trim(),
        tags: tags
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean),
      });
      onCreated(idea.id);
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not post that idea");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/50 p-4">
      <motion.div
        initial={{ opacity: 0, scale: 0.96, y: 10 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        className="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-2xl bg-white p-5 shadow-card-lg"
      >
        <div className="flex items-center gap-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-purple-light">
            <Lightbulb size={18} className="text-purple" />
          </div>
          <h3 className="text-heading-sm font-extrabold text-ink">New idea</h3>
        </div>

        <label className="mt-4 block text-label-md font-bold text-slatey">
          What&apos;s the idea?
        </label>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && submit()}
          autoFocus
          placeholder="Auto-tagging for notes"
          className="mt-1 w-full rounded-lg bg-gray-50 px-3 py-2.5 text-body-lg font-bold outline-none ring-purple/30 focus:ring-2"
        />

        {similar.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-3 rounded-lg bg-amber-light p-3"
          >
            <p className="flex items-center gap-1.5 text-label-lg font-extrabold text-ink">
              <AlertCircle size={14} className="text-amber" />
              Similar {similar.length === 1 ? "idea" : "ideas"} already on the board
            </p>
            <p className="mt-0.5 text-label-md text-slatey">
              Worth a look — adding your vote or a reply there often lands better
              than a second thread.
            </p>
            <div className="mt-2 space-y-1.5">
              {similar.map((s) => (
                <button
                  key={s.id}
                  onClick={() => onOpenIdea(s.id)}
                  className="flex w-full items-center gap-2 rounded-lg bg-white px-2.5 py-2 text-left transition hover:bg-gray-50"
                >
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-label-lg font-bold text-ink">
                      #{s.id} {s.title}
                    </p>
                    <p className="text-label-sm text-slatey">
                      {s.score} votes · {s.messageCount} messages ·{" "}
                      {Math.round(s.similarity * 100)}% overlap
                    </p>
                  </div>
                  <StatusBadge status={s.status} />
                  <ArrowRight size={13} className="shrink-0 text-slatey" />
                </button>
              ))}
            </div>
          </motion.div>
        )}

        <label className="mt-4 block text-label-md font-bold text-slatey">
          More detail <span className="font-normal">(optional)</span>
        </label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={4}
          placeholder="What problem does it solve? Rough is fine — the thread is where it gets sharpened."
          className="mt-1 w-full resize-none rounded-lg bg-gray-50 px-3 py-2 text-body-md outline-none ring-purple/30 focus:ring-2"
        />

        <label className="mt-3 block text-label-md font-bold text-slatey">
          Tags <span className="font-normal">(optional, comma separated)</span>
        </label>
        <input
          value={tags}
          onChange={(e) => setTags(e.target.value)}
          placeholder="ai, ux, v2.0"
          className="mt-1 w-full rounded-lg bg-gray-50 px-3 py-2 text-body-md outline-none ring-purple/30 focus:ring-2"
        />

        {error && (
          <p className="mt-3 rounded-lg bg-coral-light px-3 py-2 text-label-md text-coral">
            {error}
          </p>
        )}

        <div className="mt-5 flex gap-3">
          <Button variant="outline" className="flex-1" onClick={onClose}>
            Cancel
          </Button>
          <Button className="flex-1" onClick={submit} loading={busy}>
            Post idea
          </Button>
        </div>
      </motion.div>
    </div>
  );
}
