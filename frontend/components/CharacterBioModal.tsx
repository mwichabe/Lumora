"use client";

import { motion } from "framer-motion";
import { X, Volume2, Sparkles } from "lucide-react";
import { SpeakerAvatar } from "./Speaker";
import { characterInfo } from "@/lib/characters";
import { speakAs, stopSpeaking } from "@/lib/voices";
import { LANGUAGES } from "@/lib/languages";
import type { CharacterWithFriendship } from "@/lib/types";

// A short greeting each companion can say, per language. The list of buttons is
// driven by which course languages are available, so it grows automatically.
const GREETINGS: Record<string, (name: string) => string> = {
  en: (n) => `Hello! I'm ${n}. Let's learn together!`,
  es: (n) => `¡Hola! Soy ${n}. ¡Vamos a aprender juntos!`,
  de: (n) => `Hallo! Ich bin ${n}. Lass uns zusammen lernen!`,
  fr: (n) => `Bonjour ! Je suis ${n}. Apprenons ensemble !`,
  it: (n) => `Ciao! Sono ${n}. Impariamo insieme!`,
  pt: (n) => `Olá! Eu sou ${n}. Vamos aprender juntos!`,
};

const SPOKEN_LANGS = [
  { code: "en", flag: "🇬🇧", native: "English" },
  ...LANGUAGES.filter((l) => l.available && GREETINGS[l.code]).map((l) => ({
    code: l.code,
    flag: l.flag,
    native: l.native,
  })),
];

/** A tap-to-open bio card for a companion: face, role, voice, personality. */
export function CharacterBioModal({
  character,
  onClose,
}: {
  character: CharacterWithFriendship | null;
  onClose: () => void;
}) {
  if (!character) return null;
  const info = characterInfo(character.name);
  const close = () => {
    stopSpeaking();
    onClose();
  };

  const attributes: { label: string; value: string }[] = [
    { label: "Role", value: character.role },
    { label: "Species", value: character.species },
    {
      label: "Voice",
      value: info.gender === "male" ? "Male" : "Female",
    },
  ];

  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center bg-black/50 p-0 sm:items-center sm:p-4"
      onClick={close}
    >
      <motion.div
        initial={{ y: 40, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        onClick={(e) => e.stopPropagation()}
        className="w-full max-w-md overflow-hidden rounded-t-3xl bg-white shadow-card-lg sm:rounded-3xl"
      >
        {/* header band in the character's colour */}
        <div
          className="relative flex flex-col items-center px-6 pb-5 pt-7"
          style={{ backgroundColor: info.color + "1A" }}
        >
          <button
            onClick={close}
            aria-label="Close"
            className="absolute right-4 top-4 flex h-8 w-8 items-center justify-center rounded-full bg-white/70 text-gray-500 transition hover:bg-white"
          >
            <X size={18} />
          </button>
          <div
            className="rounded-full border-4 p-1"
            style={{ borderColor: info.color }}
          >
            <SpeakerAvatar name={character.name} size={96} />
          </div>
          <h2 className="mt-3 text-heading-xl font-extrabold text-ink">
            {character.name}
          </h2>
          <p className="text-body-md font-semibold" style={{ color: info.color }}>
            {character.role}
          </p>

          {/* friendship */}
          <div className="mt-2 flex items-center gap-1">
            {Array.from({ length: 3 }).map((_, i) => (
              <span
                key={i}
                className={`h-2 w-2 rounded-full ${
                  i < character.friendshipLevel ? "bg-amber" : "bg-gray-300"
                }`}
              />
            ))}
            <span className="ml-1 text-label-md text-slatey">
              Friendship {character.friendshipLevel}/3
            </span>
          </div>
        </div>

        <div className="px-6 py-5">
          {/* quote */}
          {info.quote && (
            <div className="mb-4 flex items-start gap-2 rounded-2xl bg-gray-50 p-3">
              <Sparkles size={16} className="mt-0.5 shrink-0 text-amber" />
              <p className="text-body-md italic text-ink">
                &ldquo;{info.quote}&rdquo;
              </p>
            </div>
          )}

          {/* bio */}
          <p className="text-body-md leading-relaxed text-slatey">
            {character.personality}
          </p>

          {/* attributes */}
          <div className="mt-4 grid grid-cols-3 gap-2">
            {attributes.map((a) => (
              <div key={a.label} className="rounded-xl bg-gray-50 p-2.5 text-center">
                <p className="text-label-sm font-bold uppercase tracking-wide text-gray-400">
                  {a.label}
                </p>
                <p className="mt-0.5 text-body-sm font-extrabold text-ink">
                  {a.value}
                </p>
              </div>
            ))}
          </div>

          {/* hear them speak — one tap per available language */}
          <div className="mt-4">
            <p className="mb-2 flex items-center gap-1.5 text-label-md font-bold uppercase tracking-wide text-gray-400">
              <Volume2 size={14} /> Hear {character.name} speak
            </p>
            <div className="flex flex-wrap gap-2">
              {SPOKEN_LANGS.map((l) => (
                <button
                  key={l.code}
                  onClick={() => {
                    stopSpeaking();
                    speakAs(character.name, GREETINGS[l.code](character.name), {
                      lang: l.code,
                    });
                  }}
                  className="flex items-center gap-1.5 rounded-full border-2 border-gray-100 bg-white px-3 py-2 text-body-sm font-bold text-ink transition hover:border-purple hover:bg-purple-light"
                >
                  <span className="text-base leading-none">{l.flag}</span>
                  {l.native}
                </button>
              ))}
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
