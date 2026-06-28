"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { ChevronDown, Check, Plus } from "lucide-react";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { languageMeta, languageName } from "@/lib/languages";

/**
 * Lets a learner switch between the languages they're enrolled in, or add a new
 * one. Designed to sit in a coloured (purple) header.
 */
export function LanguageSwitcher({ onChanged }: { onChanged?: () => void }) {
  const { user, setUser } = useAuth();
  const [languages, setLanguages] = useState<string[]>([]);
  const [open, setOpen] = useState(false);
  const [busy, setBusy] = useState(false);

  const active = user?.targetLanguage || "";

  useEffect(() => {
    api
      .enrollments()
      .then((d) => setLanguages(d.languages))
      .catch(() => setLanguages(active ? [active] : []));
  }, [active]);

  async function pick(code: string) {
    if (code === active || busy) {
      setOpen(false);
      return;
    }
    setBusy(true);
    try {
      const r = await api.switchLanguage(code);
      setUser(r.user);
      setLanguages(r.languages);
      onChanged?.();
    } catch {
      /* ignore */
    } finally {
      setBusy(false);
      setOpen(false);
    }
  }

  const activeMeta = languageMeta(active);

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((o) => !o)}
        className="flex items-center gap-2 rounded-full bg-white/15 px-3 py-1.5 text-body-sm font-bold text-white transition hover:bg-white/25"
      >
        <span className="text-base">{activeMeta?.flag || "🌐"}</span>
        {languageName(active)}
        <ChevronDown size={16} className={open ? "rotate-180 transition" : "transition"} />
      </button>

      {open && (
        <>
          {/* click-away */}
          <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
          <div className="absolute right-0 z-20 mt-2 w-56 overflow-hidden rounded-2xl bg-white shadow-card-lg">
            <p className="px-4 pb-1 pt-3 text-label-md font-bold uppercase tracking-wide text-gray-500">
              My languages
            </p>
            {languages.map((code) => {
              const m = languageMeta(code);
              const isActive = code === active;
              return (
                <button
                  key={code}
                  onClick={() => pick(code)}
                  className={`flex w-full items-center gap-3 px-4 py-2.5 text-left transition hover:bg-gray-50 ${
                    isActive ? "bg-purple-light/50" : ""
                  }`}
                >
                  <span className="text-xl">{m?.flag || "🌐"}</span>
                  <span className="flex-1 text-body-md font-semibold text-ink">
                    {m?.name || code}
                  </span>
                  {isActive && <Check size={16} className="text-purple" />}
                </button>
              );
            })}
            <Link
              href="/onboarding/language?add=1"
              className="flex items-center gap-3 border-t border-gray-100 px-4 py-3 text-body-md font-bold text-purple transition hover:bg-purple-light/40"
            >
              <Plus size={18} /> Add a language
            </Link>
          </div>
        </>
      )}
    </div>
  );
}
