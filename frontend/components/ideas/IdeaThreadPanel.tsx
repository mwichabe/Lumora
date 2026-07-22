"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import {
  Check,
  CircleHelp,
  Code2,
  CornerDownRight,
  Image as ImageIcon,
  Lightbulb,
  Loader2,
  Mic,
  Pencil,
  Send,
  SmilePlus,
  Sparkles,
  Square,
  Trash2,
  UserX,
  X,
} from "lucide-react";
import {
  ActionMenu,
  ActionMenuItem,
  ActionMenuNote,
} from "@/components/ActionMenu";
import { Avatar } from "@/components/Avatar";
import { TranslatedText } from "@/components/TranslatedText";
import { ImageEditorModal } from "@/components/ImageEditorModal";
import { RichText, StatusBadge, timeAgo } from "./IdeaBits";
import { api, mediaUrl } from "@/lib/api";
import type {
  BrainstormSession,
  ChatUser,
  IdeaMessage,
  IdeaThread,
  ThreadSummary,
} from "@/lib/types";

/**
 * The centre panel: one thread per idea.
 *
 * Threading is two levels deep and no more. Anything deeper reads as a tree
 * nobody can follow, and the point of thread-per-idea is that context stays
 * attached to the idea rather than scattered through one chaotic stream.
 */

export function IdeaThreadPanel({
  thread,
  loading,
  contacts,
  onRefresh,
  onOpenIdea,
  onBack,
}: {
  thread: IdeaThread | null;
  loading: boolean;
  contacts: ChatUser[];
  onRefresh: () => void;
  onOpenIdea: (id: number) => void;
  onBack?: () => void;
}) {
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const [replyTo, setReplyTo] = useState<IdeaMessage | null>(null);
  const [summary, setSummary] = useState<ThreadSummary | null>(null);
  const [summaryBusy, setSummaryBusy] = useState(false);

  const idea = thread?.idea;
  const brainstorm = thread?.brainstorm ?? null;

  useEffect(() => {
    setReplyTo(null);
    setSummary(null);
  }, [idea?.id]);

  useEffect(() => {
    const el = scrollRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [thread?.messages.length, idea?.id]);

  const loadSummary = async () => {
    if (!idea) return;
    setSummaryBusy(true);
    try {
      setSummary(await api.ideaSummary(idea.id));
    } finally {
      setSummaryBusy(false);
    }
  };

  if (!idea) {
    return (
      <div className="flex h-full items-center justify-center bg-cream p-8 text-center">
        <div>
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-purple-light">
            <Lightbulb size={24} className="text-purple" />
          </div>
          <p className="mt-3 text-heading-sm font-extrabold text-ink">
            Pick an idea
          </p>
          <p className="mt-1 max-w-xs text-body-sm text-slatey">
            Every idea has its own thread, so the discussion stays attached to
            what it&apos;s about.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full min-w-0 flex-col bg-cream">
      {/* Header */}
      <div className="shrink-0 border-b border-gray-100 bg-white px-4 py-3">
        <div className="flex items-start gap-2">
          {onBack && (
            <button
              onClick={onBack}
              className="rounded-md p-1 text-slatey lg:hidden"
              aria-label="Back to ideas"
            >
              <X size={18} />
            </button>
          )}
          <div className="min-w-0 flex-1">
            <p className="truncate font-extrabold text-ink">
              <span className="text-slatey">#{idea.id}</span> {idea.title}
            </p>
            <div className="mt-1 flex items-center gap-2">
              <StatusBadge status={idea.status} />
              <span className="text-label-md text-slatey">
                {idea.messageCount} messages
              </span>
            </div>
          </div>
          <button
            onClick={loadSummary}
            disabled={summaryBusy || idea.messageCount === 0}
            title="Summarise this thread"
            className="flex shrink-0 items-center gap-1 rounded-full bg-gray-50 px-3 py-1.5 text-label-md font-bold text-slatey transition hover:bg-gray-100 hover:text-ink disabled:opacity-40"
          >
            {summaryBusy ? (
              <Loader2 size={13} className="animate-spin" />
            ) : (
              <Sparkles size={13} />
            )}
            Summarise
          </button>
        </div>
      </div>

      {brainstorm && <BrainstormBanner session={brainstorm} />}

      <AnimatePresence>
        {summary && (
          <SummaryCard summary={summary} onClose={() => setSummary(null)} />
        )}
      </AnimatePresence>

      {/* Messages */}
      <div ref={scrollRef} className="min-h-0 flex-1 space-y-3 overflow-y-auto p-4">
        {loading && (
          <div className="space-y-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="h-16 animate-pulse rounded-2xl bg-gray-100" />
            ))}
          </div>
        )}

        {!loading && thread?.messages.length === 0 && (
          <div className="py-10 text-center">
            <p className="text-body-md font-bold text-ink">No discussion yet</p>
            <p className="mt-1 text-body-sm text-slatey">
              Say what you think, sketch it, or ask the awkward question.
            </p>
          </div>
        )}

        {thread?.messages.map((m) => (
          <MessageBlock
            key={m.id}
            message={m}
            palette={thread.reactions}
            onRefresh={onRefresh}
            onReply={() => setReplyTo(m)}
            onOpenIdea={onOpenIdea}
          />
        ))}
      </div>

      <Composer
        ideaId={idea.id}
        replyTo={replyTo}
        contacts={contacts}
        anonymous={!!brainstorm}
        onCancelReply={() => setReplyTo(null)}
        onSent={() => {
          setReplyTo(null);
          onRefresh();
        }}
      />
    </div>
  );
}

// --- banners -----------------------------------------------------------------

function BrainstormBanner({ session }: { session: BrainstormSession }) {
  const [left, setLeft] = useState(session.secondsRemaining);

  useEffect(() => setLeft(session.secondsRemaining), [session.secondsRemaining]);
  useEffect(() => {
    const t = setInterval(() => setLeft((s) => Math.max(0, s - 1)), 1000);
    return () => clearInterval(t);
  }, []);

  const mm = Math.floor(left / 60);
  const ss = String(left % 60).padStart(2, "0");

  return (
    <div className="flex shrink-0 items-center gap-2 bg-ink px-4 py-2 text-white">
      <UserX size={15} />
      <span className="text-label-lg font-extrabold">Silent brainstorm</span>
      <span className="text-label-md text-white/70">
        Everything you post is anonymous
        {session.topic ? ` — ${session.topic}` : ""}
      </span>
      <span className="ml-auto font-mono text-label-lg font-bold">
        {mm}:{ss}
      </span>
    </div>
  );
}

function SummaryCard({
  summary,
  onClose,
}: {
  summary: ThreadSummary;
  onClose: () => void;
}) {
  return (
    <motion.div
      initial={{ height: 0, opacity: 0 }}
      animate={{ height: "auto", opacity: 1 }}
      exit={{ height: 0, opacity: 0 }}
      className="shrink-0 overflow-hidden border-b border-gray-100 bg-purple-light"
    >
      <div className="p-4">
        <div className="flex items-start justify-between gap-2">
          <p className="flex items-center gap-1.5 text-label-lg font-extrabold text-purple">
            <Sparkles size={14} /> Thread summary
          </p>
          <button onClick={onClose} aria-label="Close summary" className="text-purple">
            <X size={15} />
          </button>
        </div>

        <p className="mt-2 text-body-sm text-ink">{summary.gist}</p>

        {summary.keyPoints.length > 0 && (
          <ul className="mt-3 space-y-1.5">
            {summary.keyPoints.map((k, i) => (
              <li key={i} className="flex gap-2 text-body-sm text-slatey">
                <span className="text-purple">•</span>
                <span>
                  <strong className="text-ink">{k.author.name}:</strong> {k.text}
                </span>
              </li>
            ))}
          </ul>
        )}

        {summary.questions.length > 0 && (
          <div className="mt-3">
            <p className="text-label-md font-extrabold text-coral">Open questions</p>
            <ul className="mt-1 space-y-1">
              {summary.questions.map((q, i) => (
                <li key={i} className="flex gap-1.5 text-body-sm text-slatey">
                  <CircleHelp size={13} className="mt-0.5 shrink-0 text-coral" />
                  <span>{q.text}</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Said plainly, because it changes how much you should trust it. */}
        <p className="mt-3 text-label-sm text-slatey">
          Extracted from what people actually wrote — every line above is a real
          quote, not a paraphrase.
        </p>
      </div>
    </motion.div>
  );
}

// --- messages ----------------------------------------------------------------

function MessageBlock({
  message,
  palette,
  onRefresh,
  onReply,
  onOpenIdea,
}: {
  message: IdeaMessage;
  palette: string[];
  onRefresh: () => void;
  onReply: () => void;
  onOpenIdea: (id: number) => void;
}) {
  const [showReplies, setShowReplies] = useState(true);

  return (
    <div>
      <MessageRow
        message={message}
        palette={palette}
        onRefresh={onRefresh}
        onReply={onReply}
        onOpenIdea={onOpenIdea}
      />

      {message.replies.length > 0 && (
        <div className="ml-6 mt-2 border-l-2 border-gray-100 pl-3">
          <button
            onClick={() => setShowReplies((s) => !s)}
            className="mb-2 flex items-center gap-1 text-label-md font-bold text-purple hover:underline"
          >
            <CornerDownRight size={12} />
            {showReplies ? "Hide" : "Show"} thread · {message.replies.length}{" "}
            {message.replies.length === 1 ? "reply" : "replies"}
          </button>
          {showReplies && (
            <div className="space-y-2">
              {message.replies.map((r) => (
                <MessageRow
                  key={r.id}
                  message={r}
                  palette={palette}
                  compact
                  onRefresh={onRefresh}
                  onOpenIdea={onOpenIdea}
                />
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function MessageRow({
  message,
  palette,
  compact,
  onRefresh,
  onReply,
  onOpenIdea,
}: {
  message: IdeaMessage;
  palette: string[];
  compact?: boolean;
  onRefresh: () => void;
  onReply?: () => void;
  onOpenIdea: (id: number) => void;
}) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(message.body);
  const [busy, setBusy] = useState(false);
  const [reactOpen, setReactOpen] = useState(false);

  // A deleted message keeps its place in the thread. Removing it outright would
  // silently rewrite a conversation other people replied to.
  if (message.deleted) {
    return (
      <div className="flex items-center gap-2 rounded-2xl bg-gray-50 px-3 py-2 text-body-sm italic text-gray-500">
        <Trash2 size={13} /> This message was deleted
      </div>
    );
  }

  const save = async () => {
    const body = draft.trim();
    if (!body || body === message.body) return setEditing(false);
    setBusy(true);
    try {
      await api.editIdeaMessage(message.id, body);
      setEditing(false);
      onRefresh();
    } finally {
      setBusy(false);
    }
  };

  const remove = async () => {
    setBusy(true);
    try {
      await api.deleteIdeaMessage(message.id);
      onRefresh();
    } finally {
      setBusy(false);
    }
  };

  // Retry path — the server translates in the background when a message is
  // posted, so this only fires when that didn't land.
  const translate = async () => {
    await api.translateIdeaMessage(message.id).catch(() => {});
    onRefresh();
  };

  const react = async (emoji: string) => {
    setReactOpen(false);
    await api.reactToIdeaMessage(message.id, emoji).catch(() => {});
    onRefresh();
  };

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 6 }}
      animate={{ opacity: 1, y: 0 }}
      className="group relative flex gap-2.5"
    >
      {message.anonymous ? (
        <div
          className="flex shrink-0 items-center justify-center rounded-full bg-ink text-white"
          style={{ width: compact ? 28 : 34, height: compact ? 28 : 34 }}
          title="Posted anonymously during a silent brainstorm"
        >
          <UserX size={compact ? 13 : 15} />
        </div>
      ) : (
        <Avatar
          name={message.author?.name}
          color={message.author?.avatarColor}
          url={message.author?.avatarUrl}
          size={compact ? 28 : 34}
        />
      )}

      <div className="min-w-0 flex-1">
        <div className="flex items-baseline gap-2">
          <span className="text-label-lg font-extrabold text-ink">
            {message.anonymous ? "Anonymous" : message.author?.name ?? "Someone"}
          </span>
          <span className="text-label-sm text-slatey">
            {timeAgo(message.createdAt)}
          </span>
          {message.edited && (
            <span className="text-label-sm italic text-gray-500">edited</span>
          )}
        </div>

        <div className="mt-0.5">
          {editing ? (
            <div className="rounded-2xl bg-white p-2 shadow-card">
              <textarea
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                rows={3}
                autoFocus
                className="w-full resize-none rounded-lg bg-gray-50 p-2 text-body-md text-ink outline-none ring-purple/30 focus:ring-2"
              />
              <div className="mt-2 flex justify-end gap-2">
                <button
                  onClick={() => {
                    setDraft(message.body);
                    setEditing(false);
                  }}
                  className="rounded-full px-3 py-1 text-label-lg font-bold text-slatey hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  onClick={save}
                  disabled={busy}
                  className="flex items-center gap-1 rounded-full bg-purple px-3 py-1 text-label-lg font-extrabold text-white disabled:opacity-40"
                >
                  <Check size={13} /> Save
                </button>
              </div>
            </div>
          ) : (
            <MessageBody
              message={message}
              onOpenIdea={onOpenIdea}
              onTranslate={translate}
            />
          )}
        </div>

        {/* Reactions */}
        <div className="mt-1.5 flex flex-wrap items-center gap-1">
          {message.reactions.map((r) => (
            <button
              key={r.emoji}
              onClick={() => react(r.emoji)}
              className={`flex items-center gap-1 rounded-full px-2 py-0.5 text-label-md font-bold transition ${
                r.mine
                  ? "bg-purple-light text-purple ring-1 ring-purple"
                  : "bg-gray-50 text-slatey hover:bg-gray-100"
              }`}
            >
              {r.emoji} {r.count}
            </button>
          ))}

          <div className="relative opacity-0 transition group-hover:opacity-100 focus-within:opacity-100">
            <button
              onClick={() => setReactOpen((o) => !o)}
              aria-label="Add reaction"
              className="flex items-center rounded-full bg-gray-50 px-2 py-1 text-slatey transition hover:bg-gray-100 hover:text-ink"
            >
              <SmilePlus size={13} />
            </button>
            <AnimatePresence>
              {reactOpen && (
                <motion.div
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.9 }}
                  className="absolute bottom-full left-0 z-20 mb-1 flex gap-1 rounded-full bg-white p-1.5 shadow-card-lg"
                >
                  {palette.map((e) => (
                    <button
                      key={e}
                      onClick={() => react(e)}
                      className="rounded-full px-1.5 py-0.5 text-body-md transition hover:bg-gray-50"
                    >
                      {e}
                    </button>
                  ))}
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {onReply && (
            <button
              onClick={onReply}
              className="rounded-full px-2 py-0.5 text-label-md font-bold text-slatey opacity-0 transition hover:bg-gray-50 hover:text-ink group-hover:opacity-100"
            >
              Reply
            </button>
          )}
        </div>
      </div>

      {/* Own-message actions */}
      {(message.canEdit || message.mine) && (
        <div className="shrink-0">
          <ActionMenu>
            {(close) => (
              <>
                {message.canEdit ? (
                  <ActionMenuItem
                    icon={<Pencil size={13} />}
                    onClick={() => {
                      setEditing(true);
                      close();
                    }}
                  >
                    Edit
                  </ActionMenuItem>
                ) : (
                  <ActionMenuNote>
                    {message.kind === "text" || message.kind === "code"
                      ? "Edits close after 24 hours."
                      : "Attachments can't be edited — delete and repost."}
                  </ActionMenuNote>
                )}
                <ActionMenuItem
                  icon={<Trash2 size={13} />}
                  tone="danger"
                  disabled={busy}
                  onClick={() => {
                    remove();
                    close();
                  }}
                >
                  Delete
                </ActionMenuItem>
              </>
            )}
          </ActionMenu>
        </div>
      )}
    </motion.div>
  );
}

function MessageBody({
  message,
  onOpenIdea,
  onTranslate,
}: {
  message: IdeaMessage;
  onOpenIdea: (id: number) => void;
  onTranslate: () => Promise<void>;
}) {
  if (message.kind === "image") {
    return (
      <div>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={mediaUrl(message.url)}
          alt={message.fileName || "Shared image"}
          width={message.width || undefined}
          height={message.height || undefined}
          className="max-h-80 w-auto max-w-full rounded-2xl border border-gray-100 object-contain"
        />
        {message.body && (
          <p className="mt-1 text-body-md text-ink">
            <TranslatedText body={message.body} translation={message.translation} onRetry={onTranslate}>
              {(text) => <RichText text={text} onOpenIdea={onOpenIdea} />}
            </TranslatedText>
          </p>
        )}
      </div>
    );
  }

  if (message.kind === "voice") {
    return (
      <div className="inline-flex items-center gap-2 rounded-2xl bg-white p-2 shadow-card">
        <Mic size={15} className="shrink-0 text-purple" />
        {/* eslint-disable-next-line jsx-a11y/media-has-caption */}
        <audio controls src={mediaUrl(message.url)} className="h-8 max-w-[220px]" />
        {message.duration > 0 && (
          <span className="text-label-md text-slatey">{message.duration}s</span>
        )}
      </div>
    );
  }

  if (message.kind === "code") {
    return (
      <pre className="overflow-x-auto rounded-lg bg-ink p-3 text-label-lg text-white">
        <code>{message.body}</code>
      </pre>
    );
  }

  return (
    <p className="text-body-md text-ink">
      <TranslatedText body={message.body} translation={message.translation} onRetry={onTranslate}>
        {(text) => <RichText text={text} onOpenIdea={onOpenIdea} />}
      </TranslatedText>
    </p>
  );
}

// --- composer ----------------------------------------------------------------

function Composer({
  ideaId,
  replyTo,
  contacts,
  anonymous,
  onCancelReply,
  onSent,
}: {
  ideaId: number;
  replyTo: IdeaMessage | null;
  contacts: ChatUser[];
  anonymous: boolean;
  onCancelReply: () => void;
  onSent: () => void;
}) {
  const [body, setBody] = useState("");
  const [codeMode, setCodeMode] = useState(false);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");
  const [pendingImage, setPendingImage] = useState<File | null>(null);
  const [recording, setRecording] = useState(false);
  const [seconds, setSeconds] = useState(0);

  const fileRef = useRef<HTMLInputElement | null>(null);
  const recorderRef = useRef<MediaRecorder | null>(null);
  // Mirrors `seconds` so the recorder's onstop callback can read the live
  // count rather than the value captured when recording began.
  const secondsRef = useRef(0);
  const chunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // @mention autocomplete over the word currently being typed.
  const mentionQuery = useMemo(() => {
    const m = /@([A-Za-z][A-Za-z0-9_.-]*)$/.exec(body);
    return m ? m[1].toLowerCase() : null;
  }, [body]);

  const suggestions = useMemo(() => {
    if (mentionQuery === null) return [];
    return contacts
      .filter((c) => c.name.toLowerCase().includes(mentionQuery))
      .slice(0, 5);
  }, [mentionQuery, contacts]);

  const applyMention = (name: string) => {
    setBody((b) => b.replace(/@([A-Za-z][A-Za-z0-9_.-]*)$/, `@${name.replace(/\s+/g, "-")} `));
  };

  const send = async () => {
    const text = body.trim();
    if (!text || sending) return;
    setSending(true);
    setError("");
    try {
      await api.postIdeaMessage(ideaId, {
        body: text,
        parentId: replyTo?.id ?? null,
        kind: codeMode ? "code" : "text",
      });
      setBody("");
      setCodeMode(false);
      onSent();
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not send");
    } finally {
      setSending(false);
    }
  };

  const sendImage = async (file: File) => {
    setPendingImage(null);
    setSending(true);
    setError("");
    try {
      await api.postIdeaAttachment(ideaId, file, {
        kind: "image",
        body: body.trim(),
        parentId: replyTo?.id ?? null,
      });
      setBody("");
      onSent();
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not send image");
    } finally {
      setSending(false);
    }
  };

  // Voice memos use the browser's own recorder, so nothing extra ships.
  const startRecording = async () => {
    setError("");
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const rec = new MediaRecorder(stream);
      chunksRef.current = [];
      rec.ondataavailable = (e) => e.data.size && chunksRef.current.push(e.data);
      rec.onstop = async () => {
        stream.getTracks().forEach((t) => t.stop());
        const blob = new Blob(chunksRef.current, { type: "audio/webm" });
        if (blob.size > 0) {
          setSending(true);
          try {
            await api.postIdeaAttachment(
              ideaId,
              new File([blob], "memo.webm", { type: "audio/webm" }),
              {
                kind: "voice",
                parentId: replyTo?.id ?? null,
                // From the ref, not the `seconds` state: this callback is a
                // closure created when recording started, so it would capture
                // the count as it was then — always 0 — and every memo would
                // report a zero-second duration.
                duration: secondsRef.current,
              }
            );
            onSent();
          } catch (e) {
            setError(e instanceof Error ? e.message : "could not send recording");
          } finally {
            setSending(false);
          }
        }
        setSeconds(0);
        secondsRef.current = 0;
      };
      rec.start();
      recorderRef.current = rec;
      setRecording(true);
      setSeconds(0);
      secondsRef.current = 0;
      timerRef.current = setInterval(() => {
        secondsRef.current += 1;
        setSeconds(secondsRef.current);
      }, 1000);
    } catch {
      setError("Microphone permission was declined.");
    }
  };

  const stopRecording = () => {
    recorderRef.current?.stop();
    recorderRef.current = null;
    if (timerRef.current) clearInterval(timerRef.current);
    setRecording(false);
  };

  useEffect(() => {
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
      recorderRef.current?.stop();
    };
  }, []);

  return (
    <div className="shrink-0 border-t border-gray-100 bg-white p-3">
      {replyTo && (
        <div className="mb-2 flex items-center gap-2 rounded-lg bg-gray-50 px-3 py-1.5 text-label-md">
          <CornerDownRight size={12} className="text-purple" />
          <span className="min-w-0 flex-1 truncate text-slatey">
            Replying to{" "}
            <strong className="text-ink">
              {replyTo.anonymous ? "Anonymous" : replyTo.author?.name}
            </strong>
            : {replyTo.body.slice(0, 60) || "attachment"}
          </span>
          <button onClick={onCancelReply} aria-label="Cancel reply">
            <X size={13} className="text-slatey" />
          </button>
        </div>
      )}

      {anonymous && (
        <p className="mb-2 rounded-lg bg-ink px-3 py-1.5 text-label-md text-white">
          Silent brainstorm is running — this will post without your name.
        </p>
      )}

      {error && (
        <p className="mb-2 rounded-lg bg-coral-light px-3 py-1.5 text-label-md text-coral">
          {error}
        </p>
      )}

      {suggestions.length > 0 && (
        <div className="mb-2 flex flex-wrap gap-1.5">
          {suggestions.map((s) => (
            <button
              key={s.id}
              onClick={() => applyMention(s.name)}
              className="flex items-center gap-1.5 rounded-full bg-purple-light px-2 py-1 text-label-md font-bold text-purple"
            >
              <Avatar name={s.name} color={s.avatarColor} url={s.avatarUrl} size={16} />
              {s.name}
            </button>
          ))}
        </div>
      )}

      {recording ? (
        <div className="flex items-center gap-3 rounded-2xl bg-coral-light px-4 py-3">
          <span className="h-2.5 w-2.5 animate-pulse rounded-full bg-coral" />
          <span className="font-mono text-body-md font-bold text-coral">
            {Math.floor(seconds / 60)}:{String(seconds % 60).padStart(2, "0")}
          </span>
          <span className="text-body-sm text-slatey">Recording a voice memo…</span>
          <button
            onClick={stopRecording}
            className="ml-auto flex items-center gap-1 rounded-full bg-coral px-3 py-1.5 text-label-lg font-extrabold text-white"
          >
            <Square size={12} fill="currentColor" /> Stop &amp; send
          </button>
        </div>
      ) : (
        <div className="flex items-end gap-2">
          <div className="flex shrink-0 gap-1">
            <ComposerButton
              label="Add a photo"
              onClick={() => fileRef.current?.click()}
            >
              <ImageIcon size={17} />
            </ComposerButton>
            <ComposerButton label="Record a voice memo" onClick={startRecording}>
              <Mic size={17} />
            </ComposerButton>
            <ComposerButton
              label="Code block"
              active={codeMode}
              onClick={() => setCodeMode((c) => !c)}
            >
              <Code2 size={17} />
            </ComposerButton>
          </div>

          <textarea
            value={body}
            onChange={(e) => setBody(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey && !codeMode) {
                e.preventDefault();
                send();
              }
            }}
            rows={1}
            placeholder={
              codeMode
                ? "Paste code — Shift+Enter for a new line"
                : "Type a message… @name to mention, @idea#12 to link"
            }
            className={`max-h-32 min-h-[42px] flex-1 resize-y rounded-2xl px-4 py-2.5 text-body-md text-ink outline-none ring-purple/30 placeholder:text-gray-300 focus:ring-2 ${
              codeMode ? "bg-ink font-mono text-white" : "bg-gray-50"
            }`}
          />

          <button
            onClick={send}
            disabled={!body.trim() || sending}
            aria-label="Send"
            className="flex h-[42px] w-[42px] shrink-0 items-center justify-center rounded-full bg-purple text-white shadow-float transition hover:bg-purple-dark disabled:opacity-40"
          >
            {sending ? (
              <Loader2 size={17} className="animate-spin" />
            ) : (
              <Send size={17} />
            )}
          </button>
        </div>
      )}

      <input
        ref={fileRef}
        type="file"
        accept="image/*"
        hidden
        onChange={(e) => {
          const f = e.target.files?.[0];
          if (f) setPendingImage(f);
          e.target.value = "";
        }}
      />

      {pendingImage && (
        <ImageEditorModal
          file={pendingImage}
          title="Share a photo"
          onCancel={() => setPendingImage(null)}
          onDone={sendImage}
        />
      )}
    </div>
  );
}

function ComposerButton({
  label,
  active,
  onClick,
  children,
}: {
  label: string;
  active?: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      title={label}
      aria-label={label}
      className={`flex h-[42px] w-[38px] items-center justify-center rounded-lg transition ${
        active
          ? "bg-purple-light text-purple"
          : "text-slatey hover:bg-gray-50 hover:text-ink"
      }`}
    >
      {children}
    </button>
  );
}
