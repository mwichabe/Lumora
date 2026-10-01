export interface LanguageMeta {
  code: string;
  name: string;
  native: string;
  flag: string;
  /** Whether a full course is authored for this language yet. */
  available: boolean;
}

export const LANGUAGES: LanguageMeta[] = [
  { code: "es", name: "Spanish", native: "Español", flag: "🇪🇸", available: true },
  { code: "de", name: "German", native: "Deutsch", flag: "🇩🇪", available: true },
  { code: "fr", name: "French", native: "Français", flag: "🇫🇷", available: true },
  { code: "ja", name: "Japanese", native: "日本語", flag: "🇯🇵", available: false },
  { code: "zh", name: "Mandarin", native: "中文", flag: "🇨🇳", available: true },
  { code: "ar", name: "Arabic", native: "العربية", flag: "🇸🇦", available: false },
  { code: "sw", name: "Swahili", native: "Kiswahili", flag: "🇰🇪", available: false },
  { code: "pt", name: "Portuguese", native: "Português", flag: "🇵🇹", available: false },
  { code: "it", name: "Italian", native: "Italiano", flag: "🇮🇹", available: false },
  { code: "ko", name: "Korean", native: "한국어", flag: "🇰🇷", available: false },
  { code: "hi", name: "Hindi", native: "हिन्दी", flag: "🇮🇳", available: false },
];

export function languageMeta(code?: string): LanguageMeta | undefined {
  return LANGUAGES.find((l) => l.code === code);
}

export function languageName(code?: string): string {
  return languageMeta(code)?.name || "your language";
}

/** Languages with a full course today — the only ones it's honest to advertise. */
export const AVAILABLE_LANGUAGES = LANGUAGES.filter((l) => l.available);

/** "Spanish, German & French" */
export function availableLanguageList(): string {
  const names = AVAILABLE_LANGUAGES.map((l) => l.name);
  return names.length > 1 ? `${names.slice(0, -1).join(", ")} & ${names[names.length - 1]}` : names.join("");
}

/** HSK equivalent of each CEFR level, for Mandarin (HSK 3.0). */
const HSK: Record<string, string> = {
  A1: "HSK 1",
  A2: "HSK 2",
  B1: "HSK 3",
  B2: "HSK 4",
  C1: "HSK 5",
  C2: "HSK 6",
  FINAL: "HSK 7–9",
};

/**
 * How to show an exam / certificate level for a language: Mandarin learners
 * see the HSK level they know ("HSK 1 · A1"); everyone else sees CEFR.
 */
export function levelDisplay(level: string, lang?: string): string {
  if (lang === "zh" && HSK[level]) {
    return level === "FINAL" ? `${HSK[level]} (Advanced)` : `${HSK[level]} · ${level}`;
  }
  return level === "FINAL" ? "Final Mastery" : level;
}
