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
  hints: string[]; // substrings to match a system voice name, in priority order
}

const DEFAULT_PROFILE: VoiceProfile = {
  pitch: 1,
  rate: 0.95,
  lang: "es-ES",
  hints: ["spanish", "español", "google"],
};

export const CHARACTER_VOICES: Record<string, VoiceProfile> = {
  Lumora: { pitch: 1.2, rate: 1.0, lang: "es-ES", hints: ["mónica", "monica", "female", "google español"] },
  "Professor Finch": { pitch: 0.7, rate: 0.82, lang: "es-ES", hints: ["jorge", "diego", "male"] },
  Cora: { pitch: 1.55, rate: 1.12, lang: "es-MX", hints: ["paulina", "female"] },
  Blaze: { pitch: 1.05, rate: 1.22, lang: "es-MX", hints: ["juan", "male"] },
  Mira: { pitch: 0.95, rate: 0.8, lang: "es-ES", hints: ["mónica", "monica", "female"] },
  Riko: { pitch: 0.82, rate: 1.08, lang: "es-AR", hints: ["male", "diego"] },
  Zephyr: { pitch: 0.78, rate: 0.9, lang: "es-ES", hints: ["male", "jorge"] },
  Nana: { pitch: 0.72, rate: 0.68, lang: "es-ES", hints: ["female", "mónica"] },
  Pip: { pitch: 1.6, rate: 1.4, lang: "es-MX", hints: ["male", "juan"] },
};

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

function pickVoice(profile: VoiceProfile): SpeechSynthesisVoice | undefined {
  const voices = loadVoices();
  if (!voices.length) return undefined;

  const spanish = voices.filter((v) => v.lang?.toLowerCase().startsWith("es"));
  const pool = spanish.length ? spanish : voices;

  // Prefer an exact-ish locale + name-hint match, then any hint, then locale.
  for (const hint of profile.hints) {
    const h = hint.toLowerCase();
    const match = pool.find(
      (v) =>
        v.name.toLowerCase().includes(h) &&
        v.lang.toLowerCase().startsWith(profile.lang.slice(0, 2))
    );
    if (match) return match;
  }
  const byLocale = pool.find((v) => v.lang.toLowerCase() === profile.lang.toLowerCase());
  return byLocale || pool[0];
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
