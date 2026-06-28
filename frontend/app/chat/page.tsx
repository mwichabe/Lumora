"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { MessageSquarePlus, X, MessagesSquare } from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { Avatar } from "@/components/Avatar";
import { api } from "@/lib/api";
import type { ChatThread, ChatUser } from "@/lib/types";

function timeAgo(iso: string): string {
  const t = new Date(iso).getTime();
  if (!t) return "";
  const s = Math.max(0, Math.floor((Date.now() - t) / 1000));
  if (s < 60) return "now";
  const m = Math.floor(s / 60);
  if (m < 60) return `${m}m`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h`;
  return `${Math.floor(h / 24)}d`;
}

export default function ChatPage() {
  return (
    <AppShell tabs>
      <ChatList />
    </AppShell>
  );
}

function ChatList() {
  const [threads, setThreads] = useState<ChatThread[] | null>(null);
  const [picking, setPicking] = useState(false);

  useEffect(() => {
    const load = () =>
      api.chatThreads().then((d) => setThreads(d.threads)).catch(() => setThreads([]));
    load();
    const id = setInterval(load, 8000);
    return () => clearInterval(id);
  }, []);

  return (
    <div className="min-h-full bg-cream pb-24 lg:pb-10">
      <header className="flex items-center justify-between border-b border-gray-100 bg-white px-5 py-4 lg:rounded-t-3xl lg:px-6">
        <h1 className="text-heading-lg font-extrabold text-ink">Messages</h1>
        <button
          onClick={() => setPicking(true)}
          className="flex items-center gap-1.5 rounded-full bg-purple px-3 py-2 text-label-lg font-bold text-white transition hover:bg-purple-dark"
        >
          <MessageSquarePlus size={16} /> New
        </button>
      </header>

      {threads === null ? (
        <div className="space-y-2 p-4">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-16 animate-pulse rounded-2xl bg-gray-100" />
          ))}
        </div>
      ) : threads.length === 0 ? (
        <div className="flex flex-col items-center gap-2 px-6 py-20 text-center">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gray-50">
            <MessagesSquare className="text-gray-300" size={26} />
          </div>
          <p className="mt-1 text-body-md font-bold text-ink">No conversations yet</p>
          <p className="max-w-xs text-body-sm text-slatey">
            Start chatting with fellow learners — tap “New”.
          </p>
        </div>
      ) : (
        <ul className="divide-y divide-gray-100">
          {threads.map((t) => (
            <li key={t.user.id}>
              <Link
                href={`/chat/${t.user.id}`}
                className="flex items-center gap-3 px-5 py-3.5 transition hover:bg-gray-50 lg:px-6"
              >
                <Avatar
                  name={t.user.name}
                  color={t.user.avatarColor}
                  url={t.user.avatarUrl}
                  size={48}
                />
                <div className="min-w-0 flex-1">
                  <div className="flex items-baseline justify-between gap-2">
                    <p className="truncate font-extrabold text-ink">{t.user.name}</p>
                    <span className="shrink-0 text-label-md text-gray-500">
                      {timeAgo(t.lastAt)}
                    </span>
                  </div>
                  <p
                    className={`truncate text-body-sm ${
                      t.unread > 0 ? "font-bold text-ink" : "text-slatey"
                    }`}
                  >
                    {t.lastMessage}
                  </p>
                </div>
                {t.unread > 0 && (
                  <span className="flex h-5 min-w-[20px] items-center justify-center rounded-full bg-coral px-1 text-label-sm font-extrabold text-white">
                    {t.unread > 9 ? "9+" : t.unread}
                  </span>
                )}
              </Link>
            </li>
          ))}
        </ul>
      )}

      {picking && <ContactPicker onClose={() => setPicking(false)} />}
    </div>
  );
}

function ContactPicker({ onClose }: { onClose: () => void }) {
  const router = useRouter();
  const [contacts, setContacts] = useState<ChatUser[] | null>(null);
  const [q, setQ] = useState("");

  useEffect(() => {
    api.chatContacts().then((d) => setContacts(d.contacts)).catch(() => setContacts([]));
  }, []);

  const filtered = (contacts || []).filter((c) =>
    c.name.toLowerCase().includes(q.toLowerCase())
  );

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/50 sm:items-center" onClick={onClose}>
      <div
        className="max-h-[80vh] w-full max-w-md overflow-hidden rounded-t-2xl bg-white sm:rounded-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-gray-100 px-5 py-4">
          <h3 className="text-heading-sm font-extrabold text-ink">New message</h3>
          <button onClick={onClose} aria-label="Close" className="text-gray-500">
            <X size={20} />
          </button>
        </div>
        <div className="p-3">
          <input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="Search learners…"
            className="h-11 w-full rounded-xl border border-gray-100 bg-gray-50 px-4 outline-none focus:border-purple focus:bg-white"
          />
        </div>
        <ul className="max-h-[55vh] divide-y divide-gray-100 overflow-y-auto">
          {filtered.map((c) => (
            <li key={c.id}>
              <button
                onClick={() => router.push(`/chat/${c.id}`)}
                className="flex w-full items-center gap-3 px-5 py-3 text-left transition hover:bg-gray-50"
              >
                <Avatar name={c.name} color={c.avatarColor} url={c.avatarUrl} size={40} />
                <div className="min-w-0">
                  <p className="truncate font-extrabold text-ink">{c.name}</p>
                  <p className="truncate text-label-md text-slatey">{c.levelName}</p>
                </div>
              </button>
            </li>
          ))}
          {contacts && filtered.length === 0 && (
            <li className="px-5 py-8 text-center text-body-sm text-slatey">
              No learners found.
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}
