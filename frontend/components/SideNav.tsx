"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, BookOpen, Mic, Trophy, User, Bell, LogOut } from "lucide-react";
import { FoxMascot } from "./FoxMascot";
import { useAuth } from "@/lib/auth";
import { useUnreadCount } from "@/lib/notifications";

const navItems = [
  { href: "/home", label: "Home", icon: Home },
  { href: "/learn", label: "Learn", icon: BookOpen },
  { href: "/practice", label: "Practice", icon: Mic },
  { href: "/leaderboard", label: "Leagues", icon: Trophy },
  { href: "/notifications", label: "Notifications", icon: Bell },
  { href: "/profile", label: "Profile", icon: User },
];

/**
 * Desktop-only left navigation rail. Hidden below `lg` where the BottomTabBar
 * takes over, so the experience stays phone-native on small devices and turns
 * into a proper app layout on wide screens.
 */
export function SideNav() {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const unread = useUnreadCount(user);

  return (
    <aside className="sticky top-0 hidden h-[100dvh] w-64 shrink-0 flex-col border-r border-gray-100 bg-white px-4 py-6 lg:flex">
      {/* Brand */}
      <Link href="/home" className="mb-8 flex items-center gap-2 px-2">
        <FoxMascot size={40} />
        <span className="text-heading-lg font-extrabold tracking-tight text-purple">
          Lumora
        </span>
      </Link>

      {/* Nav */}
      <nav className="flex flex-col gap-1">
        {navItems.map(({ href, label, icon: Icon }) => {
          const active = pathname?.startsWith(href);
          return (
            <Link
              key={href}
              href={href}
              className={`flex items-center gap-3 rounded-xl px-4 py-3 text-body-lg font-extrabold transition ${
                active
                  ? "bg-purple-light text-purple"
                  : "text-slatey hover:bg-gray-50 hover:text-ink"
              }`}
            >
              <Icon
                size={24}
                strokeWidth={2.2}
                fill={active ? "#EDE7F6" : "transparent"}
              />
              {label}
              {href === "/notifications" && unread > 0 && (
                <span className="ml-auto flex h-5 min-w-[20px] items-center justify-center rounded-full bg-coral px-1 text-label-sm font-extrabold text-white">
                  {unread > 9 ? "9+" : unread}
                </span>
              )}
            </Link>
          );
        })}
      </nav>

      {/* User card pinned to bottom */}
      {user && (
        <div className="mt-auto">
          <div className="flex items-center gap-3 rounded-2xl bg-gray-50 p-3">
            <span
              className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-sm font-bold text-white"
              style={{ backgroundColor: user.avatarColor || "#6C3FC5" }}
            >
              {(user.name || "L").slice(0, 1).toUpperCase()}
            </span>
            <div className="min-w-0 flex-1">
              <p className="truncate text-body-md font-extrabold text-ink">
                {user.name || "Learner"}
              </p>
              <p className="truncate text-label-md text-slatey">
                {user.levelName || "Spark"}
              </p>
            </div>
            <button
              onClick={logout}
              aria-label="Sign out"
              className="text-gray-300 transition hover:text-coral"
            >
              <LogOut size={18} />
            </button>
          </div>
        </div>
      )}
    </aside>
  );
}
