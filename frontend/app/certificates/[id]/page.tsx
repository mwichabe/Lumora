"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft, Download, ShieldCheck, Link2, Check } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { api } from "@/lib/api";
import { languageName, languageMeta } from "@/lib/languages";
import type { Certificate } from "@/lib/types";

export default function CertificatePage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [cert, setCert] = useState<Certificate | null>(null);
  const [missing, setMissing] = useState(false);
  const [copied, setCopied] = useState(false);
  const [verifyUrl, setVerifyUrl] = useState("");

  useEffect(() => {
    api
      .certificate(id)
      .then((d) => setCert(d.certificate))
      .catch(() => setMissing(true));
  }, [id]);

  useEffect(() => {
    if (cert?.serial && typeof window !== "undefined") {
      setVerifyUrl(`${window.location.origin}/verify/${cert.serial}`);
    }
  }, [cert?.serial]);

  function copyLink() {
    if (!verifyUrl) return;
    navigator.clipboard?.writeText(verifyUrl).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  }

  if (missing) {
    return (
      <div className="flex min-h-[100dvh] flex-col items-center justify-center gap-3 bg-cream px-6 text-center">
        <p className="text-heading-sm font-extrabold text-ink">
          Certificate not found
        </p>
        <Button variant="outline" onClick={() => router.push("/profile")}>
          Back to profile
        </Button>
      </div>
    );
  }

  if (!cert) {
    return (
      <div className="flex min-h-[100dvh] items-center justify-center bg-cream">
        <FoxMascot size={110} glow />
      </div>
    );
  }

  const date = new Date(cert.issuedAt).toLocaleDateString(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
  const flag = languageMeta(cert.language)?.flag || "🌐";

  return (
    <div className="flex min-h-[100dvh] w-full justify-center bg-[#eceaf3] px-4 py-6">
      <style>{`@media print {
        body * { visibility: hidden !important; }
        #cert, #cert * { visibility: visible !important; }
        #cert { position: absolute; left: 0; top: 0; margin: 0 !important; width: 100%; box-shadow: none !important; }
        .no-print { display: none !important; }
      }`}</style>

      <div className="w-full max-w-3xl">
        {/* toolbar */}
        <div className="no-print mb-4 flex items-center justify-between">
          <button
            onClick={() => router.push("/profile")}
            className="flex items-center gap-1.5 text-body-sm font-semibold text-slatey transition hover:text-purple"
          >
            <ArrowLeft size={16} /> Back
          </button>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={copyLink} className="h-10 px-4">
              <span className="flex items-center gap-2">
                {copied ? <Check size={16} /> : <Link2 size={16} />}
                {copied ? "Copied" : "Share link"}
              </span>
            </Button>
            <Button onClick={() => window.print()} className="h-10 px-5">
              <span className="flex items-center gap-2">
                <Download size={18} /> Download
              </span>
            </Button>
          </div>
        </div>

        {/* the certificate */}
        <div
          id="cert"
          className="relative overflow-hidden rounded-2xl bg-white p-2 shadow-card-lg"
        >
          {/* guilloché-style background + watermark */}
          <div
            className="pointer-events-none absolute inset-0 opacity-[0.06]"
            style={{
              backgroundImage:
                "repeating-radial-gradient(circle at 50% 50%, #6C3FC5 0, #6C3FC5 1px, transparent 1px, transparent 14px), repeating-linear-gradient(45deg, #6C3FC5 0, #6C3FC5 1px, transparent 1px, transparent 22px)",
            }}
          />
          <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
            <span className="select-none text-[7rem] font-extrabold tracking-tight text-purple/[0.04] sm:text-[10rem]">
              LUMORA
            </span>
          </div>

          <div className="relative rounded-xl border-[3px] border-purple/30 p-8 text-center sm:p-12">
            {/* corner flourishes */}
            <div className="pointer-events-none absolute left-3 top-3 h-10 w-10 rounded-tl-xl border-l-4 border-t-4 border-amber" />
            <div className="pointer-events-none absolute right-3 top-3 h-10 w-10 rounded-tr-xl border-r-4 border-t-4 border-amber" />
            <div className="pointer-events-none absolute bottom-3 left-3 h-10 w-10 rounded-bl-xl border-b-4 border-l-4 border-amber" />
            <div className="pointer-events-none absolute bottom-3 right-3 h-10 w-10 rounded-br-xl border-b-4 border-r-4 border-amber" />

            {/* verified badge */}
            <div className="absolute right-5 top-5 flex items-center gap-1 rounded-full bg-teal/10 px-2.5 py-1 text-label-sm font-bold text-teal">
              <ShieldCheck size={14} /> Certified by Lumora
            </div>

            <div className="flex items-center justify-center gap-2">
              <FoxMascot size={44} />
              <span className="text-heading-md font-extrabold tracking-tight text-purple">
                Lumora
              </span>
            </div>

            <p className="mt-6 text-label-lg font-bold uppercase tracking-[0.2em] text-gray-500">
              Certificate of Achievement
            </p>

            <p className="mt-6 text-body-md text-slatey">This certifies that</p>
            <h1 className="mt-1 text-display-lg font-extrabold text-ink">
              {cert.userName}
            </h1>

            <p className="mx-auto mt-4 max-w-md text-body-lg text-slatey">
              has successfully demonstrated{" "}
              <strong className="text-ink">
                {flag} {languageName(cert.language)}
              </strong>{" "}
              proficiency at level
            </p>

            {/* level seal */}
            <div className="relative mx-auto mt-5 flex h-28 w-28 flex-col items-center justify-center rounded-full bg-purple text-white shadow-float">
              <span className="absolute inset-1 rounded-full border-2 border-dashed border-white/40" />
              <ShieldCheck size={18} />
              <span className="text-display-lg font-extrabold leading-none">
                {cert.level}
              </span>
            </div>
            <p className="mt-3 text-body-md font-bold text-purple">
              Overall score: {cert.score}%
            </p>

            {/* section breakdown */}
            <div className="mx-auto mt-6 grid max-w-md grid-cols-4 gap-2 text-center">
              {[
                ["Listening", cert.listening],
                ["Reading", cert.reading],
                ["Writing", cert.writing],
                ["Speaking", cert.speaking],
              ].map(([name, val]) => (
                <div key={name as string} className="rounded-lg bg-gray-50 py-2">
                  <p className="text-heading-sm font-extrabold text-ink">
                    {val}%
                  </p>
                  <p className="text-label-sm text-slatey">{name}</p>
                </div>
              ))}
            </div>

            {/* footer */}
            <div className="mt-8 flex items-end justify-between">
              <div className="text-left">
                <p className="border-t-2 border-gray-200 pt-1 text-body-sm text-slatey">
                  {date}
                </p>
                <p className="text-label-sm text-gray-500">Date issued</p>
              </div>
              <div className="text-right">
                <p
                  className="pt-1 text-heading-md text-purple"
                  style={{ fontFamily: "cursive" }}
                >
                  Lumora
                </p>
                <p className="border-t-2 border-gray-200 pt-1 text-label-sm text-gray-500">
                  Lumora Academy
                </p>
              </div>
            </div>

            {/* verification strip */}
            <div className="mt-6 rounded-lg bg-gray-50 px-4 py-2 text-center">
              <p className="text-label-sm text-gray-500">
                Verification ID:{" "}
                <span className="font-bold tracking-wide text-ink">
                  {cert.serial || "—"}
                </span>
              </p>
              {verifyUrl && (
                <p className="mt-0.5 break-all text-label-sm text-slatey">
                  Verify at {verifyUrl}
                </p>
              )}
            </div>
          </div>
        </div>

        <p className="no-print mt-3 text-center text-body-sm text-slatey">
          Tip: in the print dialog choose “Save as PDF” to keep a copy. Anyone can
          confirm this certificate is genuine using the verification link.
        </p>
      </div>
    </div>
  );
}
