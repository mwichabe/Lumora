"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Check } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { SpeechBubble } from "@/components/widgets";
import { Button } from "@/components/Button";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";

const GOALS = [
  { id: "casual", minutes: "5 min/day", xp: 10, note: "Perfect for busy days!" },
  { id: "regular", minutes: "10 min/day", xp: 20, note: "A solid habit starts here." },
  { id: "serious", minutes: "15 min/day", xp: 30, note: "You mean business. I love it." },
  { id: "intense", minutes: "20+ min/day", xp: 50, note: "Are you… trying to make me cry? (happy tears)" },
];

export default function SetGoalScreen() {
  const router = useRouter();
  const { setUser } = useAuth();
  const [selected, setSelected] = useState<number | null>(1);
  const [busy, setBusy] = useState(false);

  async function finish() {
    if (selected === null) return;
    setBusy(true);
    const target =
      (typeof window !== "undefined" && sessionStorage.getItem("lumora_target")) || "es";
    try {
      const { user } = await api.setup(target, GOALS[selected].xp, "");
      setUser(user);
      router.push("/home");
    } catch {
      // If the API is unreachable, still continue to the home screen.
      router.push("/home");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="app-frame flex flex-col bg-cream px-5 pt-10">
      <div className="flex items-start gap-2">
        <FoxMascot size={64} />
        <SpeechBubble>How much time can you give me each day?</SpeechBubble>
      </div>

      <div className="mt-6 space-y-3">
        {GOALS.map((g, i) => {
          const active = selected === i;
          return (
            <button
              key={g.id}
              onClick={() => setSelected(i)}
              className={`flex w-full items-center gap-3 rounded-xl border-2 px-4 py-3 text-left transition ${
                active ? "border-purple bg-purple-light" : "border-gray-100 bg-white"
              }`}
            >
              <span className="flex h-10 w-10 items-center justify-center rounded-full bg-amber/15 text-xl">
                🔥
              </span>
              <span className="flex-1">
                <span className="flex items-center justify-between">
                  <span className="text-body-lg font-bold capitalize">{g.id}</span>
                  <span className="text-body-sm font-bold text-amber">{g.xp} XP</span>
                </span>
                <span className="block text-body-sm text-slatey">{g.minutes}</span>
                <span className="block text-body-sm italic text-purple">{g.note}</span>
              </span>
              {active && (
                <span className="flex h-6 w-6 items-center justify-center rounded-full bg-purple text-white">
                  <Check size={16} />
                </span>
              )}
            </button>
          );
        })}
      </div>

      <div className="mt-auto pb-8 pt-6">
        <Button full disabled={selected === null} loading={busy} onClick={finish}>
          Set My Goal
        </Button>
      </div>
    </div>
  );
}
