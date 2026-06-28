"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { AnimatePresence, motion } from "framer-motion";
import {
  ArrowLeft,
  ArrowRight,
  Mail,
  Lock,
  User,
  Eye,
  EyeOff,
  Globe,
  Gamepad2,
  Headphones,
  Trophy,
  ShieldCheck,
} from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { useAuth } from "@/lib/auth";

type Mode = "intro" | "signup" | "signin";

const FEATURES = [
  { icon: Gamepad2, text: "Learn through bite-sized, game-like lessons" },
  { icon: Headphones, text: "Listen & speak with distinct character voices" },
  { icon: Globe, text: "40+ languages, from Spanish to Swahili" },
  { icon: Trophy, text: "Streaks, XP and weekly leagues to keep you going" },
];

export default function WelcomeScreen() {
  const router = useRouter();
  const { login, register } = useAuth();

  const [mode, setMode] = useState<Mode>("intro");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPw, setShowPw] = useState(false);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const canSubmit =
    email.trim().length > 3 &&
    password.length >= 6 &&
    (mode !== "signup" || name.trim().length > 0);

  async function submit() {
    if (!canSubmit || busy) return;
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

  function goTo(next: Mode) {
    setError("");
    setMode(next);
  }

  return (
    <div className="flex min-h-[100dvh] w-full flex-col bg-white lg:flex-row">
      {/* ── Brand panel ───────────────────────────────────────── */}
      <aside
        className="relative flex flex-col justify-between overflow-hidden px-6 pb-8 pt-12 text-white lg:w-[46%] lg:px-12 lg:pt-16"
        style={{
          background:
            "radial-gradient(120% 120% at 20% 0%, #7B4AD6 0%, #5E33B0 45%, #3A1F8A 100%)",
        }}
      >
        <div className="pointer-events-none absolute -right-16 top-10 h-64 w-64 rounded-full bg-amber/20 blur-3xl" />
        <div className="pointer-events-none absolute -left-16 bottom-0 h-64 w-64 rounded-full bg-teal/20 blur-3xl" />

        {/* Wordmark */}
        <div className="relative flex items-center gap-2">
          <FoxMascot size={40} />
          <span className="text-heading-lg font-extrabold tracking-tight">
            Lumora
          </span>
        </div>

        {/* Hero */}
        <div className="relative my-8 lg:my-0">
          <div className="flex justify-center lg:justify-start">
            <motion.div
              animate={{ y: [0, -10, 0] }}
              transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
            >
              <FoxMascot size={150} glow bounce />
            </motion.div>
          </div>
          <h1
            className="mt-6 text-center font-extrabold leading-tight tracking-tight lg:text-left"
            style={{ fontSize: "clamp(1.75rem, 3vw, 2.75rem)" }}
          >
            Learn a language.
            <br />
            <span className="text-amber">Fall in love with it.</span>
          </h1>

          {/* Feature list — desktop only (room) */}
          <ul className="mt-7 hidden space-y-3 lg:block">
            {FEATURES.map((f) => {
              const Icon = f.icon;
              return (
                <li key={f.text} className="flex items-center gap-3">
                  <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white/15">
                    <Icon size={18} />
                  </span>
                  <span className="text-body-md text-white/90">{f.text}</span>
                </li>
              );
            })}
          </ul>
        </div>

      </aside>

      {/* ── Form panel ────────────────────────────────────────── */}
      <main className="flex flex-1 flex-col items-center justify-center px-6 py-10 lg:px-12">
        <div className="w-full max-w-md">
          {/* Back control inside auth flows */}
          <AnimatePresence>
            {mode !== "intro" && (
              <motion.button
                initial={{ opacity: 0, x: -6 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -6 }}
                onClick={() => goTo("intro")}
                className="mb-6 flex items-center gap-1.5 text-body-sm font-semibold text-slatey transition hover:text-purple"
              >
                <ArrowLeft size={16} /> Back
              </motion.button>
            )}
          </AnimatePresence>

          <AnimatePresence mode="wait">
            {mode === "intro" ? (
              <motion.div
                key="intro"
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -12 }}
                transition={{ duration: 0.25 }}
              >
                <h2 className="text-display-lg font-extrabold text-ink">
                  Start learning today
                </h2>
                <p className="mt-2 text-body-lg text-slatey">
                  Create a free account and your first lesson is ready in
                  seconds.
                </p>

                <div className="mt-8 space-y-3">
                  <Button full onClick={() => goTo("signup")}>
                    <span className="flex items-center gap-2">
                      Create free account <ArrowRight size={18} />
                    </span>
                  </Button>
                  <Button variant="outline" full onClick={() => goTo("signin")}>
                    I already have an account
                  </Button>
                </div>

                <p className="mt-6 flex items-center justify-center gap-1.5 text-body-sm text-gray-500">
                  <ShieldCheck size={15} className="text-teal" />
                  Free forever · No credit card required
                </p>
              </motion.div>
            ) : (
              <motion.div
                key="auth"
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -12 }}
                transition={{ duration: 0.25 }}
              >
                <h2 className="text-display-lg font-extrabold text-ink">
                  {mode === "signup" ? "Create your account" : "Welcome back"}
                </h2>
                <p className="mt-2 text-body-md text-slatey">
                  {mode === "signup"
                    ? "Save your progress and start your streak."
                    : "Pick up right where you left off."}
                </p>

                <div className="mt-7 space-y-3">
                  {mode === "signup" && (
                    <Field
                      label="Name"
                      value={name}
                      onChange={setName}
                      onSubmit={submit}
                      placeholder="What should I call you?"
                      icon={<User size={18} />}
                      autoFocus
                    />
                  )}
                  <Field
                    label="Email"
                    value={email}
                    onChange={setEmail}
                    onSubmit={submit}
                    placeholder="you@email.com"
                    type="email"
                    icon={<Mail size={18} />}
                    autoFocus={mode === "signin"}
                  />
                  <Field
                    label="Password"
                    value={password}
                    onChange={setPassword}
                    onSubmit={submit}
                    placeholder="6+ characters"
                    type={showPw ? "text" : "password"}
                    icon={<Lock size={18} />}
                    trailing={
                      <button
                        type="button"
                        onClick={() => setShowPw((s) => !s)}
                        aria-label={showPw ? "Hide password" : "Show password"}
                        className="text-gray-500 transition hover:text-purple"
                      >
                        {showPw ? <EyeOff size={18} /> : <Eye size={18} />}
                      </button>
                    }
                  />

                  <AnimatePresence>
                    {error && (
                      <motion.p
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: "auto" }}
                        exit={{ opacity: 0, height: 0 }}
                        className="rounded-lg bg-coral-light px-3 py-2 text-body-sm font-semibold text-coral"
                      >
                        {error}
                      </motion.p>
                    )}
                  </AnimatePresence>
                </div>

                <div className="mt-7">
                  <Button full loading={busy} disabled={!canSubmit} onClick={submit}>
                    {mode === "signup" ? "Create account" : "Sign in"}
                  </Button>
                  <button
                    onClick={() => goTo(mode === "signup" ? "signin" : "signup")}
                    className="mt-4 w-full text-center text-body-md text-slatey transition hover:text-purple"
                  >
                    {mode === "signup" ? (
                      <>
                        Already have an account?{" "}
                        <span className="font-bold text-teal">Sign in</span>
                      </>
                    ) : (
                      <>
                        New here?{" "}
                        <span className="font-bold text-teal">Create one</span>
                      </>
                    )}
                  </button>
                </div>

                {mode === "signup" && (
                  <p className="mt-6 text-center text-body-sm text-gray-500">
                    By continuing you agree to our Terms & Privacy Policy.
                  </p>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </main>
    </div>
  );
}

function Stat({ value, label }: { value: string; label: string }) {
  return (
    <span className="flex items-baseline gap-1.5">
      <strong className="text-body-lg font-extrabold text-white">{value}</strong>
      <span className="text-white/60">{label}</span>
    </span>
  );
}

function Field({
  label,
  value,
  onChange,
  onSubmit,
  placeholder,
  type = "text",
  icon,
  trailing,
  autoFocus,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  onSubmit?: () => void;
  placeholder?: string;
  type?: string;
  icon?: React.ReactNode;
  trailing?: React.ReactNode;
  autoFocus?: boolean;
}) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-body-sm font-semibold text-slatey">
        {label}
      </span>
      <div className="relative flex items-center">
        {icon && (
          <span className="pointer-events-none absolute left-3.5 text-gray-500">
            {icon}
          </span>
        )}
        <input
          type={type}
          value={value}
          placeholder={placeholder}
          autoFocus={autoFocus}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") onSubmit?.();
          }}
          className={`h-[52px] w-full rounded-xl border border-gray-100 bg-gray-50 text-body-lg outline-none transition focus:border-purple focus:bg-white focus:ring-4 focus:ring-purple/10 ${
            icon ? "pl-11" : "pl-4"
          } ${trailing ? "pr-11" : "pr-4"}`}
        />
        {trailing && <span className="absolute right-3.5">{trailing}</span>}
      </div>
    </label>
  );
}
