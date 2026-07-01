"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Lock, Eye, EyeOff, CheckCircle2 } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { api } from "@/lib/api";

export default function ResetPasswordPage() {
  const router = useRouter();
  const [token, setToken] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [show, setShow] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState(false);

  useEffect(() => {
    const t = new URLSearchParams(window.location.search).get("token") || "";
    setToken(t);
  }, []);

  const canSubmit = password.length >= 6 && password === confirm && !!token;

  async function submit() {
    if (busy) return;
    setError(null);
    if (password.length < 6) {
      setError("Password must be at least 6 characters.");
      return;
    }
    if (password !== confirm) {
      setError("Passwords don't match.");
      return;
    }
    setBusy(true);
    try {
      await api.resetPassword(token, password);
      setDone(true);
    } catch (e) {
      setError(
        (e as { message?: string })?.message ||
          "Could not reset your password. The link may be invalid or expired."
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="flex min-h-[100dvh] w-full items-center justify-center bg-[#eceaf3] px-5 py-10">
      <div className="w-full max-w-sm rounded-3xl bg-white p-7 shadow-card-lg">
        {done ? (
          <div className="text-center">
            <span className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-teal/10 text-teal">
              <CheckCircle2 size={32} />
            </span>
            <h1 className="mt-4 text-heading-xl font-extrabold text-ink">
              Password updated
            </h1>
            <p className="mt-1 text-body-md text-slatey">
              You can now sign in with your new password.
            </p>
            <Button full className="mt-6" onClick={() => router.push("/onboarding/welcome")}>
              Go to sign in
            </Button>
          </div>
        ) : !token ? (
          <div className="text-center">
            <FoxMascot size={72} glow />
            <h1 className="mt-3 text-heading-xl font-extrabold text-ink">
              Invalid reset link
            </h1>
            <p className="mt-1 text-body-md text-slatey">
              This link is missing or malformed. Please request a new one.
            </p>
            <Link href="/forgot-password" className="mt-6 block">
              <Button full>Request a new link</Button>
            </Link>
          </div>
        ) : (
          <>
            <div className="text-center">
              <FoxMascot size={72} glow />
              <h1 className="mt-3 text-heading-xl font-extrabold text-ink">
                Set a new password
              </h1>
              <p className="mt-1 text-body-md text-slatey">
                Choose a strong password you&apos;ll remember.
              </p>
            </div>

            <div className="mt-6 space-y-3">
              <PwField
                label="New password"
                value={password}
                onChange={setPassword}
                show={show}
                onToggle={() => setShow((s) => !s)}
                onSubmit={submit}
                placeholder="6+ characters"
              />
              <PwField
                label="Confirm password"
                value={confirm}
                onChange={setConfirm}
                show={show}
                onToggle={() => setShow((s) => !s)}
                onSubmit={submit}
                placeholder="Re-enter password"
              />

              {error && (
                <p className="rounded-lg bg-coral-light px-3 py-2 text-body-sm font-semibold text-coral">
                  {error}
                </p>
              )}
            </div>

            <Button full className="mt-6" loading={busy} disabled={!canSubmit} onClick={submit}>
              Reset password
            </Button>
            <Link
              href="/onboarding/welcome"
              className="mt-4 block text-center text-body-sm font-semibold text-slatey transition hover:text-purple"
            >
              Back to sign in
            </Link>
          </>
        )}
      </div>
    </div>
  );
}

function PwField({
  label,
  value,
  onChange,
  show,
  onToggle,
  onSubmit,
  placeholder,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  show: boolean;
  onToggle: () => void;
  onSubmit: () => void;
  placeholder: string;
}) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-body-sm font-semibold text-slatey">
        {label}
      </span>
      <div className="relative flex items-center">
        <span className="pointer-events-none absolute left-3.5 text-gray-500">
          <Lock size={18} />
        </span>
        <input
          type={show ? "text" : "password"}
          autoComplete="new-password"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && onSubmit()}
          placeholder={placeholder}
          className="h-[52px] w-full rounded-xl border border-gray-100 bg-gray-50 pl-11 pr-11 text-body-lg outline-none transition focus:border-purple focus:bg-white focus:ring-4 focus:ring-purple/10"
        />
        <button
          type="button"
          onClick={onToggle}
          aria-label={show ? "Hide password" : "Show password"}
          className="absolute right-3.5 text-gray-500 transition hover:text-purple"
        >
          {show ? <EyeOff size={18} /> : <Eye size={18} />}
        </button>
      </div>
    </label>
  );
}
