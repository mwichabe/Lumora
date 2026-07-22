"use client";

import { useEffect, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import {
  Archive,
  ArrowRight,
  CheckSquare,
  Clock,
  GitMerge,
  History,
  Link2,
  ListTodo,
  Pencil,
  RotateCcw,
  Square,
  Tag,
  Trash2,
  Users,
  X,
} from "lucide-react";
import { Avatar } from "@/components/Avatar";
import { Button } from "@/components/Button";
import { ConfirmDialog } from "@/components/ConfirmDialog";
import { STATUS_META, TagChip, VoteBar, timeAgo } from "./IdeaBits";
import { api } from "@/lib/api";
import type { IdeaDetail, IdeaStatus, IdeaTask } from "@/lib/types";

/**
 * The right panel: everything about the selected idea that isn't conversation —
 * where it stands, how it got there, and what to do with it next.
 */

export function IdeaDetailsPanel({
  detail,
  loading,
  onRefresh,
  onOpenIdea,
  onDeleted,
  onClose,
}: {
  detail: IdeaDetail | null;
  loading: boolean;
  onRefresh: () => void;
  onOpenIdea: (id: number) => void;
  onDeleted: () => void;
  onClose?: () => void;
}) {
  const [tab, setTab] = useState<"about" | "history">("about");
  const [editing, setEditing] = useState(false);

  useEffect(() => {
    setTab("about");
    setEditing(false);
  }, [detail?.idea.id]);

  if (loading || !detail) {
    return (
      <div className="h-full border-l border-gray-100 bg-white p-4">
        {loading ? (
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-16 animate-pulse rounded-2xl bg-gray-50" />
            ))}
          </div>
        ) : (
          <p className="pt-12 text-center text-body-sm text-slatey">
            Select an idea to see its details.
          </p>
        )}
      </div>
    );
  }

  const { idea } = detail;

  return (
    <div className="flex h-full flex-col border-l border-gray-100 bg-white">
      <div className="shrink-0 border-b border-gray-100 px-4 py-3">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <p className="text-label-md font-bold text-slatey">Idea #{idea.id}</p>
            <h3 className="mt-0.5 font-extrabold leading-snug text-ink">
              {idea.title}
            </h3>
          </div>
          <div className="flex shrink-0 gap-1">
            {detail.canEdit && (
              <button
                onClick={() => setEditing(true)}
                aria-label="Edit idea"
                className="rounded-md p-1.5 text-slatey transition hover:bg-gray-50 hover:text-ink"
              >
                <Pencil size={15} />
              </button>
            )}
            {onClose && (
              <button
                onClick={onClose}
                aria-label="Close details"
                className="rounded-md p-1.5 text-slatey xl:hidden"
              >
                <X size={16} />
              </button>
            )}
          </div>
        </div>

        <div className="mt-2 flex gap-2 rounded-full bg-gray-50 p-1">
          <TabBtn active={tab === "about"} onClick={() => setTab("about")}>
            About
          </TabBtn>
          <TabBtn active={tab === "history"} onClick={() => setTab("history")}>
            History
          </TabBtn>
        </div>
      </div>

      <div className="min-h-0 flex-1 space-y-5 overflow-y-auto p-4">
        {tab === "about" ? (
          <AboutTab
            detail={detail}
            onRefresh={onRefresh}
            onOpenIdea={onOpenIdea}
            onDeleted={onDeleted}
          />
        ) : (
          <HistoryTab detail={detail} />
        )}
      </div>

      <AnimatePresence>
        {editing && (
          <EditIdeaModal
            detail={detail}
            onClose={() => setEditing(false)}
            onSaved={() => {
              setEditing(false);
              onRefresh();
            }}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

function TabBtn({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className={`flex-1 rounded-full py-1.5 text-label-lg font-extrabold transition ${
        active ? "bg-white text-purple shadow-card" : "text-slatey"
      }`}
    >
      {children}
    </button>
  );
}

// --- about -------------------------------------------------------------------

function AboutTab({
  detail,
  onRefresh,
  onOpenIdea,
  onDeleted,
}: {
  detail: IdeaDetail;
  onRefresh: () => void;
  onOpenIdea: (id: number) => void;
  onDeleted: () => void;
}) {
  const { idea } = detail;
  const [archiving, setArchiving] = useState(false);
  const [confirmingDelete, setConfirmingDelete] = useState(false);
  const [merging, setMerging] = useState(false);
  const [reason, setReason] = useState("");
  const [mergeTarget, setMergeTarget] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const act = async (fn: () => Promise<unknown>) => {
    setBusy(true);
    setError("");
    try {
      await fn();
      onRefresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "that didn't work");
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      {idea.description && (
        <p className="whitespace-pre-wrap text-body-sm text-slatey">
          {idea.description}
        </p>
      )}

      <VoteBar idea={idea} />

      {/* Status as workflow, not as a label */}
      <Section icon={<Clock size={13} />} title="Status">
        <StatusStepper
          current={idea.status}
          flow={detail.statusFlow}
          onPick={(s) => act(() => api.updateIdea(idea.id, { status: s }))}
          disabled={busy}
        />
      </Section>

      <Section icon={<Users size={13} />} title="People">
        <div className="space-y-2">
          <Row label="Owner">
            <span className="flex items-center gap-1.5">
              <Avatar
                name={idea.owner.name}
                color={idea.owner.avatarColor}
                url={idea.owner.avatarUrl}
                size={20}
              />
              {idea.owner.name}
            </span>
          </Row>
          <Row label="Created">{timeAgo(idea.createdAt)} ago</Row>
          {detail.participants.length > 0 && (
            <div className="flex flex-wrap gap-1 pt-1">
              {detail.participants.map((p) => (
                <span
                  key={p.id}
                  title={p.name}
                  className="flex items-center gap-1 rounded-full bg-gray-50 px-2 py-0.5 text-label-sm font-bold text-slatey"
                >
                  <Avatar name={p.name} color={p.avatarColor} url={p.avatarUrl} size={14} />
                  {p.name}
                </span>
              ))}
            </div>
          )}
        </div>
      </Section>

      <Section icon={<Tag size={13} />} title="Tags">
        {idea.tags.length > 0 ? (
          <div className="flex flex-wrap gap-1.5">
            {idea.tags.map((t) => (
              <TagChip key={t} tag={t} />
            ))}
          </div>
        ) : (
          <p className="text-label-md text-gray-500">No tags yet.</p>
        )}
      </Section>

      <TasksSection detail={detail} onRefresh={onRefresh} />

      {detail.similar.length > 0 && (
        <Section icon={<Link2 size={13} />} title="Possibly related">
          <div className="space-y-1.5">
            {detail.similar.map((s) => (
              <button
                key={s.id}
                onClick={() => onOpenIdea(s.id)}
                className="flex w-full items-center gap-2 rounded-lg bg-gray-50 px-2.5 py-2 text-left transition hover:bg-gray-100"
              >
                <div className="min-w-0 flex-1">
                  <p className="truncate text-label-lg font-bold text-ink">
                    #{s.id} {s.title}
                  </p>
                  <p className="text-label-sm text-slatey">
                    {Math.round(s.similarity * 100)}% overlap · {s.score} votes
                  </p>
                </div>
                <ArrowRight size={13} className="shrink-0 text-slatey" />
              </button>
            ))}
          </div>
        </Section>
      )}

      {detail.mergedIn.length > 0 && (
        <Section icon={<GitMerge size={13} />} title="Merged in">
          {detail.mergedIn.map((m) => (
            <p key={m.id} className="text-label-md text-slatey">
              #{m.id} {m.title}
            </p>
          ))}
        </Section>
      )}

      {idea.archived && (
        <div className="rounded-lg bg-gray-50 p-3">
          <p className="text-label-lg font-extrabold text-ink">Archived</p>
          <p className="mt-0.5 text-body-sm text-slatey">{idea.archiveReason}</p>
        </div>
      )}

      {error && (
        <p className="rounded-lg bg-coral-light px-3 py-2 text-label-md text-coral">
          {error}
        </p>
      )}

      {/* Actions */}
      <div className="space-y-2 border-t border-gray-100 pt-4">
        {!idea.archived ? (
          <>
            <ActionButton
              icon={<CheckSquare size={14} />}
              onClick={() =>
                act(() => api.createIdeaTask(idea.id, { title: idea.title }))
              }
              disabled={busy}
            >
              Convert to task
            </ActionButton>

            <ActionButton
              icon={<GitMerge size={14} />}
              onClick={() => setMerging((m) => !m)}
              disabled={busy}
            >
              Merge into another idea
            </ActionButton>

            {merging && (
              <div className="rounded-lg bg-gray-50 p-2.5">
                <label className="text-label-md font-bold text-slatey">
                  Surviving idea number
                </label>
                <div className="mt-1 flex gap-2">
                  <input
                    value={mergeTarget}
                    onChange={(e) => setMergeTarget(e.target.value)}
                    placeholder="e.g. 12"
                    inputMode="numeric"
                    className="min-w-0 flex-1 rounded-lg bg-white px-2 py-1.5 text-body-sm outline-none ring-purple/30 focus:ring-2"
                  />
                  <button
                    onClick={() =>
                      act(async () => {
                        await api.mergeIdea(idea.id, Number(mergeTarget));
                        setMerging(false);
                        setMergeTarget("");
                      })
                    }
                    disabled={!mergeTarget || busy}
                    className="rounded-lg bg-purple px-3 py-1.5 text-label-lg font-extrabold text-white disabled:opacity-40"
                  >
                    Merge
                  </button>
                </div>
                <p className="mt-1.5 text-label-sm text-gray-500">
                  Votes and the whole thread move across. Anyone who backed both
                  is only counted once.
                </p>
              </div>
            )}

            <ActionButton
              icon={<Archive size={14} />}
              onClick={() => setArchiving((a) => !a)}
              disabled={busy}
              tone="warn"
            >
              Archive
            </ActionButton>

            {archiving && (
              <div className="rounded-lg bg-amber-light p-2.5">
                <label className="text-label-md font-bold text-ink">
                  Why is this being archived?
                </label>
                <textarea
                  value={reason}
                  onChange={(e) => setReason(e.target.value)}
                  rows={2}
                  placeholder="Duplicate of #12 / out of scope for v2 / superseded…"
                  className="mt-1 w-full resize-none rounded-lg bg-white p-2 text-body-sm outline-none ring-purple/30 focus:ring-2"
                />
                <p className="mt-1 text-label-sm text-slatey">
                  Required — it&apos;s what tells everyone who contributed what
                  happened to their idea.
                </p>
                <button
                  onClick={() =>
                    act(async () => {
                      await api.archiveIdea(idea.id, reason.trim());
                      setArchiving(false);
                      setReason("");
                    })
                  }
                  disabled={!reason.trim() || busy}
                  className="mt-2 w-full rounded-full bg-amber py-2 text-label-lg font-extrabold text-ink disabled:opacity-40"
                >
                  Archive idea
                </button>
              </div>
            )}
          </>
        ) : (
          <ActionButton
            icon={<RotateCcw size={14} />}
            onClick={() => act(() => api.restoreIdea(idea.id))}
            disabled={busy}
          >
            Restore to the board
          </ActionButton>
        )}

        {detail.canEdit && (
          <ActionButton
            icon={<Trash2 size={14} />}
            tone="danger"
            disabled={busy}
            onClick={() => setConfirmingDelete(true)}
          >
            Delete permanently
          </ActionButton>
        )}
      </div>

      {/* Deleting is the owner's call, but it takes other people's
          contributions with it — so the dialog counts exactly what's lost and
          offers archiving, which keeps the thread readable, as the way out. */}
      <ConfirmDialog
        open={confirmingDelete}
        danger
        title={`Delete "${truncate(idea.title, 40)}"?`}
        message={deletionWarning(detail)}
        confirmLabel="Delete forever"
        cancelLabel="Keep it"
        onCancel={() => setConfirmingDelete(false)}
        onConfirm={() => {
          setConfirmingDelete(false);
          act(async () => {
            await api.deleteIdea(idea.id);
            onDeleted();
          });
        }}
      />
    </>
  );
}

/** Spells out what deletion destroys, in the concrete rather than the abstract. */
function deletionWarning(detail: IdeaDetail): string {
  const { idea } = detail;
  const losses: string[] = [];
  if (idea.messageCount > 0) {
    losses.push(
      `${idea.messageCount} ${idea.messageCount === 1 ? "message" : "messages"}`
    );
  }
  const votes = idea.upvotes + idea.downvotes;
  if (votes > 0) losses.push(`${votes} ${votes === 1 ? "vote" : "votes"}`);
  if (detail.tasks.length > 0) {
    losses.push(
      `${detail.tasks.length} linked ${detail.tasks.length === 1 ? "task" : "tasks"}`
    );
  }

  if (losses.length === 0) {
    return "This can't be undone. Nobody has engaged with it yet, so nothing else goes with it.";
  }
  const from =
    detail.participants.length > 1
      ? ` contributed by ${detail.participants.length} people`
      : "";
  return `This permanently destroys ${listOf(losses)}${from}. It can't be undone — archiving keeps the discussion readable instead.`;
}

function listOf(items: string[]): string {
  if (items.length === 1) return items[0];
  return `${items.slice(0, -1).join(", ")} and ${items[items.length - 1]}`;
}

function truncate(s: string, n: number): string {
  return s.length > n ? `${s.slice(0, n)}…` : s;
}

function StatusStepper({
  current,
  flow,
  onPick,
  disabled,
}: {
  current: IdeaStatus;
  flow: IdeaStatus[];
  onPick: (s: IdeaStatus) => void;
  disabled?: boolean;
}) {
  // Archived isn't a step on the ladder — it's how an idea leaves it, and it
  // has its own button. The annotation keeps the array widened: TypeScript
  // otherwise infers a narrowing predicate from the filter and drops
  // "archived" from the element type, which then can't hold `current`.
  const steps: IdeaStatus[] = flow.filter((s) => s !== "archived");
  const currentIndex = steps.indexOf(current);

  return (
    <div className="space-y-1">
      {steps.map((s, i) => {
        const done = currentIndex > i;
        const active = current === s;
        const meta = STATUS_META[s];
        return (
          <button
            key={s}
            onClick={() => onPick(s)}
            disabled={disabled || active}
            className={`flex w-full items-center gap-2 rounded-lg px-2.5 py-1.5 text-left text-label-lg font-bold transition disabled:cursor-default ${
              active ? "" : "hover:bg-gray-50"
            }`}
            style={active ? { background: meta.bg, color: meta.tint } : undefined}
          >
            <span
              className="flex h-4 w-4 shrink-0 items-center justify-center rounded-full border-2"
              style={{
                borderColor: active || done ? meta.tint : "#CCCCCC",
                background: done ? meta.tint : "transparent",
              }}
            />
            <span className={active ? "" : done ? "text-slatey" : "text-gray-500"}>
              {meta.label}
            </span>
          </button>
        );
      })}
    </div>
  );
}

function TasksSection({
  detail,
  onRefresh,
}: {
  detail: IdeaDetail;
  onRefresh: () => void;
}) {
  const toggle = async (task: IdeaTask) => {
    await api
      .updateIdeaTask(task.id, { status: task.status === "done" ? "todo" : "done" })
      .catch(() => {});
    onRefresh();
  };

  return (
    <Section icon={<ListTodo size={13} />} title="Linked tasks">
      {detail.tasks.length === 0 ? (
        <p className="text-label-md text-gray-500">
          None yet — converting this idea creates one.
        </p>
      ) : (
        <div className="space-y-1">
          {detail.tasks.map((t) => (
            <button
              key={t.id}
              onClick={() => toggle(t)}
              className="flex w-full items-start gap-2 rounded-lg px-1 py-1 text-left transition hover:bg-gray-50"
            >
              {t.status === "done" ? (
                <CheckSquare size={14} className="mt-0.5 shrink-0 text-teal" />
              ) : (
                <Square size={14} className="mt-0.5 shrink-0 text-gray-300" />
              )}
              <span
                className={`text-body-sm ${
                  t.status === "done" ? "text-gray-500 line-through" : "text-ink"
                }`}
              >
                {t.title}
                {t.sprint && (
                  <span className="ml-1 text-label-sm text-slatey">· {t.sprint}</span>
                )}
              </span>
            </button>
          ))}
        </div>
      )}
    </Section>
  );
}

// --- history -----------------------------------------------------------------

function HistoryTab({ detail }: { detail: IdeaDetail }) {
  if (detail.history.length === 0) {
    return <p className="text-body-sm text-slatey">Nothing recorded yet.</p>;
  }
  return (
    <div className="space-y-3">
      <p className="flex items-center gap-1.5 text-label-md font-bold text-slatey">
        <History size={13} /> Every change, and who made it
      </p>
      {detail.history.map((e) => (
        <div key={e.id} className="flex gap-2.5">
          <div className="flex flex-col items-center">
            <span className="mt-1 h-2 w-2 shrink-0 rounded-full bg-purple" />
            <span className="w-px flex-1 bg-gray-100" />
          </div>
          <div className="min-w-0 flex-1 pb-2">
            <p className="text-body-sm text-ink">
              <strong>{e.actor?.name ?? "Someone"}</strong> {describeEvent(e)}
            </p>
            {e.note && (
              <p className="mt-0.5 text-label-md italic text-slatey">{e.note}</p>
            )}
            <p className="mt-0.5 text-label-sm text-gray-500">{timeAgo(e.at)} ago</p>
          </div>
        </div>
      ))}
    </div>
  );
}

function describeEvent(e: IdeaDetail["history"][number]): string {
  switch (e.kind) {
    case "created":
      return `posted this idea`;
    case "status":
      return `moved it from ${labelOf(e.from)} to ${labelOf(e.to)}`;
    case "vote_threshold":
      return `— it reached the review threshold and moved to ${labelOf(e.to)}`;
    case "edited":
      return `edited the ${e.field}`;
    case "tagged":
      return `changed the tags to ${e.to || "none"}`;
    case "merged":
      return `merged "${e.from}" into "${e.to}"`;
    case "archived":
      return `archived it`;
    case "restored":
      return `restored it to the board`;
    case "task":
      return `converted it to a task: "${e.to}"`;
    case "linked":
      return `linked it to "${e.to}"`;
    case "brainstorm":
      return `started a ${e.to} silent brainstorm`;
    default:
      return e.kind;
  }
}

function labelOf(status: string): string {
  return STATUS_META[status as IdeaStatus]?.label ?? status;
}

// --- edit modal --------------------------------------------------------------

function EditIdeaModal({
  detail,
  onClose,
  onSaved,
}: {
  detail: IdeaDetail;
  onClose: () => void;
  onSaved: () => void;
}) {
  const [title, setTitle] = useState(detail.idea.title);
  const [description, setDescription] = useState(detail.idea.description);
  const [tags, setTags] = useState(detail.idea.tags.join(", "));
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const save = async () => {
    setBusy(true);
    setError("");
    try {
      await api.updateIdea(detail.idea.id, {
        title: title.trim(),
        description: description.trim(),
        tags: tags
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean),
      });
      onSaved();
    } catch (e) {
      setError(e instanceof Error ? e.message : "could not save");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/50 p-4">
      <motion.div
        initial={{ opacity: 0, scale: 0.96 }}
        animate={{ opacity: 1, scale: 1 }}
        className="w-full max-w-md rounded-2xl bg-white p-5 shadow-card-lg"
      >
        <h3 className="text-heading-sm font-extrabold text-ink">Edit idea</h3>

        <label className="mt-4 block text-label-md font-bold text-slatey">Title</label>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="mt-1 w-full rounded-lg bg-gray-50 px-3 py-2 text-body-md outline-none ring-purple/30 focus:ring-2"
        />

        <label className="mt-3 block text-label-md font-bold text-slatey">
          Description
        </label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={4}
          className="mt-1 w-full resize-none rounded-lg bg-gray-50 px-3 py-2 text-body-md outline-none ring-purple/30 focus:ring-2"
        />

        <label className="mt-3 block text-label-md font-bold text-slatey">
          Tags (comma separated)
        </label>
        <input
          value={tags}
          onChange={(e) => setTags(e.target.value)}
          placeholder="ai, ux, v2.0"
          className="mt-1 w-full rounded-lg bg-gray-50 px-3 py-2 text-body-md outline-none ring-purple/30 focus:ring-2"
        />

        {error && (
          <p className="mt-3 rounded-lg bg-coral-light px-3 py-2 text-label-md text-coral">
            {error}
          </p>
        )}

        <p className="mt-3 text-label-sm text-gray-500">
          Every edit is recorded in the history, so the thread can always be read
          against what the idea said at the time.
        </p>

        <div className="mt-4 flex gap-3">
          <Button variant="outline" className="flex-1" onClick={onClose}>
            Cancel
          </Button>
          <Button className="flex-1" onClick={save} loading={busy}>
            Save
          </Button>
        </div>
      </motion.div>
    </div>
  );
}

// --- bits --------------------------------------------------------------------

function Section({
  icon,
  title,
  children,
}: {
  icon: React.ReactNode;
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <p className="mb-1.5 flex items-center gap-1.5 text-label-md font-extrabold uppercase tracking-wide text-slatey">
        {icon} {title}
      </p>
      {children}
    </div>
  );
}

function Row({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between text-body-sm">
      <span className="text-slatey">{label}</span>
      <span className="font-bold text-ink">{children}</span>
    </div>
  );
}

function ActionButton({
  icon,
  onClick,
  disabled,
  tone,
  children,
}: {
  icon: React.ReactNode;
  onClick: () => void;
  disabled?: boolean;
  tone?: "warn" | "danger";
  children: React.ReactNode;
}) {
  const toneClass =
    tone === "danger"
      ? "text-coral hover:bg-coral-light"
      : tone === "warn"
        ? "text-amber hover:bg-amber-light"
        : "text-slatey hover:bg-gray-50 hover:text-ink";
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-label-lg font-bold transition disabled:opacity-40 ${toneClass}`}
    >
      {icon} {children}
    </button>
  );
}
