/**
 * Human portrait faces for the companion cast.
 *
 * We deliberately do NOT embed real people's photos scraped from the web —
 * that carries likeness/privacy/licensing risk. Instead we use DiceBear's
 * `avataaars` set (free, CC-BY), which renders a consistent, human portrait per
 * character, deterministically by seed, with gender-appropriate styling that
 * matches each character's voice (male voice ↔ male face, female ↔ female).
 *
 * If you later license real portrait photos, just swap the `img` URLs below.
 */

export type Gender = "female" | "male";

export interface CharacterInfo {
  name: string;
  gender: Gender;
  color: string;
  img: string;
  /** A short in-character line, shown in the bio. */
  quote?: string;
}

// Valid DiceBear v9 `avataaars` hair (top) values.
const FEMALE_TOPS =
  "straight01,straight02,bob,curly,curvy,longButNotTooLong,miaWallace,bigHair";
const MALE_TOPS =
  "shortFlat,shortRound,shortCurly,shortWaved,theCaesar,sides,dreads01,frizzle";

/** Build a deterministic, gender-appropriate human portrait URL. */
function face(seed: string, gender: Gender, color: string, extra = ""): string {
  const tops = gender === "female" ? FEMALE_TOPS : MALE_TOPS;
  const facial =
    gender === "female" ? "facialHairProbability=0" : "facialHairProbability=70";
  const bg = color.replace("#", "");
  const base =
    `https://api.dicebear.com/9.x/avataaars/svg?seed=${encodeURIComponent(seed)}` +
    `&top=${tops}&${facial}&backgroundColor=${bg}&backgroundType=solid&radius=50`;
  return extra ? `${base}&${extra}` : base;
}

// The cast — gender mirrors lib/voices.ts so the face and the voice always agree.
export const CHARACTERS: Record<string, CharacterInfo> = {
  Lumora: {
    name: "Lumora",
    gender: "female",
    color: "#6C3FC5",
    img: face("Lumora-guide-7", "female", "#6C3FC5"),
    quote: "Every word you learn is a little light. Let's make you shine!",
  },
  "Professor Finch": {
    name: "Professor Finch",
    gender: "male",
    color: "#8B6F47",
    img: face(
      "Finch-prof-3",
      "male",
      "#8B6F47",
      "accessories=prescription02&accessoriesProbability=100&hairColor=a9a9a9&facialHair=beardLight"
    ),
    quote: "Grammar is not a cage — it is the architecture of meaning.",
  },
  Cora: {
    name: "Cora",
    gender: "female",
    color: "#00C2A8",
    img: face("Cora-2", "female", "#00C2A8"),
    quote: "Infinite vocabulary, infinite puns. Let's get word-y!",
  },
  Blaze: {
    name: "Blaze",
    gender: "male",
    color: "#FF5C5C",
    img: face("Blaze-5", "male", "#FF5C5C"),
    quote: "MORE FIRE! Say it louder — you've absolutely got this!",
  },
  Mira: {
    name: "Mira",
    gender: "female",
    color: "#9090A0",
    img: face("Mira-4", "female", "#9090A0"),
    quote: "Close your eyes and let the sounds tell their story.",
  },
  Riko: {
    name: "Riko",
    gender: "male",
    color: "#F5A623",
    img: face("Riko-8", "male", "#F5A623"),
    quote: "Think you can beat me? ...Fine, maybe you can. A little.",
  },
  Zephyr: {
    name: "Zephyr",
    gender: "male",
    color: "#17A3DD",
    img: face("Zephyr-1", "male", "#17A3DD"),
    quote: "Words are wind — give them shape and they carry far.",
  },
  Nana: {
    name: "Nana",
    gender: "female",
    color: "#06AECE",
    img: face("Nana-elder-6", "female", "#06AECE", "hairColor=b7b7b7"),
    quote: "Slow and steady. Every journey takes the time it needs.",
  },
  Pip: {
    name: "Pip",
    gender: "male",
    color: "#F5A623",
    img: face("Pip-9", "male", "#F5A623"),
    quote: "Quick-quick! So many quests, so little time — let's go!",
  },
};

/** Resolve a speaker name to a face/gender, with a graceful fallback. */
export function characterInfo(name?: string): CharacterInfo {
  if (name && CHARACTERS[name]) return CHARACTERS[name];
  const safe = name && name.trim() ? name : "Lumora";
  // Unknown speaker (e.g. a one-off dialogue name): give them a stable face.
  return {
    name: safe,
    gender: "female",
    color: "#6C3FC5",
    img: face(safe, "female", "#6C3FC5"),
  };
}
