"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { ShieldCheck, ShieldX } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { api } from "@/lib/api";
import { languageName, languageMeta } from "@/lib/languages";
import type { CertVerification } from "@/lib/types";

/**
 * Public, no-auth page. Anyone with a certificate's serial can confirm it is
 * genuine and see the holder's result — but no private account data.
 */
export default function VerifyPage() {
  const { serial } = useParams<{ serial: string }>();
  const [data, setData] = useState<CertVerification | null>(null);

  useEffect(() => {
    api.verifyCertificate(serial).then(setData);
  }, [serial]);

  if (!data) {
    return (
      <div className="flex min-h-[100dvh] items-center justify-center bg-cream">
        <FoxMascot size={110} glow />
      </div>
    );
  }

  const cert = data.certificate;
  const flag = cert ? languageMeta(cert.language)?.flag || "🌐" : "🌐";
  const date = cert
    ? new Date(cert.issuedAt).toLocaleDateString(undefined, {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : "";

  return (
    <div className="flex min-h-[100dvh] w-full flex-col items-center justify-center bg-[#eceaf3] px-4 py-10">
      <div className="mb-5 flex items-center gap-2">
        <FoxMascot size={40} />
        <span className="text-heading-md font-extrabold tracking-tight text-purple">
          Lumora
        </span>
      </div>

      <div className="w-full max-w-md rounded-2xl bg-white p-7 text-center shadow-card-lg">
        {data.valid && cert ? (
          <>
            <span className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-teal/10 text-teal">
              <ShieldCheck size={34} />
            </span>
            <h1 className="mt-4 text-heading-lg font-extrabold text-ink">
              Verified Certificate
            </h1>
            <p className="mt-1 text-body-sm text-slatey">
              This is a genuine certificate issued by Lumora.
            </p>

            <div className="mt-6 space-y-3 text-left">
              <Row label="Holder" value={cert.userName} />
              <Row
                label="Language"
                value={`${flag} ${languageName(cert.language)}`}
              />
              <Row label="Level" value={cert.level} />
              <Row label="Overall score" value={`${cert.score}%`} />
              <Row label="Issued" value={date} />
              <Row label="Verification ID" value={cert.serial} mono />
            </div>

            <p className="mt-6 flex items-center justify-center gap-1.5 rounded-full bg-teal/10 py-2 text-label-md font-bold text-teal">
              <ShieldCheck size={15} /> Certified by Lumora
            </p>
          </>
        ) : (
          <>
            <span className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-coral/10 text-coral">
              <ShieldX size={34} />
            </span>
            <h1 className="mt-4 text-heading-lg font-extrabold text-ink">
              Not verified
            </h1>
            <p className="mt-1 text-body-sm text-slatey">
              We couldn&apos;t find a certificate with this ID. It may be mistyped
              or not genuine.
            </p>
            <p className="mt-4 rounded-lg bg-gray-50 px-3 py-2 font-mono text-label-sm text-slatey">
              {serial}
            </p>
          </>
        )}
      </div>
    </div>
  );
}

function Row({
  label,
  value,
  mono,
}: {
  label: string;
  value: string;
  mono?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-3 border-b border-gray-100 pb-2">
      <span className="text-label-md text-slatey">{label}</span>
      <span
        className={`text-body-md font-extrabold text-ink ${
          mono ? "font-mono tracking-wide" : ""
        }`}
      >
        {value}
      </span>
    </div>
  );
}
