"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import {
  Award,
  GraduationCap,
  Trash2,
  ChevronRight,
  ShieldCheck,
  ArrowLeft,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { FoxMascot } from "@/components/FoxMascot";
import { ConfirmDialog } from "@/components/ConfirmDialog";
import { api } from "@/lib/api";
import { languageName, languageMeta } from "@/lib/languages";
import type { Certificate } from "@/lib/types";

const LEVEL_ORDER = ["A1", "A2", "B1", "B2", "C1", "C2"];

export default function CertificatesPage() {
  return (
    <AppShell tabs>
      <CertificatesContent />
    </AppShell>
  );
}

function CertificatesContent() {
  const router = useRouter();
  const [certs, setCerts] = useState<Certificate[] | null>(null);
  const [filter, setFilter] = useState<string>("all");
  const [pendingDelete, setPendingDelete] = useState<Certificate | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    api
      .certificates()
      .then((r) => setCerts(r.certificates))
      .catch(() => setCerts([]));
  }, []);

  const languages = useMemo(() => {
    if (!certs) return [];
    return Array.from(new Set(certs.map((c) => c.language)));
  }, [certs]);

  const visible = useMemo(() => {
    if (!certs) return [];
    const list =
      filter === "all" ? certs : certs.filter((c) => c.language === filter);
    return [...list].sort(
      (a, b) =>
        a.language.localeCompare(b.language) ||
        LEVEL_ORDER.indexOf(a.level) - LEVEL_ORDER.indexOf(b.level)
    );
  }, [certs, filter]);

  async function confirmDelete() {
    if (!pendingDelete) return;
    setDeleting(true);
    try {
      await api.deleteCertificate(pendingDelete.id);
      setCerts((cur) => (cur || []).filter((c) => c.id !== pendingDelete.id));
    } catch {
      /* ignore */
    } finally {
      setDeleting(false);
      setPendingDelete(null);
    }
  }

  if (certs === null) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center bg-cream">
        <FoxMascot size={110} glow />
      </div>
    );
  }

  return (
    <div className="bg-cream pb-24 lg:pb-10">
      <header className="px-6 pb-2 pt-14 lg:px-8 lg:pt-8">
        <button
          onClick={() => router.push("/profile")}
          className="mb-3 flex items-center gap-1.5 text-body-sm font-semibold text-slatey transition hover:text-purple"
        >
          <ArrowLeft size={16} /> Back to profile
        </button>
        <div className="flex items-start justify-between gap-3">
          <div>
            <h1 className="text-display-lg font-extrabold text-ink">
              Certificates
            </h1>
            <p className="mt-1 text-body-md text-slatey">
              {certs.length === 0
                ? "Your achievements will appear here."
                : `${certs.length} earned across ${languages.length} language${
                    languages.length === 1 ? "" : "s"
                  }.`}
            </p>
          </div>
          <Link
            href="/exam"
            className="flex shrink-0 items-center gap-1.5 rounded-full bg-purple px-4 py-2 text-label-lg font-bold text-white transition hover:bg-purple-dark"
          >
            <GraduationCap size={16} /> Get certified
          </Link>
        </div>
      </header>

      {certs.length === 0 ? (
        <div className="px-6 lg:px-8">
          <Link
            href="/exam"
            className="flex items-center gap-3 rounded-2xl border-2 border-dashed border-purple/30 bg-white p-5 transition hover:border-purple/60"
          >
            <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-purple-light text-purple">
              <Award size={24} />
            </span>
            <div>
              <p className="font-extrabold text-ink">
                Earn your first certificate
              </p>
              <p className="text-body-sm text-slatey">
                Take the proficiency exam to certify your level.
              </p>
            </div>
          </Link>
        </div>
      ) : (
        <>
          {/* language filter chips (only when more than one language) */}
          {languages.length > 1 && (
            <div className="flex flex-wrap gap-2 px-6 lg:px-8">
              <Chip
                active={filter === "all"}
                onClick={() => setFilter("all")}
                label={`All (${certs.length})`}
              />
              {languages.map((l) => (
                <Chip
                  key={l}
                  active={filter === l}
                  onClick={() => setFilter(l)}
                  label={`${languageMeta(l)?.flag || "🌐"} ${languageName(l)}`}
                />
              ))}
            </div>
          )}

          <div className="mt-4 grid gap-3 px-6 sm:grid-cols-2 lg:grid-cols-3 lg:px-8">
            {visible.map((cert, i) => (
              <CertCard
                key={cert.id}
                cert={cert}
                index={i}
                onDelete={() => setPendingDelete(cert)}
              />
            ))}
          </div>
        </>
      )}

      <ConfirmDialog
        open={!!pendingDelete}
        title="Delete certificate?"
        message={
          pendingDelete
            ? `This removes your ${languageName(pendingDelete.language)} ${
                pendingDelete.level
              } certificate. You'll be able to retake that level.`
            : undefined
        }
        confirmLabel={deleting ? "Deleting…" : "Delete"}
        danger
        onConfirm={confirmDelete}
        onCancel={() => (deleting ? null : setPendingDelete(null))}
      />
    </div>
  );
}

function Chip({
  active,
  onClick,
  label,
}: {
  active: boolean;
  onClick: () => void;
  label: string;
}) {
  return (
    <button
      onClick={onClick}
      className={`rounded-full px-3.5 py-1.5 text-label-lg font-bold transition ${
        active
          ? "bg-purple text-white"
          : "bg-white text-slatey shadow-card hover:text-ink"
      }`}
    >
      {label}
    </button>
  );
}

function CertCard({
  cert,
  index,
  onDelete,
}: {
  cert: Certificate;
  index: number;
  onDelete: () => void;
}) {
  const flag = languageMeta(cert.language)?.flag || "🌐";
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: Math.min(index * 0.04, 0.3) }}
      className="relative flex flex-col rounded-2xl bg-white p-4 shadow-card transition hover:shadow-card-lg"
    >
      <button
        onClick={onDelete}
        aria-label="Delete certificate"
        className="absolute right-3 top-3 flex h-8 w-8 items-center justify-center rounded-full text-gray-300 transition hover:bg-coral/10 hover:text-coral"
      >
        <Trash2 size={16} />
      </button>

      <Link href={`/certificates/${cert.id}`} className="flex flex-1 flex-col">
        <div className="flex items-center gap-3">
          <span className="flex h-14 w-14 shrink-0 flex-col items-center justify-center rounded-xl bg-purple text-white">
            <ShieldCheck size={14} />
            <span className="text-body-md font-extrabold leading-none">
              {cert.level}
            </span>
          </span>
          <div className="min-w-0 flex-1 pr-6">
            <p className="truncate font-extrabold text-ink">
              {flag} {languageName(cert.language)}
            </p>
            <p className="text-body-sm text-slatey">Level {cert.level}</p>
          </div>
        </div>

        <div className="mt-4 flex items-end justify-between">
          <div>
            <p className="text-heading-xl font-extrabold text-purple">
              {cert.score}%
            </p>
            <p className="text-label-md text-slatey">
              {new Date(cert.issuedAt).toLocaleDateString()}
            </p>
          </div>
          <span className="flex items-center gap-1 text-label-md font-bold text-purple">
            View <ChevronRight size={14} />
          </span>
        </div>
      </Link>
    </motion.div>
  );
}
