"use client";

import { useState } from "react";
import Link from "next/link";
import { ArrowLeft, Mail, MailCheck } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { api } from "@/lib/api";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [busy, setBusy] = useState(false);
  const [sent, setSent] = useState(false);

  async function submit() {
    if (busy || email.trim().length < 4) return;
    setBusy(true);
    try {
      await api.forgotPassword(email.trim());
    } catch {
      /* always show the same confirmation — don't reveal whether the email exists */
    } finally {
      setBusy(false);
      setSent(true);
    }
  }

  return (
    <div className="flex min-h-[100dvh] w-full items-center justify-center bg-[#eceaf3] px-5 py-10">
      <div className="w-full max-w-sm">
        <Link
          href="/onboarding/welcome"
          className="mb-4 inline-flex items-center gap-1.5 text-body-sm font-semibold text-slatey transition hover:text-purple"
        >
          <ArrowLeft size={16} /> Back to sign in
        </Link>

        <div className="rounded-3xl bg-white p-7 shadow-card-lg">
          {sent ? (
            <div className="text-center">
              <span className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-teal/10 text-teal">
                <MailCheck size={32} />
              </span>
              <h1 className="mt-4 text-heading-xl font-extrabold text-ink">
                Check your email
              </h1>
              <p className="mt-1 text-body-md text-slatey">
                If an account exists for <strong>{email.trim()}</strong>, we&apos;ve
                sent a link to reset your password. It expires in 1 hour.
              </p>
              <Link href="/onboarding/welcome" className="mt-6 block">
                <Button full>Back to sign in</Button>
              </Link>
              <button
                onClick={() => setSent(false)}
                className="mt-3 w-full text-body-sm font-semibold text-slatey transition hover:text-purple"
              >
                Use a different email
              </button>
            </div>
          ) : (
            <>
              <div className="text-center">
                <FoxMascot size={72} glow />
                <h1 className="mt-3 text-heading-xl font-extrabold text-ink">
                  Forgot your password?
                </h1>
                <p className="mt-1 text-body-md text-slatey">
                  Enter your email and we&apos;ll send you a link to reset it.
                </p>
              </div>

              <div className="mt-6">
                <span className="mb-1.5 block text-body-sm font-semibold text-slatey">
                  Email
                </span>
                <div className="relative flex items-center">
                  <span className="pointer-events-none absolute left-3.5 text-gray-500">
                    <Mail size={18} />
                  </span>
                  <input
                    type="email"
                    autoComplete="email"
                    autoFocus
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && submit()}
                    placeholder="you@email.com"
                    className="h-[52px] w-full rounded-xl border border-gray-100 bg-gray-50 pl-11 pr-4 text-body-lg outline-none transition focus:border-purple focus:bg-white focus:ring-4 focus:ring-purple/10"
                  />
                </div>
              </div>

              <div className="mt-6">
                <Button
                  full
                  loading={busy}
                  disabled={email.trim().length < 4}
                  onClick={submit}
                >
                  Send reset link
                </Button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
