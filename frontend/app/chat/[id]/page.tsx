"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft, Send } from "lucide-react";
import { Avatar } from "@/components/Avatar";
import { FoxMascot } from "@/components/FoxMascot";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
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
  const bottomRef = useRef<HTMLDivElement | null>(null);

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
    try {
      const r = await api.sendChatMessage(id, body);
      setMessages((m) => [...m, r.message]);
    } catch {
      setText(body); // restore on failure
    } finally {
      setSending(false);
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
            messages.map((m) => {
              const mine = m.senderId === myId;
              return (
                <div
                  key={m.id}
                  className={`flex ${mine ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`max-w-[78%] rounded-2xl px-4 py-2.5 text-body-md ${
                      mine
                        ? "rounded-br-md bg-purple text-white"
                        : "rounded-bl-md bg-white text-ink shadow-card"
                    }`}
                  >
                    {m.body}
                  </div>
                </div>
              );
            })
          )}
          <div ref={bottomRef} />
        </div>

        {/* Composer */}
        <div className="flex items-center gap-2 border-t border-gray-100 bg-white px-3 py-3 pb-[max(0.75rem,env(safe-area-inset-bottom))]">
          <input
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") send();
            }}
            placeholder="Type a message…"
            className="h-12 flex-1 rounded-full border border-gray-100 bg-gray-50 px-4 text-body-lg outline-none transition focus:border-purple focus:bg-white"
          />
          <button
            onClick={send}
            disabled={!text.trim() || sending}
            aria-label="Send"
            className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-purple text-white transition hover:bg-purple-dark disabled:opacity-40"
          >
            <Send size={20} />
          </button>
        </div>
      </div>
    </div>
  );
}
