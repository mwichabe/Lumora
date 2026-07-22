"use client";

import Link from "next/link";
import { Lightbulb } from "lucide-react";

/**
 * Ideas shortcut for coloured headers — the sibling of ChatBell and
 * NotificationBell.
 *
 * On a phone Ideas lives behind the More sheet, which is fine for a
 * destination you already know about and poor for one you don't. Sitting it
 * beside the bells on the home screen means the feature is discoverable from
 * the first screen anyone sees, rather than depending on someone opening a
 * menu to find out it exists.
 */
export function IdeasBell({ className = "" }: { className?: string }) {
  return (
    <Link
      href="/ideas"
      aria-label="Ideas"
      title="Ideas"
      className={`relative flex h-10 w-10 items-center justify-center rounded-full bg-white/15 text-white transition hover:bg-white/25 ${className}`}
    >
      <Lightbulb size={20} strokeWidth={2.2} />
    </Link>
  );
}
