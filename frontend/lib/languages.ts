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
  { code: "zh", name: "Mandarin", native: "中文", flag: "🇨🇳", available: false },
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
