import { mediaUrl } from "@/lib/api";

/**
 * Shows a user's uploaded profile photo, falling back to a coloured initial
 * when none is set.
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
  const initial = (name || "L").slice(0, 1).toUpperCase();

  if (src) {
    // eslint-disable-next-line @next/next/no-img-element
    return (
      <img
        src={src}
        alt={name || "Profile"}
        width={size}
        height={size}
        style={{ width: size, height: size }}
        className={`shrink-0 rounded-full object-cover ${className}`}
      />
    );
  }

  return (
    <span
      style={{ width: size, height: size, backgroundColor: color || "#6C3FC5", fontSize: size * 0.4 }}
      className={`flex shrink-0 items-center justify-center rounded-full font-extrabold text-white ${className}`}
    >
      {initial}
    </span>
  );
}
