"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ChevronLeft, ChevronRight, Search } from "lucide-react";

const LANGUAGES = [
  { code: "es", flag: "🇪🇸", name: "Spanish", native: "Español" },
  { code: "fr", flag: "🇫🇷", name: "French", native: "Français" },
  { code: "ja", flag: "🇯🇵", name: "Japanese", native: "日本語" },
  { code: "zh", flag: "🇨🇳", name: "Mandarin", native: "中文" },
  { code: "ar", flag: "🇸🇦", name: "Arabic", native: "العربية" },
  { code: "sw", flag: "🇰🇪", name: "Swahili", native: "Kiswahili" },
  { code: "pt", flag: "🇵🇹", name: "Portuguese", native: "Português" },
  { code: "de", flag: "🇩🇪", name: "German", native: "Deutsch" },
  { code: "it", flag: "🇮🇹", name: "Italian", native: "Italiano" },
  { code: "ko", flag: "🇰🇷", name: "Korean", native: "한국어" },
  { code: "hi", flag: "🇮🇳", name: "Hindi", native: "हिन्दी" },
];

export default function ChooseLanguageScreen() {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [selected, setSelected] = useState<string | null>(null);

  const filtered = LANGUAGES.filter((l) =>
    `${l.name} ${l.native}`.toLowerCase().includes(query.toLowerCase())
  );

  function choose(code: string) {
    setSelected(code);
    // Spanish has full seeded content; others are previews in the MVP.
    sessionStorage.setItem("lumora_target", code);
    setTimeout(() => router.push("/onboarding/goal"), 220);
  }

  return (
    <div className="app-frame flex flex-col bg-cream">
      <header className="flex h-14 items-center px-4">
        <button onClick={() => router.back()} aria-label="Back">
          <ChevronLeft className="text-purple" />
        </button>
        <h1 className="flex-1 text-center text-heading-md font-bold">Choose a Language</h1>
        <span className="w-6" />
      </header>

      <div className="px-5">
        <div className="flex items-center gap-2 rounded-sm bg-gray-50 px-3 py-2.5">
          <Search size={18} className="text-gray-500" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search languages…"
            className="w-full bg-transparent text-body-md outline-none"
          />
        </div>
      </div>

      <div className="mt-3 flex-1 overflow-y-auto px-3 pb-6">
        {filtered.map((l) => {
          const active = selected === l.code;
          return (
            <button
              key={l.code}
              onClick={() => choose(l.code)}
              className={`flex h-[72px] w-full items-center gap-3 rounded-md px-3 text-left transition ${
                active ? "border-l-4 border-purple bg-purple-light" : "border-l-4 border-transparent"
              }`}
            >
              <span className="text-3xl">{l.flag}</span>
              <span className="flex-1">
                <span className="block text-body-lg font-bold">{l.name}</span>
                <span className="block text-body-sm text-slatey">{l.native}</span>
              </span>
              {l.code === "es" && (
                <span className="rounded-full bg-teal-light px-2 py-0.5 text-label-sm font-bold text-teal">
                  Full course
                </span>
              )}
              <ChevronRight size={18} className="text-gray-300" />
            </button>
          );
        })}
      </div>
    </div>
  );
}
