"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  ListChecks,
  Headphones,
  Mic,
  RefreshCw,
  ChevronRight,
  Sparkles,
  BookMarked,
  MessagesSquare,
  BookOpenText,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { SpeechBubble } from "@/components/widgets";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { languageName } from "@/lib/languages";

export default function PracticePage() {
  return (
    <AppShell tabs>
      <PracticeContent />
    </AppShell>
  );
}

function PracticeContent() {
  const { user } = useAuth();
  const [vocabCount, setVocabCount] = useState<number | null>(null);
  const [mistakeCount, setMistakeCount] = useState(0);
  const [listeningCount, setListeningCount] = useState(0);
  const [readingCount, setReadingCount] = useState(0);

  useEffect(() => {
    api
      .practice()
      .then((d) => {
        setVocabCount(d.vocab.length);
        setMistakeCount(d.mistakes.length);
        setListeningCount(d.listeningCount || 0);
        setReadingCount(d.readingCount || 0);
      })
      .catch(() => setVocabCount(0));
  }, [user?.targetLanguage]);

  const hasVocab = (vocabCount ?? 0) > 0;

  const modes = [
    {
      key: "quiz",
      title: "Vocabulary Quiz",
      desc: "Fresh words every time — choose the meaning.",
      icon: ListChecks,
      tint: "#6C3FC5",
      disabled: !hasVocab,
    },
    {
      key: "listening",
      title: "Listening Comprehension",
      desc:
        listeningCount > 0
          ? `${listeningCount} unlocked conversations · listen & answer`
          : "Unlock units to practise listening.",
      icon: MessagesSquare,
      tint: "#00C2A8",
      disabled: listeningCount === 0,
    },
    {
      key: "reading",
      title: "Reading Comprehension",
      desc:
        readingCount > 0
          ? `${readingCount} unlocked passages · read & answer`
          : "Unlock units to practise reading.",
      icon: BookOpenText,
      tint: "#17A3DD",
      disabled: readingCount === 0,
    },
    {
      key: "listen",
      title: "Listening Drill",
      desc: "Hear a word, choose what it means.",
      icon: Headphones,
      tint: "#0E9F8A",
      disabled: !hasVocab,
    },
    {
      key: "speak",
      title: "Speaking Practice",
      desc: "Say phrases aloud and get a fluency score.",
      icon: Mic,
      tint: "#FF5C5C",
      disabled: !hasVocab,
    },
    {
      key: "mistakes",
      title: "Review Mistakes",
      desc:
        mistakeCount > 0
          ? `${mistakeCount} to review`
          : "Anything you miss appears here.",
      icon: RefreshCw,
      tint: "#F5A623",
      disabled: mistakeCount === 0,
      badge: mistakeCount > 0 ? mistakeCount : undefined,
    },
  ];

  return (
    <div className="bg-cream pb-24 lg:pb-10">
      <header className="px-6 pb-2 pt-14 lg:px-8 lg:pt-8">
        <h1 className="text-display-lg font-extrabold text-ink">Practice</h1>
        <p className="mt-1 text-body-md text-slatey">
          Strengthen your {languageName(user?.targetLanguage)} between lessons.
        </p>
      </header>

      <div className="px-6 pt-2 lg:px-8">
        <SpeechBubble className="max-w-none">
          <span className="font-extrabold text-coral">Blaze:</span> Ready to warm
          up? Pick a drill and let&apos;s get fluent! 🔥
        </SpeechBubble>
      </div>

      {/* Stats */}
      <div className="mt-5 grid grid-cols-2 gap-3 px-6 lg:px-8">
        <div className="rounded-2xl bg-white p-4 shadow-card">
          <div className="flex items-center gap-2 text-purple">
            <BookMarked size={18} />
            <span className="text-label-md font-bold uppercase tracking-wide">
              Words ready
            </span>
          </div>
          <p className="mt-1 text-heading-xl font-extrabold text-ink">
            {vocabCount === null ? "—" : vocabCount}
          </p>
          <p className="text-body-sm text-slatey">grows as you unlock units</p>
        </div>
        <div className="rounded-2xl bg-white p-4 shadow-card">
          <div className="flex items-center gap-2 text-amber">
            <RefreshCw size={18} />
            <span className="text-label-md font-bold uppercase tracking-wide">
              To review
            </span>
          </div>
          <p className="mt-1 text-heading-xl font-extrabold text-ink">
            {mistakeCount}
          </p>
          <p className="text-body-sm text-slatey">
            {mistakeCount > 0 ? "mistakes to fix" : "all caught up"}
          </p>
        </div>
      </div>

      {/* Daily Mix hero */}
      {hasVocab && (
        <div className="mt-3 px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <Link href="/practice/run?mode=mix">
              <div className="flex items-center gap-4 rounded-2xl bg-purple p-5 text-left text-white shadow-card-lg transition hover:brightness-105">
                <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-white/20">
                  <Sparkles size={28} />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="text-heading-md font-extrabold">
                      Daily Mix
                    </span>
                    <span className="rounded-full bg-white/25 px-2 py-0.5 text-label-sm font-extrabold">
                      Recommended
                    </span>
                  </div>
                  <p className="mt-0.5 text-body-sm text-white/85">
                    A fresh blend of new words — quiz, listening &amp; speaking.
                  </p>
                </div>
                <ChevronRight size={22} className="text-white/70" />
              </div>
            </Link>
          </motion.div>
        </div>
      )}

      <div className="mt-3 grid gap-3 px-6 lg:grid-cols-2 lg:px-8">
        {modes.map((m, i) => {
          const Icon = m.icon;
          const card = (
            <div
              className={`flex w-full items-center gap-4 rounded-2xl bg-white p-4 text-left shadow-card transition ${
                m.disabled ? "opacity-60" : "hover:shadow-card-lg"
              }`}
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
                  {m.badge && (
                    <span className="rounded-full bg-amber px-2 py-0.5 text-label-sm font-extrabold text-ink">
                      {m.badge}
                    </span>
                  )}
                </div>
                <p className="mt-0.5 text-body-sm text-slatey">{m.desc}</p>
              </div>
              {!m.disabled && <ChevronRight size={20} className="text-gray-300" />}
            </div>
          );

          if (m.disabled) {
            return (
              <div key={m.key}>{card}</div>
            );
          }
          return (
            <motion.div
              key={m.key}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.06 }}
            >
              <Link href={`/practice/run?mode=${m.key}`}>{card}</Link>
            </motion.div>
          );
        })}
      </div>

      {vocabCount === 0 && (
        <p className="mt-5 px-6 text-center text-body-sm text-slatey lg:px-8">
          Complete a lesson first to unlock practice drills.
        </p>
      )}
    </div>
  );
}
