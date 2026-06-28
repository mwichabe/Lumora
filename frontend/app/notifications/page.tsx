"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowLeft, BellOff } from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { api } from "@/lib/api";
import type { AppNotification } from "@/lib/types";

const SYNC_EVENT = "lumora:notifications";

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

  useEffect(() => {
    let active = true;
    api
      .notifications()
      .then((d) => {
        if (active) setItems(d.notifications);
      })
      .catch(() => {});

    // Clear the unread badge shortly after the screen is viewed.
    const t = setTimeout(() => {
      api
        .markNotificationsRead()
        .then(() => window.dispatchEvent(new Event(SYNC_EVENT)))
        .catch(() => {});
    }, 800);

    return () => {
      active = false;
      clearTimeout(t);
    };
  }, []);

  const unread = items.filter((n) => !n.read).length;
  const newItems = useMemo(() => items.filter((n) => !n.read), [items]);
  const earlierItems = useMemo(() => items.filter((n) => n.read), [items]);

  function clearHighlights() {
    api
      .markNotificationsRead()
      .then(() => window.dispatchEvent(new Event(SYNC_EVENT)))
      .catch(() => {});
    setItems((prev) => prev.map((n) => ({ ...n, read: true })));
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
                <NotificationItem key={n.id} n={n} />
              ))}
            </Section>
          )}
          {earlierItems.length > 0 && (
            <Section label="Earlier">
              {earlierItems.map((n) => (
                <NotificationItem key={n.id} n={n} />
              ))}
            </Section>
          )}
        </>
      )}
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

function NotificationItem({ n }: { n: AppNotification }) {
  const cls = `relative flex items-start gap-3 px-4 py-4 transition hover:bg-gray-50 lg:px-6 ${
    n.read ? "" : "bg-purple-light/40"
  }`;

  const content = (
    <>
      {!n.read && <span className="absolute left-0 top-0 h-full w-1 bg-purple" />}

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
        <p className="mt-0.5 text-body-sm leading-snug text-slatey">{n.body}</p>
      </div>

      {!n.read && (
        <span className="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-purple" />
      )}
    </>
  );

  return (
    <li>
      {n.link ? (
        <Link href={n.link} className={cls}>
          {content}
        </Link>
      ) : (
        <div className={cls}>{content}</div>
      )}
    </li>
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
