"use client";

import Link from "next/link";
import { Bell } from "lucide-react";
import { useAuth } from "@/lib/auth";
import { useUnreadCount } from "@/lib/notifications";

/**
 * Bell button that links to the notifications screen and shows an unread badge.
 * Designed to sit in a coloured header (white icon) by default.
 */
export function NotificationBell({ className = "" }: { className?: string }) {
  const { user } = useAuth();
  const count = useUnreadCount(user);

  return (
    <Link
      href="/notifications"
      aria-label={
        count > 0 ? `Notifications, ${count} unread` : "Notifications"
      }
      className={`relative flex h-10 w-10 items-center justify-center rounded-full bg-white/15 text-white transition hover:bg-white/25 ${className}`}
    >
      <Bell size={20} strokeWidth={2.2} />
      {count > 0 && (
        <span className="absolute -right-0.5 -top-0.5 flex h-5 min-w-[20px] items-center justify-center rounded-full border-2 border-purple bg-coral px-1 text-label-sm font-extrabold text-white">
          {count > 9 ? "9+" : count}
        </span>
      )}
    </Link>
  );
}
