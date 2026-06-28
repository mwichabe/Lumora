"use client";

import Link from "next/link";
import { MessageCircle } from "lucide-react";
import { useChatUnread } from "@/lib/chat";

/** Chat icon with an unread-message badge, for coloured headers. */
export function ChatBell({ className = "" }: { className?: string }) {
  const count = useChatUnread();
  return (
    <Link
      href="/chat"
      aria-label={count > 0 ? `Messages, ${count} unread` : "Messages"}
      className={`relative flex h-10 w-10 items-center justify-center rounded-full bg-white/15 text-white transition hover:bg-white/25 ${className}`}
    >
      <MessageCircle size={20} strokeWidth={2.2} />
      {count > 0 && (
        <span className="absolute -right-0.5 -top-0.5 flex h-5 min-w-[20px] items-center justify-center rounded-full border-2 border-purple bg-coral px-1 text-label-sm font-extrabold text-white">
          {count > 9 ? "9+" : count}
        </span>
      )}
    </Link>
  );
}
