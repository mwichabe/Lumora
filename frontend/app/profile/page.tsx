"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  Flame,
  Zap,
  BookOpen,
  Globe,
  LogOut,
  ChevronRight,
  Plus,
  Check,
  GraduationCap,
  Award,
  Camera,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { Avatar } from "@/components/Avatar";
import { CropModal } from "@/components/CropModal";
import { ConfirmDialog } from "@/components/ConfirmDialog";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { languageMeta, languageName } from "@/lib/languages";
import { CharacterWithFriendship, Certificate, User } from "@/lib/types";

export default function ProfilePage() {
  const { user, logout, setUser } = useAuth();
  const [characters, setCharacters] = useState<CharacterWithFriendship[]>([]);
  const [lessonsDone, setLessonsDone] = useState(0);
  const [languages, setLanguages] = useState<string[]>([]);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [switching, setSwitching] = useState(false);
  const [confirmOut, setConfirmOut] = useState(false);

  useEffect(() => {
    api
      .characters()
      .then((r) => setCharacters(r.characters))
      .catch(() => setCharacters([]));
    api
      .enrollments()
      .then((r) => setLanguages(r.languages))
      .catch(() => setLanguages([]));
    api
      .certificates()
      .then((r) => setCertificates(r.certificates))
      .catch(() => setCertificates([]));
    // Lessons completed is derived from total XP as a friendly estimate.
    api
      .home()
      .then(() => {})
      .catch(() => {});
  }, []);

  async function switchTo(code: string) {
    if (code === user?.targetLanguage || switching) return;
    setSwitching(true);
    try {
      const r = await api.switchLanguage(code);
      setUser(r.user);
      setLanguages(r.languages);
    } catch {
      /* ignore */
    } finally {
      setSwitching(false);
    }
  }

  useEffect(() => {
    if (user) setLessonsDone(Math.max(0, Math.floor(user.xp / 10)));
  }, [user]);

  if (!user) return null;

  const fluencyPct = Math.min(100, Math.round((user.fluencyScore / 1000) * 100));

  return (
    <AppShell tabs wide>
      <div className="bg-cream pb-24 lg:pb-10">
        {/* Purple header */}
        <div className="relative rounded-b-[32px] bg-gradient-to-b from-purple to-purple-dark px-6 pb-8 pt-14 text-white">
          <div className="flex items-center gap-4">
            <AvatarUploader
              name={user.name}
              color={user.avatarColor}
              url={user.avatarUrl}
              onUploaded={setUser}
            />
            <div className="min-w-0">
              <h1 className="truncate text-heading-xl font-extrabold">
                {user.name}
              </h1>
              <p className="truncate text-body-sm text-purple-light">
                {user.email}
              </p>
              <span className="mt-1 inline-block rounded-full bg-amber px-3 py-0.5 text-label-md font-extrabold text-ink">
                {user.cefrLevel} · {user.levelName}
              </span>
            </div>
          </div>
        </div>

        {/* Fluency ring */}
        <div className="-mt-6 px-6">
          <div className="flex items-center gap-5 rounded-2xl bg-white p-5 shadow-card">
            <FluencyRing pct={fluencyPct} score={user.fluencyScore} />
            <div>
              <div className="text-heading-md font-extrabold text-ink">
                Fluency Score
              </div>
              <p className="text-body-sm text-slatey">
                {user.targetLanguage
                  ? `Your ${user.targetLanguage} mastery, out of 1000.`
                  : "Your mastery, out of 1000."}
              </p>
            </div>
          </div>
        </div>

        {/* Stats row */}
        <div className="mt-5 grid grid-cols-2 gap-3 px-6 lg:grid-cols-4">
          <StatTile icon={<Flame size={20} />} label="Day Streak" value={user.streak} tint="#FF5C5C" />
          <StatTile icon={<Zap size={20} />} label="Total XP" value={user.xp} tint="#F5A623" />
          <StatTile icon={<BookOpen size={20} />} label="Lessons" value={lessonsDone} tint="#6C3FC5" />
          <StatTile icon={<Globe size={20} />} label="Languages" value={languages.length || (user.targetLanguage ? 1 : 0)} tint="#00C2A8" />
        </div>

        {/* My Languages */}
        <div className="mt-7 px-6">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-heading-md font-extrabold text-ink">
              My Languages
            </h2>
            <Link
              href="/onboarding/language?add=1"
              className="flex items-center gap-1.5 rounded-full bg-purple-light px-3 py-1.5 text-label-lg font-bold text-purple transition hover:brightness-95"
            >
              <Plus size={16} /> Add
            </Link>
          </div>

          <div className="overflow-hidden rounded-2xl bg-white shadow-card">
            {languages.length === 0 && (
              <p className="px-5 py-4 text-body-sm text-slatey">
                No languages yet. Tap “Add” to start a course.
              </p>
            )}
            {languages.map((code, i) => {
              const m = languageMeta(code);
              const active = code === user.targetLanguage;
              return (
                <button
                  key={code}
                  onClick={() => switchTo(code)}
                  disabled={switching}
                  className={`flex w-full items-center gap-3 px-5 py-4 text-left transition hover:bg-gray-50 ${
                    i > 0 ? "border-t border-gray-100" : ""
                  } ${active ? "bg-purple-light/40" : ""}`}
                >
                  <span className="text-2xl">{m?.flag || "🌐"}</span>
                  <div className="flex-1">
                    <p className="font-extrabold text-ink">{m?.name || code}</p>
                    <p className="text-body-sm text-slatey">
                      {active ? "Active course" : "Tap to switch"}
                    </p>
                  </div>
                  {active ? (
                    <span className="flex items-center gap-1 rounded-full bg-purple px-2.5 py-0.5 text-label-md font-bold text-white">
                      <Check size={14} /> Active
                    </span>
                  ) : (
                    <ChevronRight size={18} className="text-gray-300" />
                  )}
                </button>
              );
            })}
          </div>
        </div>

        {/* Certificates */}
        <div className="mt-7 px-6">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-heading-md font-extrabold text-ink">
              Certificates
            </h2>
            <Link
              href="/exam"
              className="flex items-center gap-1.5 rounded-full bg-purple px-3 py-1.5 text-label-lg font-bold text-white transition hover:bg-purple-dark"
            >
              <GraduationCap size={16} /> Get certified
            </Link>
          </div>

          {certificates.length === 0 ? (
            <Link
              href="/exam"
              className="flex items-center gap-3 rounded-2xl border-2 border-dashed border-purple/30 bg-white p-4 transition hover:border-purple/60"
            >
              <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-purple-light text-purple">
                <Award size={24} />
              </span>
              <div>
                <p className="font-extrabold text-ink">Earn your first certificate</p>
                <p className="text-body-sm text-slatey">
                  Take the proficiency exam to certify your level.
                </p>
              </div>
            </Link>
          ) : (
            <div className="space-y-2">
              {certificates.map((cert) => (
                <Link
                  key={cert.id}
                  href={`/certificates/${cert.id}`}
                  className="flex items-center gap-3 rounded-2xl bg-white p-4 shadow-card transition hover:shadow-card-lg"
                >
                  <span className="flex h-12 w-12 shrink-0 flex-col items-center justify-center rounded-xl bg-purple text-white">
                    <span className="text-body-md font-extrabold leading-none">
                      {cert.level}
                    </span>
                  </span>
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-extrabold text-ink">
                      {languageName(cert.language)} · {cert.level}
                    </p>
                    <p className="text-body-sm text-slatey">
                      Score {cert.score}% ·{" "}
                      {new Date(cert.issuedAt).toLocaleDateString()}
                    </p>
                  </div>
                  <ChevronRight size={18} className="text-gray-300" />
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Characters */}
        <div className="mt-7 px-6">
          <h2 className="mb-3 text-heading-md font-extrabold text-ink">
            Your Companions
          </h2>
          <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 lg:grid-cols-6">
            {characters.map((c) => (
              <motion.div
                key={c.id}
                whileTap={{ scale: 0.95 }}
                className="flex flex-col items-center rounded-2xl bg-white p-3 shadow-card"
              >
                <div
                  className="flex h-14 w-14 items-center justify-center rounded-full text-3xl"
                  style={{ backgroundColor: (c.color || "#EDE7F6") + "33" }}
                >
                  {c.emoji}
                </div>
                <div className="mt-2 text-center text-label-md font-extrabold text-ink">
                  {c.name}
                </div>
                <div className="mt-1 flex gap-0.5">
                  {Array.from({ length: 3 }).map((_, i) => (
                    <span
                      key={i}
                      className={`h-1.5 w-1.5 rounded-full ${
                        i < c.friendshipLevel ? "bg-amber" : "bg-gray-300"
                      }`}
                    />
                  ))}
                </div>
              </motion.div>
            ))}
            {characters.length === 0 &&
              Array.from({ length: 6 }).map((_, i) => (
                <div
                  key={i}
                  className="h-28 animate-pulse rounded-2xl bg-gray-100"
                />
              ))}
          </div>
        </div>

        {/* Settings */}
        <div className="mt-7 px-6">
          <div className="overflow-hidden rounded-2xl bg-white shadow-card">
            <SettingsRow label="Account Settings" href="/profile/settings" />
            <SettingsRow label="Notifications" href="/notifications" />
            <SettingsRow label="Help & Support" href="/profile/help" />
            <button
              onClick={() => setConfirmOut(true)}
              className="flex w-full items-center justify-between px-5 py-4 text-coral hover:bg-coral-light"
            >
              <span className="flex items-center gap-3 font-extrabold">
                <LogOut size={18} /> Sign Out
              </span>
              <ChevronRight size={18} />
            </button>
          </div>
        </div>

        <ConfirmDialog
          open={confirmOut}
          title="Sign out?"
          message="You'll need to sign in again to continue learning."
          confirmLabel="Sign out"
          danger
          onConfirm={() => {
            setConfirmOut(false);
            logout();
          }}
          onCancel={() => setConfirmOut(false)}
        />
      </div>
    </AppShell>
  );
}

function AvatarUploader({
  name,
  color,
  url,
  onUploaded,
}: {
  name: string;
  color: string;
  url: string;
  onUploaded: (u: User) => void;
}) {
  const [busy, setBusy] = useState(false);
  const [pending, setPending] = useState<File | null>(null);

  function onPick(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    e.target.value = ""; // allow re-selecting the same file
    if (file) setPending(file); // open the cropper
  }

  async function uploadCropped(file: File) {
    setPending(null);
    setBusy(true);
    try {
      const r = await api.uploadAvatar(file);
      onUploaded(r.user);
    } catch {
      /* ignore — could surface a toast */
    } finally {
      setBusy(false);
    }
  }

  return (
    <>
      <label className="relative cursor-pointer" aria-label="Change profile photo">
        <span className="block rounded-full border-4 border-white/30">
          <Avatar name={name} color={color} url={url} size={80} />
        </span>
        <span className="absolute -bottom-0.5 -right-0.5 flex h-7 w-7 items-center justify-center rounded-full border-2 border-purple bg-white text-purple">
          {busy ? (
            <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-purple/40 border-t-purple" />
          ) : (
            <Camera size={14} />
          )}
        </span>
        <input
          type="file"
          accept="image/png,image/jpeg,image/webp,image/gif"
          onChange={onPick}
          className="hidden"
        />
      </label>

      {pending && (
        <CropModal
          file={pending}
          onCancel={() => setPending(null)}
          onCropped={uploadCropped}
        />
      )}
    </>
  );
}

function FluencyRing({ pct, score }: { pct: number; score: number }) {
  const r = 42;
  const c = 2 * Math.PI * r;
  const offset = c - (pct / 100) * c;
  return (
    <div className="relative h-24 w-24 shrink-0">
      <svg width={96} height={96} className="-rotate-90">
        <circle cx={48} cy={48} r={r} stroke="#EDE7F6" strokeWidth={8} fill="none" />
        <motion.circle
          cx={48}
          cy={48}
          r={r}
          stroke="#6C3FC5"
          strokeWidth={8}
          fill="none"
          strokeLinecap="round"
          strokeDasharray={c}
          initial={{ strokeDashoffset: c }}
          animate={{ strokeDashoffset: offset }}
          transition={{ duration: 1, ease: "easeOut" }}
        />
      </svg>
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <span className="text-heading-md font-extrabold text-purple">{score}</span>
        <span className="text-label-sm text-slatey">/ 1000</span>
      </div>
    </div>
  );
}

function StatTile({
  icon,
  label,
  value,
  tint,
}: {
  icon: React.ReactNode;
  label: string;
  value: number;
  tint: string;
}) {
  return (
    <div className="flex items-center gap-3 rounded-2xl bg-white p-4 shadow-card">
      <div
        className="flex h-10 w-10 items-center justify-center rounded-xl"
        style={{ backgroundColor: tint + "1A", color: tint }}
      >
        {icon}
      </div>
      <div>
        <div className="text-heading-md font-extrabold text-ink">{value}</div>
        <div className="text-label-md text-slatey">{label}</div>
      </div>
    </div>
  );
}

function SettingsRow({ label, href }: { label: string; href: string }) {
  return (
    <Link
      href={href}
      className="flex w-full items-center justify-between border-b border-gray-100 px-5 py-4 text-ink hover:bg-gray-50"
    >
      <span className="font-semibold">{label}</span>
      <ChevronRight size={18} className="text-gray-300" />
    </Link>
  );
}
