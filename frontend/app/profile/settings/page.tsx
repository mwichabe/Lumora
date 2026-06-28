"use client";

import { useState } from "react";
import Link from "next/link";
import {
  ArrowLeft,
  Eye,
  EyeOff,
  LogOut,
  Trash2,
  AlertTriangle,
  Camera,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { Button } from "@/components/Button";
import { Avatar } from "@/components/Avatar";
import { CropModal } from "@/components/CropModal";
import { ConfirmDialog } from "@/components/ConfirmDialog";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";

const GOALS = [
  { label: "Casual", xp: 10, mins: "5 min/day" },
  { label: "Regular", xp: 20, mins: "10 min/day" },
  { label: "Serious", xp: 30, mins: "15 min/day" },
  { label: "Intense", xp: 50, mins: "20+ min/day" },
];

export default function SettingsPage() {
  return (
    <AppShell tabs>
      <SettingsContent />
    </AppShell>
  );
}

function SettingsContent() {
  const { user, setUser, logout } = useAuth();
  const [confirmOut, setConfirmOut] = useState(false);

  if (!user) return null;

  return (
    <div className="min-h-full bg-cream pb-24 lg:pb-10">
      <header className="flex items-center gap-3 border-b border-gray-100 bg-white px-4 py-3.5 lg:rounded-t-3xl lg:px-6">
        <Link
          href="/profile"
          aria-label="Back"
          className="flex h-9 w-9 items-center justify-center rounded-full text-slatey transition hover:bg-gray-50"
        >
          <ArrowLeft size={20} />
        </Link>
        <h1 className="text-heading-lg font-extrabold text-ink">
          Account Settings
        </h1>
      </header>

      <div className="space-y-6 px-5 py-5 lg:px-6">
        <ProfileSection user={user} setUser={setUser} />
        <PasswordSection />
        <DangerSection logout={logout} />

        <button
          onClick={() => setConfirmOut(true)}
          className="flex w-full items-center justify-center gap-2 rounded-2xl bg-white p-4 font-extrabold text-slatey shadow-card transition hover:text-ink"
        >
          <LogOut size={18} /> Sign out
        </button>
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
  );
}

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="rounded-2xl bg-white p-5 shadow-card">
      <h2 className="mb-4 text-heading-sm font-extrabold text-ink">{title}</h2>
      {children}
    </section>
  );
}

function Notice({ kind, text }: { kind: "ok" | "err"; text: string }) {
  return (
    <p
      className={`mt-3 rounded-lg px-3 py-2 text-body-sm font-semibold ${
        kind === "ok" ? "bg-teal-light text-teal" : "bg-coral-light text-coral"
      }`}
    >
      {text}
    </p>
  );
}

function ProfileSection({
  user,
  setUser,
}: {
  user: NonNullable<ReturnType<typeof useAuth>["user"]>;
  setUser: ReturnType<typeof useAuth>["setUser"];
}) {
  const [name, setName] = useState(user.name || "");
  const [goal, setGoal] = useState(user.dailyGoalXp || 20);
  const [busy, setBusy] = useState(false);
  const [msg, setMsg] = useState<{ kind: "ok" | "err"; text: string } | null>(null);

  // Profile photo
  const [pending, setPending] = useState<File | null>(null);
  const [photoBusy, setPhotoBusy] = useState(false);

  const changed =
    name.trim() !== (user.name || "") || goal !== user.dailyGoalXp;

  async function save() {
    setBusy(true);
    setMsg(null);
    try {
      const r = await api.updateProfile({
        name: name.trim(),
        dailyGoalXp: goal,
      });
      setUser(r.user);
      setMsg({ kind: "ok", text: "Profile updated." });
    } catch (e: any) {
      setMsg({ kind: "err", text: e?.message || "Could not save." });
    } finally {
      setBusy(false);
    }
  }

  async function uploadCropped(file: File) {
    setPending(null);
    setPhotoBusy(true);
    try {
      const r = await api.uploadAvatar(file);
      setUser(r.user);
    } catch {
      /* ignore */
    } finally {
      setPhotoBusy(false);
    }
  }

  async function removePhoto() {
    setPhotoBusy(true);
    try {
      const r = await api.removeAvatar();
      setUser(r.user);
    } catch {
      /* ignore */
    } finally {
      setPhotoBusy(false);
    }
  }

  return (
    <Card title="Profile">
      {/* Profile photo */}
      <p className="mb-2 text-body-sm font-semibold text-slatey">Profile photo</p>
      <div className="flex items-center gap-4">
        <Avatar name={user.name} color={user.avatarColor} url={user.avatarUrl} size={64} />
        <div className="flex flex-wrap gap-2">
          <label className="flex cursor-pointer items-center gap-2 rounded-full bg-purple px-4 py-2 text-label-lg font-bold text-white transition hover:bg-purple-dark">
            <Camera size={16} />
            {photoBusy ? "Working…" : user.avatarUrl ? "Change photo" : "Upload photo"}
            <input
              type="file"
              accept="image/png,image/jpeg,image/webp,image/gif"
              className="hidden"
              disabled={photoBusy}
              onChange={(e) => {
                const f = e.target.files?.[0];
                e.target.value = "";
                if (f) setPending(f);
              }}
            />
          </label>
          {user.avatarUrl && (
            <button
              onClick={removePhoto}
              disabled={photoBusy}
              className="rounded-full border-2 border-gray-200 px-4 py-2 text-label-lg font-bold text-slatey transition hover:border-coral hover:text-coral disabled:opacity-50"
            >
              Remove
            </button>
          )}
        </div>
      </div>

      <label className="mt-5 block">
        <span className="mb-1.5 block text-body-sm font-semibold text-slatey">
          Display name
        </span>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="h-[52px] w-full rounded-xl border border-gray-100 bg-gray-50 px-4 text-body-lg outline-none transition focus:border-purple focus:bg-white"
        />
      </label>

      <p className="mb-1.5 mt-4 text-body-sm font-semibold text-slatey">
        Daily goal
      </p>
      <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
        {GOALS.map((g) => (
          <button
            key={g.xp}
            onClick={() => setGoal(g.xp)}
            className={`rounded-xl border-2 p-3 text-center transition ${
              goal === g.xp
                ? "border-purple bg-purple-light"
                : "border-gray-100 bg-white"
            }`}
          >
            <span className="block text-body-md font-extrabold text-ink">
              {g.label}
            </span>
            <span className="block text-label-md text-slatey">{g.mins}</span>
          </button>
        ))}
      </div>

      <div className="mt-5">
        <Button full loading={busy} disabled={!changed} onClick={save}>
          Save changes
        </Button>
        {msg && <Notice kind={msg.kind} text={msg.text} />}
      </div>

      {pending && (
        <CropModal
          file={pending}
          onCancel={() => setPending(null)}
          onCropped={uploadCropped}
        />
      )}
    </Card>
  );
}

function PasswordSection() {
  const [current, setCurrent] = useState("");
  const [next, setNext] = useState("");
  const [confirm, setConfirm] = useState("");
  const [show, setShow] = useState(false);
  const [busy, setBusy] = useState(false);
  const [msg, setMsg] = useState<{ kind: "ok" | "err"; text: string } | null>(null);

  const valid = current.length > 0 && next.length >= 6 && next === confirm;

  async function save() {
    if (!valid) {
      if (next !== confirm) setMsg({ kind: "err", text: "Passwords don't match." });
      else setMsg({ kind: "err", text: "New password must be 6+ characters." });
      return;
    }
    setBusy(true);
    setMsg(null);
    try {
      await api.changePassword(current, next);
      setMsg({ kind: "ok", text: "Password changed." });
      setCurrent("");
      setNext("");
      setConfirm("");
    } catch (e: any) {
      setMsg({ kind: "err", text: e?.message || "Could not change password." });
    } finally {
      setBusy(false);
    }
  }

  const inputCls =
    "h-[52px] w-full rounded-xl border border-gray-100 bg-gray-50 px-4 pr-11 text-body-lg outline-none transition focus:border-purple focus:bg-white";

  return (
    <Card title="Change password">
      <div className="space-y-3">
        <input
          type={show ? "text" : "password"}
          value={current}
          onChange={(e) => setCurrent(e.target.value)}
          placeholder="Current password"
          className={inputCls}
        />
        <div className="relative">
          <input
            type={show ? "text" : "password"}
            value={next}
            onChange={(e) => setNext(e.target.value)}
            placeholder="New password (6+ characters)"
            className={inputCls}
          />
          <button
            type="button"
            onClick={() => setShow((s) => !s)}
            aria-label={show ? "Hide" : "Show"}
            className="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-500"
          >
            {show ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        </div>
        <input
          type={show ? "text" : "password"}
          value={confirm}
          onChange={(e) => setConfirm(e.target.value)}
          placeholder="Confirm new password"
          className={inputCls}
        />
      </div>
      <div className="mt-5">
        <Button full loading={busy} onClick={save}>
          Update password
        </Button>
        {msg && <Notice kind={msg.kind} text={msg.text} />}
      </div>
    </Card>
  );
}

function DangerSection({ logout }: { logout: () => void }) {
  const [open, setOpen] = useState(false);
  const [password, setPassword] = useState("");
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState("");

  async function destroy() {
    if (!password) {
      setErr("Enter your password to confirm.");
      return;
    }
    setBusy(true);
    setErr("");
    try {
      await api.deleteAccount(password);
      logout(); // clears the session and returns to the start screen
    } catch (e: any) {
      setErr(e?.message || "Could not delete account.");
      setBusy(false);
    }
  }

  return (
    <section className="rounded-2xl border-2 border-coral/30 bg-coral-light/40 p-5">
      <div className="flex items-center gap-2">
        <AlertTriangle size={18} className="text-coral" />
        <h2 className="text-heading-sm font-extrabold text-coral">Danger zone</h2>
      </div>
      <p className="mt-1 text-body-sm text-slatey">
        Deleting your account is permanent and removes all your progress,
        streaks and data. This cannot be undone.
      </p>

      {!open ? (
        <button
          onClick={() => setOpen(true)}
          className="mt-4 flex items-center justify-center gap-2 rounded-full border-2 border-coral px-5 py-2.5 font-extrabold text-coral transition hover:bg-coral hover:text-white"
        >
          <Trash2 size={18} /> Delete account
        </button>
      ) : (
        <div className="mt-4 space-y-3">
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Enter your password to confirm"
            className="h-[52px] w-full rounded-xl border border-coral/40 bg-white px-4 text-body-lg outline-none focus:border-coral"
          />
          {err && <p className="text-body-sm font-semibold text-coral">{err}</p>}
          <div className="flex gap-3">
            <Button
              variant="outline"
              className="flex-1"
              onClick={() => {
                setOpen(false);
                setPassword("");
                setErr("");
              }}
            >
              Cancel
            </Button>
            <Button
              variant="danger"
              className="flex-1"
              loading={busy}
              onClick={destroy}
            >
              Permanently delete
            </Button>
          </div>
        </div>
      )}
    </section>
  );
}
