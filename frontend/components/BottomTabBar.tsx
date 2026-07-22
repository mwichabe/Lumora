"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { AnimatePresence, motion } from "framer-motion";
import {
  Bell,
  BookOpen,
  ChevronRight,
  Home,
  Lightbulb,
  MessageCircle,
  Mic,
  MoreHorizontal,
  Trophy,
  User,
  X,
} from "lucide-react";
import { useUnreadCount } from "@/lib/notifications";
import { useChatUnread } from "@/lib/chat";

/**
 * The phone navigation.
 *
 * There are nine destinations and room for about five. Cramming them all in
 * gives five-point tap targets and unreadable labels, so the four the app is
 * actually *for* stay on the bar and the rest live one tap away behind More.
 *
 * The compromise only works if nothing important can hide there, so the More
 * tab carries a badge for anything unread inside it — you can see there's a
 * message waiting without opening the sheet to find out.
 */

const tabs = [
  { href: "/home", label: "Home", icon: Home },
  { href: "/learn", label: "Learn", icon: BookOpen },
  { href: "/practice", label: "Practice", icon: Mic },
  { href: "/leaderboard", label: "Leagues", icon: Trophy },
];

/** Everything reachable from the More sheet. */
const moreRoutes = ["/ideas", "/chat", "/notifications", "/profile"];

export function BottomTabBar() {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);

  const notifications = useUnreadCount();
  const messages = useChatUnread();
  const hidden = notifications + messages;

  // Close the sheet on navigation, so it never lingers over the new screen.
  useEffect(() => setOpen(false), [pathname]);

  // The sheet is a modal surface; the page behind it shouldn't scroll.
  useEffect(() => {
    if (!open) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = prev;
    };
  }, [open]);

  const inMore = moreRoutes.some((r) => pathname?.startsWith(r));

  return (
    <>
      <AnimatePresence>
        {open && (
          <MoreSheet
            onClose={() => setOpen(false)}
            notifications={notifications}
            messages={messages}
          />
        )}
      </AnimatePresence>

      <nav className="sticky bottom-0 z-40 mt-auto flex h-16 items-stretch border-t border-gray-100 bg-white pb-[env(safe-area-inset-bottom)] lg:hidden">
        {tabs.map(({ href, label, icon: Icon }) => (
          <Tab
            key={href}
            href={href}
            label={label}
            icon={Icon}
            active={!!pathname?.startsWith(href) && !open}
          />
        ))}

        <button
          onClick={() => setOpen((o) => !o)}
          aria-label="More"
          aria-expanded={open}
          className="relative flex flex-1 flex-col items-center justify-center gap-0.5"
        >
          {(inMore || open) && (
            <motion.span
              layoutId="tab-dot"
              className="absolute top-1.5 h-1.5 w-1.5 rounded-full bg-purple"
            />
          )}
          <span className="relative">
            <MoreHorizontal
              size={24}
              strokeWidth={2}
              className={inMore || open ? "text-purple" : "text-gray-300"}
            />
            {/* Nothing urgent may hide behind More. */}
            {hidden > 0 && !open && (
              <span className="absolute -right-2 -top-1 flex h-4 min-w-[16px] items-center justify-center rounded-full bg-coral px-1 text-[10px] font-extrabold leading-none text-white">
                {hidden > 9 ? "9+" : hidden}
              </span>
            )}
          </span>
          <span
            className={`text-label-sm font-semibold ${
              inMore || open ? "text-purple" : "text-gray-300"
            }`}
          >
            More
          </span>
        </button>
      </nav>
    </>
  );
}

function Tab({
  href,
  label,
  icon: Icon,
  active,
}: {
  href: string;
  label: string;
  icon: typeof Home;
  active: boolean;
}) {
  return (
    <Link
      href={href}
      className="relative flex flex-1 flex-col items-center justify-center gap-0.5"
    >
      {active && (
        <motion.span
          layoutId="tab-dot"
          className="absolute top-1.5 h-1.5 w-1.5 rounded-full bg-purple"
        />
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
}

function MoreSheet({
  onClose,
  notifications,
  messages,
}: {
  onClose: () => void;
  notifications: number;
  messages: number;
}) {
  const pathname = usePathname();

  const items = [
    {
      href: "/ideas",
      label: "Ideas",
      hint: "Propose, vote, discuss",
      icon: Lightbulb,
      tint: "#F5A623",
      badge: 0,
    },
    {
      href: "/chat",
      label: "Messages",
      hint: "Chat with other learners",
      icon: MessageCircle,
      tint: "#6C3FC5",
      badge: messages,
    },
    {
      href: "/notifications",
      label: "Notifications",
      hint: "League results, milestones",
      icon: Bell,
      tint: "#00C2A8",
      badge: notifications,
    },
    {
      href: "/profile",
      label: "Profile",
      hint: "Certificates, settings",
      icon: User,
      tint: "#17A3DD",
      badge: 0,
    },
  ];

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      onClick={onClose}
      className="fixed inset-0 z-30 flex items-end bg-black/40 backdrop-blur-[2px] lg:hidden"
    >
      <motion.div
        initial={{ y: "100%" }}
        animate={{ y: 0 }}
        exit={{ y: "100%" }}
        transition={{ type: "spring", stiffness: 320, damping: 32 }}
        onClick={(e) => e.stopPropagation()}
        // Clears the tab bar so the sheet never sits under it.
        className="w-full rounded-t-[28px] bg-white pb-[calc(4rem+env(safe-area-inset-bottom))] shadow-card-lg"
        role="dialog"
        aria-label="More destinations"
      >
        <div className="flex items-center justify-between px-5 pb-2 pt-3">
          <span className="mx-auto h-1 w-10 rounded-full bg-gray-100" />
        </div>

        <div className="flex items-center justify-between px-5 pb-1">
          <h2 className="text-heading-md font-extrabold text-ink">More</h2>
          <button
            onClick={onClose}
            aria-label="Close"
            className="rounded-full p-1.5 text-slatey transition hover:bg-gray-50"
          >
            <X size={18} />
          </button>
        </div>

        <div className="space-y-1 p-3 pt-1">
          {items.map((item, i) => {
            const active = pathname?.startsWith(item.href);
            return (
              <motion.div
                key={item.href}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.04 + i * 0.045 }}
              >
                <Link
                  href={item.href}
                  onClick={onClose}
                  className={`flex items-center gap-3 rounded-2xl px-3 py-3 transition ${
                    active ? "bg-purple-light" : "hover:bg-gray-50"
                  }`}
                >
                  <span
                    className="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl"
                    style={{ background: `${item.tint}1A`, color: item.tint }}
                  >
                    <item.icon size={21} strokeWidth={2.2} />
                  </span>
                  <span className="min-w-0 flex-1">
                    <span className="flex items-center gap-2">
                      <span
                        className={`font-extrabold ${
                          active ? "text-purple" : "text-ink"
                        }`}
                      >
                        {item.label}
                      </span>
                      {item.badge > 0 && (
                        <span className="flex h-5 min-w-[20px] items-center justify-center rounded-full bg-coral px-1 text-label-sm font-extrabold text-white">
                          {item.badge > 9 ? "9+" : item.badge}
                        </span>
                      )}
                    </span>
                    <span className="block truncate text-label-md text-slatey">
                      {item.hint}
                    </span>
                  </span>
                  <ChevronRight size={17} className="shrink-0 text-gray-300" />
                </Link>
              </motion.div>
            );
          })}
        </div>
      </motion.div>
    </motion.div>
  );
}
