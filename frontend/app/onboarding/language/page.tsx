"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { ChevronLeft, ChevronRight, Search } from "lucide-react";
import { LANGUAGES } from "@/lib/languages";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";

function ChooseLanguageInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const adding = searchParams.get("add") === "1";
  const { setUser } = useAuth();
  const [query, setQuery] = useState("");
  const [selected, setSelected] = useState<string | null>(null);

  const filtered = LANGUAGES.filter((l) =>
    `${l.name} ${l.native}`.toLowerCase().includes(query.toLowerCase())
  );

  async function choose(code: string) {
    setSelected(code);
    if (adding) {
      // Existing learner enrolling in an additional language → straight in.
      try {
        const r = await api.enrollLanguage(code);
        setUser(r.user);
      } catch {
        /* ignore */
      }
      router.push("/learn");
      return;
    }
    // First-time onboarding → continue to the daily-goal step.
    sessionStorage.setItem("lumora_target", code);
    setTimeout(() => router.push("/onboarding/goal"), 220);
  }

  return (
    <div className="mx-auto flex min-h-[100dvh] w-full max-w-3xl flex-col bg-cream lg:my-8 lg:min-h-[calc(100dvh-4rem)] lg:rounded-3xl lg:shadow-card-lg">
      <header className="flex h-14 items-center px-4 lg:pt-4">
        <button onClick={() => router.back()} aria-label="Back">
          <ChevronLeft className="text-purple" />
        </button>
        <h1 className="flex-1 text-center text-heading-md font-bold">
          {adding ? "Add a Language" : "Choose a Language"}
        </h1>
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

      <div className="mt-3 flex-1 overflow-y-auto px-3 pb-6 lg:grid lg:grid-cols-2 lg:gap-x-3 lg:gap-y-1 lg:content-start">
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
              {l.available ? (
                <span className="rounded-full bg-teal-light px-2 py-0.5 text-label-sm font-bold text-teal">
                  Full course
                </span>
              ) : (
                <span className="rounded-full bg-gray-100 px-2 py-0.5 text-label-sm font-bold text-gray-500">
                  Soon
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

export default function ChooseLanguageScreen() {
  return (
    <Suspense fallback={null}>
      <ChooseLanguageInner />
    </Suspense>
  );
}
