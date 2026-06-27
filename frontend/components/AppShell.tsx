"use client";

import { ReactNode, useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth";
import { setLastRoute } from "@/lib/api";
import { BottomTabBar } from "./BottomTabBar";
import { SideNav } from "./SideNav";
import { FoxMascot } from "./FoxMascot";

/**
 * AppShell wraps authenticated screens. It is responsive:
 *  • On phones it renders a single column with the BottomTabBar.
 *  • On desktops (lg+) it renders a left SideNav rail and a centred, wider
 *    content column so the app uses the available space instead of staying
 *    locked to a phone-width strip.
 *
 * Pass `wide` for screens that benefit from a roomier multi-column layout
 * (e.g. the home dashboard and profile).
 */
export function AppShell({
  children,
  tabs = true,
  wide = false,
}: {
  children: ReactNode;
  tabs?: boolean;
  wide?: boolean;
}) {
  const { user, loading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  // Remember where the user is so the splash can restore this exact screen if
  // they reload or manually edit the URL.
  useEffect(() => {
    if (user && pathname) setLastRoute(pathname);
  }, [user, pathname]);

  // Only send to the splash once we're certain there's no session — never while
  // still loading, and never when a cached user is present.
  useEffect(() => {
    if (!loading && !user) router.replace("/");
  }, [loading, user, router]);

  if (loading || !user) {
    return (
      <div className="flex min-h-[100dvh] flex-col items-center justify-center gap-4 bg-cream">
        <FoxMascot size={120} glow />
        <p className="text-body-md text-slatey">Loading your world…</p>
      </div>
    );
  }

  return (
    <div className="min-h-[100dvh] bg-[#eceaf3] lg:flex">
      {tabs && <SideNav />}

      <div className="relative flex min-h-[100dvh] w-full flex-1 flex-col lg:min-w-0">
        <main
          className={`mx-auto flex w-full flex-1 flex-col bg-cream lg:my-6 lg:overflow-hidden lg:rounded-3xl lg:shadow-card-lg ${
            wide ? "lg:max-w-5xl" : "lg:max-w-3xl"
          }`}
        >
          {children}
        </main>
        {tabs && <BottomTabBar />}
      </div>
    </div>
  );
}
