"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { AnimatePresence } from "framer-motion";
import {
  CircleHelp,
  Info,
  List,
  MessagesSquare,
  PanelRight,
  UserX,
} from "lucide-react";
import { AppShell } from "@/components/AppShell";
import { IdeasListPanel } from "@/components/ideas/IdeasListPanel";
import { IdeaThreadPanel } from "@/components/ideas/IdeaThreadPanel";
import { IdeaDetailsPanel } from "@/components/ideas/IdeaDetailsPanel";
import { NewIdeaModal } from "@/components/ideas/NewIdeaModal";
import { IdeasIntro, useIdeasIntro } from "@/components/ideas/IdeasIntro";
import { api } from "@/lib/api";
import type { ChatUser, IdeaBoard, IdeaDetail, IdeaThread } from "@/lib/types";

/**
 * The ideas workspace: board on the left, thread in the middle, details on the
 * right.
 *
 * Three panels side by side on a wide screen. On a phone that's unreadable, so
 * it collapses to one panel at a time with the selection driving which — list
 * until you pick an idea, then the thread, with details a tap away.
 */
function IdeasWorkspace() {
  const router = useRouter();
  const params = useSearchParams();

  const [board, setBoard] = useState<IdeaBoard | null>(null);
  const [detail, setDetail] = useState<IdeaDetail | null>(null);
  const [thread, setThread] = useState<IdeaThread | null>(null);
  const [contacts, setContacts] = useState<ChatUser[]>([]);

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [boardLoading, setBoardLoading] = useState(true);
  const [ideaLoading, setIdeaLoading] = useState(false);

  const [filter, setFilter] = useState("");
  const [sort, setSort] = useState("hot");
  const [tag, setTag] = useState("");
  const [query, setQuery] = useState("");

  const [creating, setCreating] = useState(false);
  const [mobilePane, setMobilePane] = useState<"list" | "thread" | "details">("list");
  const [showIntro, dismissIntro, replayIntro] = useIdeasIntro();

  // A notification deep-links straight to an idea (/ideas?idea=12).
  useEffect(() => {
    const id = Number(params.get("idea"));
    if (id > 0) {
      setSelectedId(id);
      setMobilePane("thread");
    }
  }, [params]);

  const loadBoard = useCallback(async () => {
    setBoardLoading(true);
    try {
      setBoard(await api.ideas({ status: filter, sort, tag, q: query }));
    } catch {
      setBoard(null);
    } finally {
      setBoardLoading(false);
    }
  }, [filter, sort, tag, query]);

  // Debounced so typing in the search box isn't a request per keystroke.
  useEffect(() => {
    const t = setTimeout(loadBoard, query ? 350 : 0);
    return () => clearTimeout(t);
  }, [loadBoard, query]);

  useEffect(() => {
    api
      .chatContacts()
      .then((r) => setContacts(r.contacts))
      .catch(() => setContacts([]));
  }, []);

  const loadIdea = useCallback(async (id: number, showSpinner = true) => {
    if (showSpinner) setIdeaLoading(true);
    try {
      const [d, t] = await Promise.all([api.idea(id), api.ideaMessages(id)]);
      setDetail(d);
      setThread(t);
    } catch {
      setDetail(null);
      setThread(null);
    } finally {
      setIdeaLoading(false);
    }
  }, []);

  useEffect(() => {
    if (selectedId) loadIdea(selectedId);
    else {
      setDetail(null);
      setThread(null);
    }
  }, [selectedId, loadIdea]);

  const openIdea = (id: number) => {
    setSelectedId(id);
    setMobilePane("thread");
    router.replace(`/ideas?idea=${id}`, { scroll: false });
  };

  // Votes and stars update the row optimistically — a one-tap action that waits
  // on a round trip doesn't feel like one tap.
  const vote = async (id: number, value: number) => {
    setBoard((b) =>
      b
        ? {
            ...b,
            ideas: b.ideas.map((i) =>
              i.id === id
                ? {
                    ...i,
                    myVote: value,
                    score: i.score - i.myVote + value,
                    upvotes: i.upvotes - (i.myVote === 1 ? 1 : 0) + (value === 1 ? 1 : 0),
                    downvotes:
                      i.downvotes - (i.myVote === -1 ? 1 : 0) + (value === -1 ? 1 : 0),
                  }
                : i
            ),
          }
        : b
    );
    await api.voteIdea(id, value).catch(() => {});
    loadBoard();
    if (selectedId === id) loadIdea(id, false);
  };

  const star = async (id: number) => {
    setBoard((b) =>
      b
        ? {
            ...b,
            ideas: b.ideas.map((i) =>
              i.id === id ? { ...i, starred: !i.starred } : i
            ),
          }
        : b
    );
    await api.starIdea(id).catch(() => {});
  };

  const refreshIdea = () => {
    if (selectedId) loadIdea(selectedId, false);
    loadBoard();
  };

  const startBrainstorm = async () => {
    if (!selectedId) return;
    await api.startBrainstorm(selectedId, 10, "").catch(() => {});
    refreshIdea();
  };

  return (
    <AppShell tabs wide="full">
      <div className="flex h-[calc(100dvh-4rem-env(safe-area-inset-bottom))] min-h-0 flex-col lg:h-[calc(100dvh-3rem)]">
        {/* Toolbar */}
        <div className="flex shrink-0 items-center gap-2 border-b border-gray-100 bg-white px-4 py-2">
          <h1 className="text-heading-sm font-extrabold text-ink">Ideas</h1>
          <span className="hidden text-label-md text-slatey sm:inline">
            One thread per idea — team chat stays in Chat.
          </span>
          <button
            onClick={replayIntro}
            aria-label="How Ideas work"
            title="How Ideas work"
            className="rounded-full p-1 text-gray-300 transition hover:bg-gray-50 hover:text-slatey"
          >
            <CircleHelp size={15} />
          </button>

          {selectedId && !thread?.brainstorm && (
            <button
              onClick={startBrainstorm}
              title="Everyone posts anonymously for 10 minutes"
              className="ml-auto flex items-center gap-1.5 rounded-full bg-ink px-3 py-1.5 text-label-md font-bold text-white transition hover:opacity-90"
            >
              <UserX size={13} />
              <span className="hidden md:inline">Silent brainstorm</span>
            </button>
          )}

          {selectedId && (
            <div
              className={`flex gap-0.5 rounded-full bg-gray-50 p-0.5 xl:hidden ${
                thread?.brainstorm ? "ml-auto" : ""
              }`}
            >
              {(
                [
                  { key: "list", label: "Ideas", icon: List, hideAt: "lg" },
                  { key: "thread", label: "Thread", icon: MessagesSquare, hideAt: "lg" },
                  { key: "details", label: "Details", icon: PanelRight, hideAt: "xl" },
                ] as const
              ).map((p) => (
                <button
                  key={p.key}
                  onClick={() => setMobilePane(p.key)}
                  aria-pressed={mobilePane === p.key}
                  aria-label={p.label}
                  title={p.label}
                  // List and Thread are both on screen from lg up, so only
                  // Details still needs a switch there.
                  className={`flex items-center gap-1.5 rounded-full px-2.5 py-1.5 text-label-md font-bold transition ${
                    p.hideAt === "lg" ? "lg:hidden" : ""
                  } ${
                    mobilePane === p.key
                      ? "bg-white text-purple shadow-card"
                      : "text-slatey hover:text-ink"
                  }`}
                >
                  <p.icon size={13} />
                  <span className="hidden sm:inline">{p.label}</span>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Panels */}
        {/* minmax(0,1fr), not 1fr: a bare 1fr track is minmax(auto, 1fr), so
            the thread's own content sets a floor the column can't shrink below
            — which pushes the details panel off the right edge of the screen. */}
        <div className="grid min-h-0 flex-1 grid-cols-1 lg:grid-cols-[300px_minmax(0,1fr)] xl:grid-cols-[300px_minmax(0,1fr)_320px]">
          <div
            className={`min-h-0 ${
              mobilePane === "list" ? "block" : "hidden"
            } lg:block`}
          >
            <IdeasListPanel
              board={board}
              loading={boardLoading}
              selectedId={selectedId}
              filter={filter}
              sort={sort}
              tag={tag}
              query={query}
              onSelect={openIdea}
              onFilter={setFilter}
              onSort={setSort}
              onTag={setTag}
              onQuery={setQuery}
              onVote={vote}
              onStar={star}
              onNew={() => setCreating(true)}
            />
          </div>

          <div
            className={`min-h-0 ${
              mobilePane === "thread" ? "block" : "hidden"
            } lg:block`}
          >
            <IdeaThreadPanel
              thread={thread}
              loading={ideaLoading}
              contacts={contacts}
              onRefresh={refreshIdea}
              onOpenIdea={openIdea}
            />
          </div>

          <div
            className={`min-h-0 ${
              mobilePane === "details" ? "block" : "hidden"
            } xl:block`}
          >
            <IdeaDetailsPanel
              detail={detail}
              loading={ideaLoading && !detail}
              onRefresh={refreshIdea}
              onOpenIdea={openIdea}
              onDeleted={() => {
                setSelectedId(null);
                setMobilePane("list");
                loadBoard();
              }}
              onClose={() => setMobilePane("thread")}
            />
          </div>
        </div>
      </div>

      <AnimatePresence>
        {showIntro && <IdeasIntro onClose={dismissIntro} />}
      </AnimatePresence>

      <AnimatePresence>
        {creating && (
          <NewIdeaModal
            onClose={() => setCreating(false)}
            onCreated={(id) => {
              setCreating(false);
              loadBoard();
              openIdea(id);
            }}
            onOpenIdea={(id) => {
              setCreating(false);
              openIdea(id);
            }}
          />
        )}
      </AnimatePresence>
    </AppShell>
  );
}

export default function IdeasPage() {
  // useSearchParams needs a Suspense boundary for the static build.
  return (
    <Suspense
      fallback={
        <div className="flex min-h-[60vh] items-center justify-center text-body-md text-slatey">
          <Info size={16} className="mr-2" /> Loading ideas…
        </div>
      }
    >
      <IdeasWorkspace />
    </Suspense>
  );
}
