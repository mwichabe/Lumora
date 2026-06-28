"use client";

import Link from "next/link";
import { ArrowLeft, Mail, ChevronDown } from "lucide-react";
import { useState } from "react";
import { AppShell } from "@/components/AppShell";

const FAQ = [
  {
    q: "How do I add another language?",
    a: "Open the Learn tab and tap the language switcher (top-right), or go to Profile → My Languages → Add. Your progress is saved per language.",
  },
  {
    q: "How are streaks calculated?",
    a: "Complete at least one lesson each day to keep your streak alive. Miss a day and it resets — so even a quick 5-minute lesson counts!",
  },
  {
    q: "Why can't I open a listening or reading session?",
    a: "Those unlock once you've completed all the lessons in that unit. Finish the unit's lessons first, then they open up.",
  },
  {
    q: "How does the fluency score work?",
    a: "In speaking practice we listen to what you say and compare it to the target phrase, giving a 0–100% score. It needs a Chrome or Edge browser.",
  },
  {
    q: "How do I change my daily goal or password?",
    a: "Go to Profile → Account Settings. You can edit your name, avatar, daily goal, change your password, or delete your account there.",
  },
];

export default function HelpPage() {
  return (
    <AppShell tabs>
      <HelpContent />
    </AppShell>
  );
}

function HelpContent() {
  const [open, setOpen] = useState<number | null>(0);

  return (
    <div className="min-h-full bg-cream pb-24 lg:pb-10">
      <header className="flex items-center gap-3 border-b border-gray-100 bg-white px-4 py-3.5 lg:rounded-t-3xl lg:px-6">
        <Link
          href="/profile"
          aria-label="Back"
          className="flex h-9 w-9 items-center justify-center rounded-full text-slatey transition hover:bg-gray-50"
        >
          <ArrowLeft size={20} />
        </Link>
        <h1 className="text-heading-lg font-extrabold text-ink">Help &amp; Support</h1>
      </header>

      <div className="space-y-6 px-5 py-5 lg:px-6">
        <section>
          <h2 className="mb-3 text-label-lg font-bold uppercase tracking-wide text-gray-500">
            Frequently asked
          </h2>
          <div className="overflow-hidden rounded-2xl bg-white shadow-card">
            {FAQ.map((item, i) => {
              const isOpen = open === i;
              return (
                <div key={i} className={i > 0 ? "border-t border-gray-100" : ""}>
                  <button
                    onClick={() => setOpen(isOpen ? null : i)}
                    className="flex w-full items-center justify-between gap-3 px-5 py-4 text-left"
                  >
                    <span className="font-extrabold text-ink">{item.q}</span>
                    <ChevronDown
                      size={18}
                      className={`shrink-0 text-gray-400 transition ${
                        isOpen ? "rotate-180" : ""
                      }`}
                    />
                  </button>
                  {isOpen && (
                    <p className="px-5 pb-4 text-body-md text-slatey">{item.a}</p>
                  )}
                </div>
              );
            })}
          </div>
        </section>

        <section>
          <h2 className="mb-3 text-label-lg font-bold uppercase tracking-wide text-gray-500">
            Still need help?
          </h2>
          <a
            href="mailto:support@lumora.app?subject=Lumora%20Support"
            className="flex items-center gap-3 rounded-2xl bg-white p-4 shadow-card transition hover:shadow-card-lg"
          >
            <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-purple-light text-purple">
              <Mail size={22} />
            </span>
            <div>
              <p className="font-extrabold text-ink">Contact support</p>
              <p className="text-body-sm text-slatey">support@lumora.app</p>
            </div>
          </a>
        </section>
      </div>
    </div>
  );
}
