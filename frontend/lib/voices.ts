"use client";

/**
 * Character voices via the browser's built-in Web Speech API (no API key
 * required). Each character gets a distinct profile — a preferred Spanish voice
 * plus a unique pitch/rate — so they sound recognisably different even on
 * machines that only ship one Spanish voice.
 *
 * To upgrade to studio-quality, per-character cloned voices later, swap
 * `speakAs` to call a backend TTS proxy (see the notes shared in chat).
 */

export interface VoiceProfile {
  pitch: number; // 0–2 (1 = default)
  rate: number; // 0.1–10 (1 = default)
  lang: string; // preferred BCP-47 tag
  gender: "female" | "male"; // preferred voice gender
  hints: string[]; // substrings to match a system voice name, in priority order
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
const FEMALE_HINTS = ["female", "mónica", "monica", "paulina", "helena", "elvira", "sabina", "marisol", "lucia", "laura", "mujer"];
const MALE_HINTS = ["male", "jorge", "diego", "juan", "carlos", "enrique", "pablo", "hombre"];

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

// Score a voice for a profile. Higher = better. Prioritises clear, natural
// voices (Google / neural / online, non-local) of the correct gender, which
// fixes thin/robotic/"whispering" output from low-quality local engines.
function scoreVoice(v: SpeechSynthesisVoice, profile: VoiceProfile): number {
  const name = v.name.toLowerCase();
  const lang = (v.lang || "").toLowerCase();
  let s = 0;

  if (lang.startsWith("es")) s += 5;
  if (lang === profile.lang.toLowerCase()) s += 3;

  // Specific name hints, in priority order.
  profile.hints.forEach((h, i) => {
    if (name.includes(h.toLowerCase())) s += 30 - i;
  });

  // Gender: reward the right gender, punish the wrong one.
  const wanted = profile.gender === "male" ? MALE_HINTS : FEMALE_HINTS;
  const opposite = profile.gender === "male" ? FEMALE_HINTS : MALE_HINTS;
  if (wanted.some((g) => name.includes(g))) s += 8;
  if (opposite.some((g) => name.includes(g))) s -= 10;

  // Quality signals — these are the clear, natural-sounding engines.
  if (/(natural|neural|online|wavenet|premium)/.test(name)) s += 12;
  if (name.includes("google")) s += 8;
  if (name.includes("microsoft")) s += 5;
  if (v.localService === false) s += 4; // network voices are usually clearer

  // Penalise known thin/robotic local engines.
  if (/(espeak|festival|pico|compact)/.test(name)) s -= 8;

  return s;
}

function pickVoice(profile: VoiceProfile): SpeechSynthesisVoice | undefined {
  const voices = loadVoices();
  if (!voices.length) return undefined;

  const spanish = voices.filter((v) => v.lang?.toLowerCase().startsWith("es"));
  const pool = spanish.length ? spanish : voices;

  let best: SpeechSynthesisVoice | undefined;
  let bestScore = -Infinity;
  for (const v of pool) {
    const sc = scoreVoice(v, profile);
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

/** Speak a single line in a character's voice. Resolves when finished. */
export function speakAs(
  character: string | undefined,
  text: string,
  opts: { onStart?: () => void; onEnd?: () => void } = {}
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
    const u = new SpeechSynthesisUtterance(text);
    const voice = pickVoice(profile);
    if (voice) u.voice = voice;
    u.lang = voice?.lang || profile.lang;
    u.pitch = profile.pitch;
    u.rate = profile.rate;
    u.volume = 1; // full volume — avoids the faint/"whispering" feel

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

/** Listen to the microphone once and resolve with the recognised transcript. */
export function recognizeSpeech(lang = "es-ES"): Promise<string> {
  return new Promise((resolve, reject) => {
    const Rec = getRecognition();
    if (!Rec) {
      reject(new Error("unsupported"));
      return;
    }
    stopSpeaking();
    const rec = new Rec();
    rec.lang = lang;
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
  shouldContinue?: () => boolean
): Promise<void> {
  for (let i = 0; i < lines.length; i++) {
    if (shouldContinue && !shouldContinue()) return;
    onLine?.(i);
    // small gap between speakers
    await speakAs(lines[i].character, lines[i].text);
    await new Promise((r) => setTimeout(r, 250));
  }
}
