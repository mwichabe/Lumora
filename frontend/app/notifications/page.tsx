"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, BellOff, X, ChevronRight, Trash2 } from "lucide-react";
import { motion } from "framer-motion";
import { AppShell } from "@/components/AppShell";
import { api } from "@/lib/api";
import type { AppNotification } from "@/lib/types";

const SYNC_EVENT = "lumora:notifications";

const KIND_LABEL: Record<string, string> = {
  welcome: "Welcome",
  welcome_back: "Welcome back",
  milestone: "Milestone",
  exam: "Exam",
  hearts: "Hearts",
  payment: "Payment",
  tip: "Tip",
  feature: "New feature",
  language: "Languages",
  streak: "Streak",
  league: "League",
  chat: "Message",
  message: "Message",
};

function kindLabel(kind: string): string {
  return KIND_LABEL[kind] || "Notification";
}

function timeAgo(iso: string): string {
  const t = new Date(iso).getTime();
  if (!t) return "";
  const s = Math.max(0, Math.floor((Date.now() - t) / 1000));
  if (s < 60) return "just now";
  const m = Math.floor(s / 60);
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  const d = Math.floor(h / 24);
  return `${d}d ago`;
}

function fullDate(iso: string): string {
  const d = new Date(iso);
  if (isNaN(d.getTime())) return "";
  return d.toLocaleString(undefined, { dateStyle: "medium", timeStyle: "short" });
}

export default function NotificationsPage() {
  return (
    <AppShell tabs>
      <NotificationsContent />
    </AppShell>
  );
}

function NotificationsContent() {
  const router = useRouter();
  const [items, setItems] = useState<AppNotification[]>([]);
  const [selected, setSelected] = useState<AppNotification | null>(null);

  useEffect(() => {
    // Just load the list — viewing the tab no longer marks anything read.
    let active = true;
    api
      .notifications()
      .then((d) => {
        if (active) setItems(d.notifications);
      })
      .catch(() => {});
    return () => {
      active = false;
    };
  }, []);

  const unread = items.filter((n) => !n.read).length;
  const newItems = useMemo(() => items.filter((n) => !n.read), [items]);
  const earlierItems = useMemo(() => items.filter((n) => n.read), [items]);

  // Open one notification: show its full content + metadata, and mark ONLY that
  // one as read.
  function openNotification(n: AppNotification) {
    setSelected(n);
    if (!n.read) {
      setItems((prev) =>
        prev.map((x) => (x.id === n.id ? { ...x, read: true } : x))
      );
      api
        .markNotificationRead(n.id)
        .then(() => window.dispatchEvent(new Event(SYNC_EVENT)))
        .catch(() => {});
    }
  }

  function clearHighlights() {
    api
      .markNotificationsRead()
      .then(() => window.dispatchEvent(new Event(SYNC_EVENT)))
      .catch(() => {});
    setItems((prev) => prev.map((n) => ({ ...n, read: true })));
  }

  function deleteOne(id: number) {
    setItems((prev) => prev.filter((x) => x.id !== id));
    setSelected((s) => (s && s.id === id ? null : s));
    api
      .deleteNotification(id)
      .then(() => window.dispatchEvent(new Event(SYNC_EVENT)))
      .catch(() => {});
  }

  return (
    <div className="min-h-full bg-white">
      {/* Header */}
      <header className="sticky top-0 z-10 flex items-center gap-3 border-b border-gray-100 bg-white px-4 py-3.5 lg:rounded-t-3xl lg:px-6 lg:py-4">
        <button
          onClick={() => router.back()}
          aria-label="Back"
          className="flex h-9 w-9 items-center justify-center rounded-full text-slatey transition hover:bg-gray-50 lg:hidden"
        >
          <ArrowLeft size={20} />
        </button>
        <div className="flex-1">
          <h1 className="text-heading-lg font-extrabold text-ink">
            Notifications
          </h1>
        </div>
        {unread > 0 && (
          <button
            onClick={clearHighlights}
            className="rounded-lg px-3 py-1.5 text-body-sm font-bold text-purple transition hover:bg-purple-light"
          >
            Mark all as read
          </button>
        )}
      </header>

      {items.length === 0 ? (
        <EmptyState />
      ) : (
        <>
          {newItems.length > 0 && (
            <Section label="New">
              {newItems.map((n) => (
                <NotificationItem
                  key={n.id}
                  n={n}
                  onOpen={openNotification}
                  onDelete={deleteOne}
                />
              ))}
            </Section>
          )}
          {earlierItems.length > 0 && (
            <Section label="Earlier">
              {earlierItems.map((n) => (
                <NotificationItem
                  key={n.id}
                  n={n}
                  onOpen={openNotification}
                  onDelete={deleteOne}
                />
              ))}
            </Section>
          )}
        </>
      )}

      <NotificationDetail
        n={selected}
        onClose={() => setSelected(null)}
        onDelete={deleteOne}
        onFollowLink={(link) => {
          setSelected(null);
          router.push(link);
        }}
      />
    </div>
  );
}

function Section({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <section>
      <h2 className="px-4 pb-1 pt-5 text-label-md font-bold uppercase tracking-wide text-gray-500 lg:px-6">
        {label}
      </h2>
      <ul className="divide-y divide-gray-100">{children}</ul>
    </section>
  );
}

function NotificationItem({
  n,
  onOpen,
  onDelete,
}: {
  n: AppNotification;
  onOpen: (n: AppNotification) => void;
  onDelete: (id: number) => void;
}) {
  return (
    <li
      className={`group relative flex items-center ${
        n.read ? "" : "bg-purple-light/40"
      }`}
    >
      {!n.read && <span className="absolute left-0 top-0 h-full w-1 bg-purple" />}

      <button
        onClick={() => onOpen(n)}
        className="flex flex-1 items-start gap-3 px-4 py-4 text-left transition hover:bg-black/[0.02] lg:px-6"
      >
        <span
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-xl"
          style={{ backgroundColor: n.tint + "1F" }}
        >
          {n.emoji}
        </span>

        <div className="min-w-0 flex-1">
          <div className="flex items-baseline justify-between gap-2">
            <p
              className={`truncate text-body-md text-ink ${
                n.read ? "font-semibold" : "font-extrabold"
              }`}
            >
              {n.title}
            </p>
            <span className="shrink-0 text-label-md text-gray-500">
              {timeAgo(n.createdAt)}
            </span>
          </div>
          <p className="mt-0.5 line-clamp-1 text-body-sm leading-snug text-slatey">
            {n.body}
          </p>
        </div>

        <ChevronRight size={18} className="mt-1 shrink-0 text-gray-300" />
      </button>

      <button
        onClick={() => onDelete(n.id)}
        aria-label="Delete notification"
        className="mr-2 flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-gray-300 transition hover:bg-coral/10 hover:text-coral"
      >
        <Trash2 size={16} />
      </button>
    </li>
  );
}

/** Full-detail view of a single notification with its metadata. */
function NotificationDetail({
  n,
  onClose,
  onDelete,
  onFollowLink,
}: {
  n: AppNotification | null;
  onClose: () => void;
  onDelete: (id: number) => void;
  onFollowLink: (link: string) => void;
}) {
  if (!n) return null;
  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center bg-black/50 p-0 sm:items-center sm:p-4"
      onClick={onClose}
    >
      <motion.div
        initial={{ y: 40, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        onClick={(e) => e.stopPropagation()}
        className="w-full max-w-md overflow-hidden rounded-t-3xl bg-white shadow-card-lg sm:rounded-3xl"
      >
        <div
          className="relative flex flex-col items-center px-6 pb-5 pt-7"
          style={{ backgroundColor: n.tint + "1A" }}
        >
          <button
            onClick={onClose}
            aria-label="Close"
            className="absolute right-4 top-4 flex h-8 w-8 items-center justify-center rounded-full bg-white/70 text-gray-500 transition hover:bg-white"
          >
            <X size={18} />
          </button>
          <span
            className="flex h-16 w-16 items-center justify-center rounded-full text-3xl"
            style={{ backgroundColor: n.tint + "2E" }}
          >
            {n.emoji}
          </span>
          <span
            className="mt-3 rounded-full px-2.5 py-1 text-label-sm font-bold uppercase tracking-wide"
            style={{ backgroundColor: n.tint + "22", color: n.tint }}
          >
            {kindLabel(n.kind)}
          </span>
          <h2 className="mt-2 text-center text-heading-lg font-extrabold text-ink">
            {n.title}
          </h2>
        </div>

        <div className="px-6 py-5">
          <p className="whitespace-pre-line text-body-md leading-relaxed text-ink">
            {n.body}
          </p>

          {/* metadata */}
          <div className="mt-5 space-y-2 rounded-2xl bg-gray-50 p-4">
            <MetaRow label="Type" value={kindLabel(n.kind)} />
            <MetaRow label="Received" value={fullDate(n.createdAt)} />
            <MetaRow label="Status" value="Read" />
          </div>

          <div className="mt-4 flex gap-2">
            <button
              onClick={() => onDelete(n.id)}
              className="flex items-center justify-center gap-1.5 rounded-full border-2 border-gray-100 px-4 py-3 font-bold text-slatey transition hover:border-coral hover:text-coral"
            >
              <Trash2 size={16} /> Delete
            </button>
            {n.link ? (
              <button
                onClick={() => onFollowLink(n.link)}
                className="flex flex-1 items-center justify-center gap-2 rounded-full bg-purple py-3 font-extrabold text-white shadow-float"
              >
                Open <ChevronRight size={18} />
              </button>
            ) : (
              <button
                onClick={onClose}
                className="flex-1 rounded-full bg-purple py-3 font-extrabold text-white shadow-float"
              >
                Done
              </button>
            )}
          </div>
        </div>
      </motion.div>
    </div>
  );
}

function MetaRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="text-label-md text-slatey">{label}</span>
      <span className="text-body-sm font-bold text-ink">{value}</span>
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center gap-2 px-6 py-24 text-center">
      <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gray-50">
        <BellOff className="text-gray-300" size={26} />
      </div>
      <p className="mt-1 text-body-md font-bold text-ink">No notifications yet</p>
      <p className="max-w-xs text-body-sm text-slatey">
        When something happens — a streak reminder, a new quest, a league update
        — you&apos;ll see it here.
      </p>
    </div>
  );
}
