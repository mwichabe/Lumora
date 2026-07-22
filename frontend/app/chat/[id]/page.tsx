"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  ArrowLeft,
  Check,
  Image as ImageIcon,
  Loader2,
  Pencil,
  Send,
  Trash2,
} from "lucide-react";
import {
  ActionMenu,
  ActionMenuItem,
  ActionMenuNote,
} from "@/components/ActionMenu";
import { Avatar } from "@/components/Avatar";
import { TranslatedText } from "@/components/TranslatedText";
import { FoxMascot } from "@/components/FoxMascot";
import { ImageEditorModal } from "@/components/ImageEditorModal";
import { useAuth } from "@/lib/auth";
import { api, mediaUrl } from "@/lib/api";
import type { ChatMessage, ChatUser } from "@/lib/types";

export default function ChatThreadPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { user } = useAuth();
  const myId = user?.id;

  const [other, setOther] = useState<ChatUser | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [text, setText] = useState("");
  const [sending, setSending] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState("");
  const [pendingImage, setPendingImage] = useState<File | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);

  const bottomRef = useRef<HTMLDivElement | null>(null);
  const fileRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    let active = true;
    const load = () =>
      api
        .chatMessages(id)
        .then((d) => {
          if (!active) return;
          setOther(d.user);
          setMessages(d.messages);
          setLoaded(true);
        })
        .catch(() => setLoaded(true));
    load();
    const t = setInterval(load, 4000);
    return () => {
      active = false;
      clearInterval(t);
    };
  }, [id]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages.length]);

  async function send() {
    const body = text.trim();
    if (!body || sending) return;
    setSending(true);
    setText("");
    setError("");
    try {
      const r = await api.sendChatMessage(id, body);
      setMessages((m) => [...m, r.message]);
    } catch (e) {
      setText(body); // restore on failure
      setError(e instanceof Error ? e.message : "could not send");
    } finally {
      setSending(false);
    }
  }

  async function sendImage(file: File) {
    setPendingImage(null);
    setSending(true);
    setError("");
    const caption = text.trim();
    try {
      const r = await api.sendChatImage(id, file, caption);
      setMessages((m) => [...m, r.message]);
      setText("");
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not send photo");
    } finally {
      setSending(false);
    }
  }

  async function saveEdit(messageId: number, body: string) {
    try {
      const r = await api.editChatMessage(messageId, body);
      setMessages((m) => m.map((x) => (x.id === messageId ? r.message : x)));
      setEditingId(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not edit");
    }
  }

  // Retry path: the server translates in the background on send, so this only
  // fires when that didn't land.
  async function translate(messageId: number) {
    try {
      const r = await api.translateChatMessage(messageId);
      if (r.translation) {
        setMessages((m) =>
          m.map((x) =>
            x.id === messageId ? { ...x, translation: r.translation! } : x
          )
        );
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not translate");
    }
  }

  async function remove(messageId: number) {
    try {
      const r = await api.deleteChatMessage(messageId);
      setMessages((m) => m.map((x) => (x.id === messageId ? r.message : x)));
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not delete");
    }
  }

  return (
    <div className="flex min-h-[100dvh] w-full justify-center bg-cream lg:bg-[#eceaf3]">
      <div className="flex h-[100dvh] w-full max-w-2xl flex-col bg-cream lg:my-8 lg:h-[calc(100dvh-4rem)] lg:overflow-hidden lg:rounded-3xl lg:shadow-card-lg">
        {/* Header */}
        <header className="flex items-center gap-3 border-b border-gray-100 bg-white px-3 py-3 pt-10 lg:rounded-t-3xl lg:pt-3">
          <button
            onClick={() => router.push("/chat")}
            aria-label="Back"
            className="flex h-9 w-9 items-center justify-center rounded-full text-slatey transition hover:bg-gray-50"
          >
            <ArrowLeft size={20} />
          </button>
          {other && (
            <>
              <Avatar
                name={other.name}
                color={other.avatarColor}
                url={other.avatarUrl}
                size={40}
              />
              <div className="min-w-0">
                <p className="truncate font-extrabold text-ink">{other.name}</p>
                <p className="truncate text-label-md text-slatey">
                  {other.levelName || "Learner"}
                </p>
              </div>
            </>
          )}
        </header>

        {/* Messages */}
        <div className="flex-1 space-y-2 overflow-y-auto px-4 py-4">
          {!loaded ? (
            <div className="flex h-full items-center justify-center">
              <FoxMascot size={90} glow />
            </div>
          ) : messages.length === 0 ? (
            <div className="flex h-full flex-col items-center justify-center text-center">
              <p className="text-body-md font-bold text-ink">
                Say hi to {other?.name || "your new friend"} 👋
              </p>
              <p className="text-body-sm text-slatey">
                Practising together makes it stick.
              </p>
            </div>
          ) : (
            messages.map((m) => (
              <Bubble
                key={m.id}
                message={m}
                mine={m.senderId === myId}
                editing={editingId === m.id}
                onStartEdit={() => setEditingId(m.id)}
                onCancelEdit={() => setEditingId(null)}
                onSaveEdit={(body) => saveEdit(m.id, body)}
                onDelete={() => remove(m.id)}
                onTranslate={() => translate(m.id)}
              />
            ))
          )}
          <div ref={bottomRef} />
        </div>

        {error && (
          <p className="mx-4 mb-2 rounded-lg bg-coral-light px-3 py-2 text-label-md text-coral">
            {error}
          </p>
        )}

        {/* Composer */}
        <div className="flex items-center gap-2 border-t border-gray-100 bg-white px-3 py-3 pb-[max(0.75rem,env(safe-area-inset-bottom))]">
          <button
            onClick={() => fileRef.current?.click()}
            aria-label="Send a photo"
            title="Send a photo"
            className="flex h-12 w-11 shrink-0 items-center justify-center rounded-full text-slatey transition hover:bg-gray-50 hover:text-ink"
          >
            <ImageIcon size={20} />
          </button>
          <input
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") send();
            }}
            placeholder="Type a message…"
            className="h-12 min-w-0 flex-1 rounded-full border border-gray-100 bg-gray-50 px-4 text-body-lg outline-none transition focus:border-purple focus:bg-white"
          />
          <button
            onClick={send}
            disabled={!text.trim() || sending}
            aria-label="Send"
            className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-purple text-white transition hover:bg-purple-dark disabled:opacity-40"
          >
            {sending ? (
              <Loader2 size={20} className="animate-spin" />
            ) : (
              <Send size={20} />
            )}
          </button>
        </div>

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
            title="Send a photo"
            onCancel={() => setPendingImage(null)}
            onDone={sendImage}
          />
        )}
      </div>
    </div>
  );
}

function Bubble({
  message,
  mine,
  editing,
  onStartEdit,
  onCancelEdit,
  onSaveEdit,
  onDelete,
  onTranslate,
}: {
  message: ChatMessage;
  mine: boolean;
  editing: boolean;
  onStartEdit: () => void;
  onCancelEdit: () => void;
  onSaveEdit: (body: string) => void;
  onDelete: () => void;
  onTranslate: () => Promise<void>;
}) {
  const [draft, setDraft] = useState(message.body);

  useEffect(() => setDraft(message.body), [message.body]);

  // A deleted message keeps its place — the other person already read the
  // reply that answers it, and a silent gap would make that reply nonsense.
  if (message.deleted) {
    return (
      <div className={`flex ${mine ? "justify-end" : "justify-start"}`}>
        <div className="flex items-center gap-1.5 rounded-2xl bg-gray-100 px-4 py-2 text-body-sm italic text-gray-500">
          <Trash2 size={12} /> Message deleted
        </div>
      </div>
    );
  }

  if (editing) {
    return (
      <div className="flex justify-end">
        <div className="w-[80%] rounded-2xl bg-white p-2 shadow-card">
          <textarea
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            rows={2}
            autoFocus
            className="w-full resize-none rounded-lg bg-gray-50 p-2 text-body-md text-ink outline-none ring-purple/30 focus:ring-2"
          />
          <div className="mt-2 flex justify-end gap-2">
            <button
              onClick={onCancelEdit}
              className="rounded-full px-3 py-1 text-label-lg font-bold text-slatey hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              onClick={() => onSaveEdit(draft.trim())}
              disabled={!draft.trim()}
              className="flex items-center gap-1 rounded-full bg-purple px-3 py-1 text-label-lg font-extrabold text-white disabled:opacity-40"
            >
              <Check size={13} /> Save
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={`group flex items-end gap-1 ${mine ? "justify-end" : "justify-start"}`}>
      {mine && (
        <ActionMenu>
          {(close) => (
            <>
              {message.canEdit ? (
                <ActionMenuItem
                  icon={<Pencil size={13} />}
                  onClick={() => {
                    onStartEdit();
                    close();
                  }}
                >
                  Edit
                </ActionMenuItem>
              ) : (
                <ActionMenuNote>
                  {message.kind === "image"
                    ? "Photos can't be edited."
                    : "Edits close after 24 hours."}
                </ActionMenuNote>
              )}
              <ActionMenuItem
                icon={<Trash2 size={13} />}
                tone="danger"
                onClick={() => {
                  onDelete();
                  close();
                }}
              >
                Delete
              </ActionMenuItem>
            </>
          )}
        </ActionMenu>
      )}

      <div
        className={`max-w-[78%] overflow-hidden rounded-2xl text-body-md ${
          mine
            ? "rounded-br-md bg-purple text-white"
            : "rounded-bl-md bg-white text-ink shadow-card"
        } ${message.kind === "image" ? "p-1" : "px-4 py-2.5"}`}
      >
        {message.kind === "image" ? (
          <>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={mediaUrl(message.url)}
              alt={message.fileName || "Shared photo"}
              width={message.width || undefined}
              height={message.height || undefined}
              className="max-h-72 w-auto max-w-full rounded-xl object-contain"
            />
            {message.body && (
              <p className="px-3 py-1.5">
                <TranslatedText
                  body={message.body}
                  translation={message.translation}
                  tone={mine ? "dark" : "light"}
                  onRetry={onTranslate}
                >
                  {(text) => text}
                </TranslatedText>
              </p>
            )}
          </>
        ) : (
          <TranslatedText
            body={message.body}
            translation={message.translation}
            tone={mine ? "dark" : "light"}
            onRetry={onTranslate}
          >
            {(text) => text}
          </TranslatedText>
        )}

        {message.edited && (
          <span
            className={`ml-2 text-label-sm italic ${
              mine ? "text-white/60" : "text-gray-500"
            }`}
          >
            edited
          </span>
        )}
      </div>
    </div>
  );
}
