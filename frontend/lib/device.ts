"use client";

import { useEffect, useState } from "react";

/**
 * Whether this browser is running on a phone or tablet. Exams need a laptop or
 * desktop (the proctoring screen-share isn't available on handhelds, and a
 * phone screen invites a second device off-camera), so the exam page uses this
 * to refuse to start. The backend applies the same rule from the user agent
 * (backend/middleware/web_only.go); this side also catches iPads, whose Safari
 * reports itself as a Mac.
 */
export function isHandheldDevice(): boolean {
  if (typeof navigator === "undefined") return false;
  const nav = navigator as Navigator & { userAgentData?: { mobile?: boolean } };

  // Client hints, where supported (Chromium): the browser says so outright.
  if (nav.userAgentData?.mobile) return true;

  const ua = nav.userAgent || "";
  if (
    /android|iphone|ipad|ipod|mobile|tablet|silk|kindle|playbook|bb10|opera mini|iemobile|windows phone/i.test(
      ua
    )
  ) {
    return true;
  }
  // iPadOS Safari claims to be a Mac; a real Mac has no touch screen.
  if (/macintosh/i.test(ua) && nav.maxTouchPoints > 1) return true;

  return false;
}

/** `null` until checked on the client (it can't be known during SSR). */
export function useIsHandheld(): boolean | null {
  const [handheld, setHandheld] = useState<boolean | null>(null);
  useEffect(() => setHandheld(isHandheldDevice()), []);
  return handheld;
}
