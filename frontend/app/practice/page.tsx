"use client";

import { motion } from "framer-motion";
import { MessageCircle, Headphones, Repeat } from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { SpeechBubble } from "@/components/widgets";

const MODES = [
  {
    title: "AI Conversation",
    desc: "Chat with Blaze to practise speaking in real scenarios.",
    icon: MessageCircle,
    tint: "#6C3FC5",
    soon: true,
  },
  {
    title: "Listening Drills",
    desc: "Sharpen your ear with Mira's audio challenges.",
    icon: Headphones,
    tint: "#00C2A8",
    soon: true,
  },
  {
    title: "Review Mistakes",
    desc: "Revisit tricky words Cora has saved for you.",
    icon: Repeat,
    tint: "#F5A623",
    soon: true,
  },
];

export default function PracticePage() {
  return (
    <AppShell tabs>
      <div className="bg-cream pb-24 lg:pb-10">
        <header className="px-6 pb-2 pt-14">
          <h1 className="text-display-lg font-extrabold text-ink">Practice</h1>
          <p className="mt-1 text-body-md text-slatey">
            Strengthen your skills between lessons.
          </p>
        </header>

        <div className="px-6 pt-2">
          <SpeechBubble className="max-w-none">
            <span className="font-extrabold text-purple">Blaze:</span> Ready to
            warm up? Pick a practice mode and let&apos;s get fluent! 🔥
          </SpeechBubble>
        </div>

        <div className="mt-5 space-y-3 px-6">
          {MODES.map((m, i) => {
            const Icon = m.icon;
            return (
              <motion.button
                key={m.title}
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.08 }}
                whileTap={{ scale: 0.98 }}
                className="flex w-full items-center gap-4 rounded-2xl bg-white p-4 text-left shadow-card"
              >
                <div
                  className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl"
                  style={{ backgroundColor: m.tint + "1A", color: m.tint }}
                >
                  <Icon size={24} />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="font-extrabold text-ink">{m.title}</span>
                    {m.soon && (
                      <span className="rounded-full bg-amber-light px-2 py-0.5 text-label-sm font-extrabold text-amber">
                        Soon
                      </span>
                    )}
                  </div>
                  <p className="mt-0.5 text-body-sm text-slatey">{m.desc}</p>
                </div>
              </motion.button>
            );
          })}
        </div>
      </div>
    </AppShell>
  );
}
