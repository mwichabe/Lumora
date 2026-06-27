"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, BookOpen, Mic, Trophy, User } from "lucide-react";

const tabs = [
  { href: "/home", label: "Home", icon: Home },
  { href: "/learn", label: "Learn", icon: BookOpen },
  { href: "/practice", label: "Practice", icon: Mic },
  { href: "/leaderboard", label: "Leagues", icon: Trophy },
  { href: "/profile", label: "Profile", icon: User },
];

export function BottomTabBar() {
  const pathname = usePathname();
  return (
    <nav className="sticky bottom-0 z-20 mt-auto flex h-16 items-stretch border-t border-gray-100 bg-white lg:hidden">
      {tabs.map(({ href, label, icon: Icon }) => {
        const active = pathname?.startsWith(href);
        return (
          <Link
            key={href}
            href={href}
            className="relative flex flex-1 flex-col items-center justify-center gap-0.5"
          >
            {active && (
              <span className="absolute top-1.5 h-1.5 w-1.5 rounded-full bg-purple" />
            )}
            <Icon
              size={24}
              strokeWidth={2}
              className={active ? "text-purple" : "text-gray-300"}
              fill={active ? "#EDE7F6" : "transparent"}
            />
            <span
              className={`text-label-sm font-semibold ${
                active ? "text-purple" : "text-gray-300"
              }`}
            >
              {label}
            </span>
          </Link>
        );
      })}
    </nav>
  );
}
