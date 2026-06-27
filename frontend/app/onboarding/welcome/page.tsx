"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { FoxMascot } from "@/components/FoxMascot";
import { SpeechBubble } from "@/components/widgets";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";

export default function WelcomeScreen() {
  const router = useRouter();
  const { login, register } = useAuth();

  const [mode, setMode] = useState<"intro" | "signup" | "signin">("intro");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  async function submit() {
    setError("");
    setBusy(true);
    try {
      if (mode === "signup") {
        await register(email, password, name || "Learner");
        router.push("/onboarding/language");
      } else {
        await login(email, password);
        router.push("/home");
      }
    } catch (e: any) {
      setError(e?.message || "Something went wrong");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="app-frame flex flex-col bg-cream">
      {/* Top illustration area */}
      <div className="relative flex h-[46%] flex-col items-center justify-end bg-purple pb-6">
        <SpeechBubble className="mb-3">Hi! I&apos;m Lumora 👋</SpeechBubble>
        <FoxMascot size={180} glow bounce />
      </div>

      {/* Bottom card */}
      <motion.div
        initial={{ y: 40, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ type: "spring", stiffness: 300, damping: 26 }}
        className="-mt-6 flex flex-1 flex-col rounded-t-[32px] bg-cream px-6 pt-7"
      >
        {mode === "intro" ? (
          <>
            <h1 className="text-heading-xl font-extrabold">Hi! I&apos;m Lumora.</h1>
            <p className="mt-2 text-body-lg text-slatey">
              I&apos;m going to help you learn a whole new language. Ready?
            </p>
            <div className="mt-auto pb-8 pt-8">
              <Button full onClick={() => setMode("signup")}>
                Let&apos;s go!
              </Button>
              <button
                onClick={() => setMode("signin")}
                className="mt-4 w-full text-center text-body-sm text-teal"
              >
                Already have an account? Sign in
              </button>
            </div>
          </>
        ) : (
          <>
            <h1 className="text-heading-xl font-extrabold">
              {mode === "signup" ? "Create your account" : "Welcome back"}
            </h1>
            <p className="mt-1 text-body-md text-slatey">
              {mode === "signup"
                ? "Save your progress and start your streak."
                : "Pick up right where you left off."}
            </p>

            <div className="mt-5 space-y-3">
              {mode === "signup" && (
                <Field label="Name" value={name} onChange={setName} placeholder="What should I call you?" />
              )}
              <Field label="Email" value={email} onChange={setEmail} placeholder="you@email.com" type="email" />
              <Field
                label="Password"
                value={password}
                onChange={setPassword}
                placeholder="6+ characters"
                type="password"
              />
              {error && <p className="text-body-sm text-coral">{error}</p>}
            </div>

            <div className="mt-auto pb-8 pt-8">
              <Button full loading={busy} onClick={submit}>
                {mode === "signup" ? "Create account" : "Sign in"}
              </Button>
              <button
                onClick={() => setMode(mode === "signup" ? "signin" : "signup")}
                className="mt-4 w-full text-center text-body-sm text-teal"
              >
                {mode === "signup"
                  ? "Already have an account? Sign in"
                  : "New here? Create an account"}
              </button>
            </div>
          </>
        )}
      </motion.div>
    </div>
  );
}

function Field({
  label,
  value,
  onChange,
  placeholder,
  type = "text",
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  placeholder?: string;
  type?: string;
}) {
  return (
    <label className="block">
      <span className="mb-1 block text-body-sm font-semibold text-slatey">{label}</span>
      <input
        type={type}
        value={value}
        placeholder={placeholder}
        onChange={(e) => onChange(e.target.value)}
        className="h-[52px] w-full rounded-lg border border-gray-100 bg-gray-50 px-4 text-body-lg outline-none transition focus:border-purple focus:bg-white"
      />
    </label>
  );
}
