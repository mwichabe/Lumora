"use client";

/* eslint-disable @next/next/no-img-element */
import { characterInfo } from "@/lib/characters";

/** A character's circular human portrait. */
export function SpeakerAvatar({
  name,
  size = 40,
  className = "",
}: {
  name?: string;
  size?: number;
  className?: string;
}) {
  const info = characterInfo(name);
  return (
    <img
      src={info.img}
      alt={info.name}
      width={size}
      height={size}
      loading="lazy"
      className={`shrink-0 rounded-full bg-white object-cover ${className}`}
      style={{ width: size, height: size }}
    />
  );
}

/** A pill showing the speaker's face + name — used wherever a character speaks. */
export function SpeakerChip({
  name,
  size = 28,
  className = "",
}: {
  name?: string;
  size?: number;
  className?: string;
}) {
  const info = characterInfo(name);
  return (
    <span
      className={`inline-flex items-center gap-2 rounded-full bg-white py-1 pl-1 pr-3 shadow-card ${className}`}
    >
      <SpeakerAvatar name={name} size={size} />
      <span className="text-label-lg font-extrabold text-ink">{info.name}</span>
    </span>
  );
}
