"use client";

import { useEffect, useState } from "react";
import { mediaUrl } from "@/lib/api";

/**
 * Shows a user's uploaded profile photo, falling back to a coloured initial.
 *
 * The fallback covers two cases, and the second is the one that bites: a user
 * with no photo at all, and a user whose photo URL no longer resolves. Profile
 * photos used to be written to disk and referenced as /uploads/avatars/…;
 * they now live in the database behind /api/avatars/:id, and nothing serves the
 * old path. Without an onError handler those accounts render the browser's
 * broken-image glyph next to the alt text, which looks like a bug in the app
 * rather than a stale row.
 *
 * (The rows themselves are repaired at boot — see healLegacyAvatars in
 * backend/database/database.go. This is the belt to that braces: any future
 * unreachable URL degrades to something that still looks deliberate.)
 */
export function Avatar({
  name,
  color,
  url,
  size = 40,
  className = "",
}: {
  name?: string;
  color?: string;
  url?: string;
  size?: number;
  className?: string;
}) {
  const src = mediaUrl(url);
  const [failed, setFailed] = useState(false);

  // A changed URL deserves a fresh attempt — otherwise uploading a new photo
  // would keep showing the initial for the rest of the session.
  useEffect(() => setFailed(false), [src]);

  const initial = (name || "L").slice(0, 1).toUpperCase();

  if (src && !failed) {
    // eslint-disable-next-line @next/next/no-img-element
    return (
      <img
        src={src}
        alt={name || "Profile"}
        width={size}
        height={size}
        onError={() => setFailed(true)}
        style={{ width: size, height: size, backgroundColor: color || "#6C3FC5" }}
        className={`shrink-0 rounded-full object-cover ${className}`}
      />
    );
  }

  return (
    <span
      aria-label={name || "Profile"}
      style={{
        width: size,
        height: size,
        backgroundColor: color || "#6C3FC5",
        fontSize: size * 0.4,
      }}
      className={`flex shrink-0 select-none items-center justify-center rounded-full font-extrabold leading-none text-white ${className}`}
    >
      {initial}
    </span>
  );
}
