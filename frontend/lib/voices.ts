"use client";

/**
 * Character voices via the browser's built-in Web Speech API (no API key
 * required). Each character gets a distinct profile — a unique pitch/rate plus
 * voice-name hints — so they sound recognisably different.
 *
 * IMPORTANT: pronunciation must match the CONTENT language. The active learning
 * language is set via `setSpeechLanguage()` (wired to the user's target
 * language), so German text is spoken with a German voice, French with French,
 * etc. The per-character pitch/rate still differentiates speakers within any
 * language. To upgrade to studio-quality cloned voices later, swap `speakAs` to
 * call a backend TTS proxy.
 */

export interface VoiceProfile {
  pitch: number; // 0–2 (1 = default)
  rate: number; // 0.1–10 (1 = default)
  lang: string; // preferred BCP-47 tag (region nudge within the active language)
  gender: "female" | "male"; // preferred voice gender
  hints: string[]; // substrings to match a system voice name, in priority order
}

// Map a language code to a sensible default locale for speech.
const LANG_LOCALE: Record<string, string> = {
  es: "es-ES",
  de: "de-DE",
  fr: "fr-FR",
  en: "en-US",
  it: "it-IT",
  pt: "pt-PT",
  ja: "ja-JP",
  zh: "zh-CN",
  ar: "ar-SA",
  sw: "sw-KE",
};

/** Normalise a language code or locale into a BCP-47 locale (e.g. "de" → "de-DE"). */
function toLocale(x?: string): string {
  if (!x) return activeLocale;
  if (x.includes("-")) return x;
  return LANG_LOCALE[x] || x;
}

// The active learning language drives which system voice is used. Defaults to
// Spanish for backward compatibility until setSpeechLanguage() is called.
let activeLocale = "es-ES";

/** Set the language all subsequent speech/recognition should use. */
export function setSpeechLanguage(code?: string) {
  if (!code) return;
  activeLocale = toLocale(code);
}

/** The locale currently used for speech (e.g. for SpeechRecognition). */
export function getSpeechLocale(): string {
  return activeLocale;
}

const DEFAULT_PROFILE: VoiceProfile = {
  pitch: 1,
  rate: 1.0,
  lang: "es-ES",
  gender: "female",
  hints: ["google español", "español", "spanish"],
};

// Pitches are kept in a natural range (no thin/breathy extremes). Characters are
// differentiated by GENDER + a distinct VOICE pick first, then pitch/rate, so
// they stay distinct even when the OS ships only one Spanish voice.
export const CHARACTER_VOICES: Record<string, VoiceProfile> = {
  // Warm, clear, friendly — medium female.
  Lumora: { pitch: 1.05, rate: 1.0, lang: "es-ES", gender: "female", hints: ["google español", "mónica", "monica", "helena", "elvira", "sabina", "paulina"] },
  // Bubbly and quick, noticeably brighter & faster than Lumora — Latin-American female.
  Cora: { pitch: 1.28, rate: 1.22, lang: "es-US", gender: "female", hints: ["google español de estados unidos", "paulina", "us"] },
  // Deep, deliberate male.
  "Professor Finch": { pitch: 0.8, rate: 0.88, lang: "es-ES", gender: "male", hints: ["jorge", "google español", "diego"] },
  // Energetic male.
  Blaze: { pitch: 1.1, rate: 1.2, lang: "es-US", gender: "male", hints: ["juan", "google español de estados unidos"] },
  // Calm, measured female.
  Mira: { pitch: 0.96, rate: 0.86, lang: "es-ES", gender: "female", hints: ["mónica", "monica", "elvira"] },
  // Cocky mid male.
  Riko: { pitch: 0.9, rate: 1.08, lang: "es-MX", gender: "male", hints: ["diego", "juan"] },
  // Smooth, poetic male.
  Zephyr: { pitch: 0.86, rate: 0.94, lang: "es-ES", gender: "male", hints: ["jorge", "diego"] },
  // Gentle, slow elder female.
  Nana: { pitch: 0.84, rate: 0.74, lang: "es-ES", gender: "female", hints: ["mónica", "elvira"] },
  // Fast, excitable male (not shrill).
  Pip: { pitch: 1.2, rate: 1.35, lang: "es-MX", gender: "male", hints: ["juan", "diego"] },
};

// Common female/male voice-name fragments used to honour `gender` when picking.
// Spans several languages so gender matching also works for German/French/English
// system voices, not just Spanish.
const FEMALE_HINTS = [
  "female", "mujer", "femme", "frau", "weiblich",
  // Spanish
  "mónica", "monica", "paulina", "helena", "elvira", "sabina", "marisol", "lucia", "laura",
  // German
  "katja", "hedda", "marlene", "gisela", "vicki", "amala", "petra", "anna",
  // French
  "amelie", "amélie", "audrey", "julie", "hortense", "céline", "celine", "denise", "léa",
  // English
  "zira", "samantha", "susan", "karen", "fiona", "moira", "tessa", "serena", "aria", "jenny", "michelle",
];
const MALE_HINTS = [
  "male", "hombre", "homme", "mann", "männlich",
  // Spanish
  "jorge", "diego", "juan", "carlos", "enrique", "pablo",
  // German
  "conrad", "stefan", "hans", "klaus", "bernd", "yannick",
  // French
  "thomas", "paul", "claude", "henri", "nicolas", "mathieu",
  // English
  "david", "mark", "george", "daniel", "alex", "fred", "guy", "ryan", "eric", "brian",
];

export function profileFor(character?: string): VoiceProfile {
  return (character && CHARACTER_VOICES[character]) || DEFAULT_PROFILE;
}

function synth(): SpeechSynthesis | null {
  if (typeof window === "undefined") return null;
  return window.speechSynthesis || null;
}

let cachedVoices: SpeechSynthesisVoice[] = [];

function loadVoices(): SpeechSynthesisVoice[] {
  const s = synth();
  if (!s) return [];
  const v = s.getVoices();
  if (v.length) cachedVoices = v;
  return cachedVoices;
}

// Voices populate asynchronously in some browsers; prime the cache.
if (typeof window !== "undefined" && window.speechSynthesis) {
  loadVoices();
  window.speechSynthesis.onvoiceschanged = () => loadVoices();
}

// Score a voice for a profile at a given locale. Higher = better. The CONTENT
// language is by far the most important factor (so German is spoken by a German
// voice); after that we prefer clear, natural voices of the correct gender.
function scoreVoice(
  v: SpeechSynthesisVoice,
  profile: VoiceProfile,
  locale: string
): number {
  const name = v.name.toLowerCase();
  const vlang = (v.lang || "").toLowerCase().replace("_", "-");
  const base = locale.split("-")[0].toLowerCase();
  let s = 0;

  // Matching the spoken language is critical for correct pronunciation.
  if (vlang.startsWith(base)) s += 40;
  if (vlang === locale.toLowerCase()) s += 8; // exact region match
  if (vlang === profile.lang.toLowerCase()) s += 3; // character's region nudge

  // Specific name hints, in priority order.
  profile.hints.forEach((h, i) => {
    if (name.includes(h.toLowerCase())) s += 20 - i;
  });

  // Gender: strongly reward the right gender, heavily punish the wrong one, so a
  // male character never ends up on a clearly-female named voice when a male one
  // exists for the language.
  const wanted = profile.gender === "male" ? MALE_HINTS : FEMALE_HINTS;
  const opposite = profile.gender === "male" ? FEMALE_HINTS : MALE_HINTS;
  if (wanted.some((g) => name.includes(g))) s += 25;
  if (opposite.some((g) => name.includes(g))) s -= 30;

  // Quality signals — these are the clear, natural-sounding engines.
  if (/(natural|neural|online|wavenet|premium)/.test(name)) s += 12;
  if (name.includes("google")) s += 8;
  if (name.includes("microsoft")) s += 5;
  if (v.localService === false) s += 4; // network voices are usually clearer

  // Penalise known thin/robotic local engines.
  if (/(espeak|festival|pico|compact)/.test(name)) s -= 8;

  return s;
}

/** Web Speech pitch is valid in [0, 2]; keep it in a usable, non-robotic range. */
function clampPitch(p: number): number {
  return Math.max(0.3, Math.min(2, p));
}

/** Guess a system voice's gender from its name; "unknown" when not encoded. */
export function voiceGender(
  v?: SpeechSynthesisVoice
): "male" | "female" | "unknown" {
  if (!v) return "unknown";
  const n = v.name.toLowerCase();
  if (MALE_HINTS.some((g) => n.includes(g))) return "male";
  if (FEMALE_HINTS.some((g) => n.includes(g))) return "female";
  return "unknown";
}

function pickVoice(
  profile: VoiceProfile,
  locale: string
): SpeechSynthesisVoice | undefined {
  const voices = loadVoices();
  if (!voices.length) return undefined;

  // Only consider voices in the content language when any are installed, so we
  // never read German text with a Spanish/English voice.
  const base = locale.split("-")[0].toLowerCase();
  const matching = voices.filter((v) =>
    (v.lang || "").toLowerCase().replace("_", "-").startsWith(base)
  );
  const langPool = matching.length ? matching : voices;

  // Prefer voices of the requested gender; otherwise gender-neutral names;
  // never fall onto an opposite-gender named voice unless nothing else exists.
  const sameGender = langPool.filter((v) => voiceGender(v) === profile.gender);
  const neutral = langPool.filter((v) => voiceGender(v) === "unknown");
  const pool = sameGender.length
    ? sameGender
    : neutral.length
    ? neutral
    : langPool;

  let best: SpeechSynthesisVoice | undefined;
  let bestScore = -Infinity;
  for (const v of pool) {
    const sc = scoreVoice(v, profile, locale);
    if (sc > bestScore) {
      bestScore = sc;
      best = v;
    }
  }
  return best;
}

/** Debug helper: run `debugVoices()` in the browser console to list what's
 *  installed (name + lang + local/remote) so we can tune the picker. */
export function debugVoices(): { name: string; lang: string; local: boolean }[] {
  return loadVoices().map((v) => ({
    name: v.name,
    lang: v.lang,
    local: v.localService,
  }));
}

export function stopSpeaking() {
  synth()?.cancel();
}

/** Speak a single line in a character's voice. Resolves when finished. Pass
 *  `opts.lang` to override the active language for this utterance. */
export function speakAs(
  character: string | undefined,
  text: string,
  opts: { onStart?: () => void; onEnd?: () => void; lang?: string } = {}
): Promise<void> {
  return new Promise((resolve) => {
    const s = synth();
    if (!s || !text) {
      opts.onEnd?.();
      resolve();
      return;
    }
    s.cancel();

    const profile = profileFor(character);
    const locale = toLocale(opts.lang);
    const u = new SpeechSynthesisUtterance(text);
    const voice = pickVoice(profile, locale);
    if (voice) u.voice = voice;
    u.lang = voice?.lang || locale;
    u.rate = profile.rate;
    u.volume = 1; // full volume — avoids the faint/"whispering" feel

    // Make gender unmistakable across every language. Per-language system voices
    // often have no gender in their name ("unknown") and may sound male or
    // female depending on the OS — so we can't rely on the voice alone. Whenever
    // the chosen voice isn't *confirmed* to match the character's gender, we bend
    // the pitch to it: deep for men, high for women. (Per-character base pitch
    // keeps them distinct.) A voice already confirmed to match is left natural.
    const vg = voiceGender(voice);
    let pitch = profile.pitch;
    if (profile.gender === "male" && vg !== "male") {
      pitch = clampPitch(0.5 + (profile.pitch - 0.8) * 0.25); // ~0.50–0.60, deep
    } else if (profile.gender === "female" && vg !== "female") {
      pitch = clampPitch(1.4 + (profile.pitch - 0.8) * 0.4); // ~1.42–1.60, high
    }
    u.pitch = pitch;

    u.onstart = () => opts.onStart?.();
    u.onend = () => {
      opts.onEnd?.();
      resolve();
    };
    u.onerror = () => {
      opts.onEnd?.();
      resolve();
    };
    s.speak(u);
  });
}

// --- Speech recognition (pronunciation / fluency scoring) -------------------

/* eslint-disable @typescript-eslint/no-explicit-any */
function getRecognition(): any {
  if (typeof window === "undefined") return null;
  const w = window as any;
  return w.SpeechRecognition || w.webkitSpeechRecognition || null;
}

export function speechRecognitionSupported(): boolean {
  return !!getRecognition();
}

/** Listen to the microphone once and resolve with the recognised transcript.
 *  Defaults to the active learning language so scoring matches what's spoken. */
export function recognizeSpeech(lang?: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const Rec = getRecognition();
    if (!Rec) {
      reject(new Error("unsupported"));
      return;
    }
    stopSpeaking();
    const rec = new Rec();
    rec.lang = toLocale(lang);
    rec.interimResults = false;
    rec.maxAlternatives = 1;
    let done = false;
    rec.onresult = (e: any) => {
      done = true;
      resolve(e.results?.[0]?.[0]?.transcript ?? "");
    };
    rec.onerror = (e: any) => reject(new Error(e?.error || "error"));
    rec.onend = () => {
      if (!done) resolve("");
    };
    try {
      rec.start();
    } catch {
      reject(new Error("start-failed"));
    }
  });
}

function normalize(s: string): string {
  return s
    .toLowerCase()
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "") // strip accents for fair comparison
    .replace(/[.,!¡¿?"]/g, "")
    .trim();
}

/**
 * Score how closely the spoken phrase matched the target, 0–100. Combines a
 * word-overlap measure with a character-level similarity so partial credit is
 * fair for beginners.
 */
export function scorePronunciation(target: string, said: string): number {
  const t = normalize(target);
  const s = normalize(said);
  if (!s) return 0;
  if (t === s) return 100;

  const tWords = t.split(/\s+/).filter(Boolean);
  const sWords = new Set(s.split(/\s+/).filter(Boolean));
  const hit = tWords.filter((w) => sWords.has(w)).length;
  const wordScore = tWords.length ? hit / tWords.length : 0;

  // character similarity (1 - normalized edit distance)
  const charScore = 1 - editDistance(t, s) / Math.max(t.length, s.length, 1);

  return Math.round((wordScore * 0.6 + charScore * 0.4) * 100);
}

function editDistance(a: string, b: string): number {
  const m = a.length;
  const n = b.length;
  const dp = Array.from({ length: m + 1 }, () => new Array(n + 1).fill(0));
  for (let i = 0; i <= m; i++) dp[i][0] = i;
  for (let j = 0; j <= n; j++) dp[0][j] = j;
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      dp[i][j] = Math.min(
        dp[i - 1][j] + 1,
        dp[i][j - 1] + 1,
        dp[i - 1][j - 1] + (a[i - 1] === b[j - 1] ? 0 : 1)
      );
    }
  }
  return dp[m][n];
}

export interface SpokenLine {
  character?: string;
  text: string;
}

/**
 * Speak a dialogue line by line, each in its character's voice. `onLine` fires
 * with the index as each line starts (for highlighting). Returns a stop fn via
 * `signal`-like cancellation: call stopSpeaking() to abort.
 */
export async function speakSequence(
  lines: SpokenLine[],
  onLine?: (index: number) => void,
  shouldContinue?: () => boolean,
  lang?: string
): Promise<void> {
  for (let i = 0; i < lines.length; i++) {
    if (shouldContinue && !shouldContinue()) return;
    onLine?.(i);
    // small gap between speakers
    await speakAs(lines[i].character, lines[i].text, { lang });
    await new Promise((r) => setTimeout(r, 250));
  }
}
