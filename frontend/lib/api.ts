import type {
  User,
  Lesson,
  Skill,
  UserQuest,
  CharacterWithFriendship,
  LeaderRow,
  LeagueStandings,
  LeagueResult,
  LeagueHistoryEntry,
  HomeData,
  ListeningSession,
  ReadingSession,
  VocabItem,
  Mistake,
  AppNotification,
  Certificate,
  ExamResult,
  ExamMeta,
  PaymentStatus,
  HeartsStatus,
  ExamPaper,
  CertVerification,
  ChatUser,
  ChatMessage,
  ChatThread,
  Idea,
  IdeaBoard,
  IdeaDetail,
  IdeaMessage,
  IdeaStatus,
  WritingCheck,
  IdeaTask,
  IdeaThread,
  SimilarIdea,
  ThreadSummary,
  BrainstormSession,
  MessageTranslation,
} from "./types";

export const API_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

/** Resolve a backend-relative media path (e.g. an avatar URL) to an absolute URL. */
export function mediaUrl(path?: string): string {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  return `${API_URL}${path}`;
}

const TOKEN_KEY = "lumora_token";
const USER_KEY = "lumora_user";
const ROUTE_KEY = "lumora_last_route";

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string) {
  window.localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  window.localStorage.removeItem(TOKEN_KEY);
}

/**
 * The signed-in user is cached in localStorage so a full page reload (e.g. the
 * user manually editing the URL) can rehydrate state instantly instead of
 * dropping to a logged-out view and rerouting through onboarding.
 */
export function getStoredUser(): User | null {
  if (typeof window === "undefined") return null;
  try {
    const raw = window.localStorage.getItem(USER_KEY);
    return raw ? (JSON.parse(raw) as User) : null;
  } catch {
    return null;
  }
}

export function setStoredUser(user: User | null) {
  if (typeof window === "undefined") return;
  if (user) window.localStorage.setItem(USER_KEY, JSON.stringify(user));
  else window.localStorage.removeItem(USER_KEY);
}

/** The last authenticated screen the user was on, so we can restore it. */
export function getLastRoute(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(ROUTE_KEY);
}

export function setLastRoute(path: string) {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(ROUTE_KEY, path);
}

export function clearSession() {
  clearToken();
  setStoredUser(null);
}

class ApiError extends Error {
  status: number;
  /** True only when the session in use was rejected — i.e. the user must sign in again. */
  sessionExpired: boolean;
  constructor(message: string, status: number, sessionExpired = false) {
    super(message);
    this.status = status;
    this.sessionExpired = sessionExpired;
  }
}

const NETWORK_ERROR =
  "We can't reach Lumora right now. Check your internet connection and try again.";

/** fetch, with a network failure turned into a readable ApiError. */
async function send(url: string, init: RequestInit): Promise<Response> {
  try {
    return await fetch(url, init);
  } catch {
    throw new ApiError(NETWORK_ERROR, 0);
  }
}

/** What to tell the user when the server gave no message of its own. */
function fallbackMessage(status: number): string {
  if (status === 401) return "Please sign in to continue.";
  if (status === 403) return "You don't have access to that.";
  if (status === 404) return "We couldn't find what you were looking for.";
  if (status === 429) return "Too many attempts. Please wait a moment and try again.";
  if (status >= 500)
    return "Lumora is having trouble right now — it may be waking up. Please try again in a moment.";
  return "Something went wrong. Please try again.";
}

/**
 * Turns a failed response into an ApiError carrying a message fit to show the
 * user: the server's own message when it sent one, otherwise a plain-language
 * fallback for the status.
 *
 * A 401 only ends the session when it's about the session in use right now.
 * Sign-in with a wrong password is a 401 too, and so is a request that went out
 * with an older token just before the user signed in again — neither of those
 * means "your session expired", and treating them that way showed a bare
 * "unauthorized" on the login form.
 */
async function toApiError(res: Response, sentToken: string | null): Promise<ApiError> {
  let serverMessage = "";
  try {
    serverMessage = (await res.json())?.error || "";
  } catch {
    /* not JSON — e.g. the host's own error page while the API wakes up */
  }
  if (res.status === 401 && sentToken && sentToken === getToken()) {
    clearSession();
    return new ApiError("Your session has ended. Please sign in again.", 401, true);
  }
  return new ApiError(serverMessage || fallbackMessage(res.status), res.status);
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };
  const token = getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await send(`${API_URL}${path}`, { ...options, headers });
  if (!res.ok) throw await toApiError(res, token);
  return (await res.json()) as T;
}

/**
 * Guarantees an idea message's collection fields are arrays.
 *
 * The API is supposed to emit `[]` rather than `null` for these, and does — but
 * a nil slice in Go marshals to `null`, so a single missed initialisation on
 * the server (or an older backend still running against a newer client) turns
 * `message.replies.length` into a crash that takes the whole thread down.
 * Normalising once at the boundary means no render site has to guard.
 */
function normaliseMessage(m: IdeaMessage): IdeaMessage {
  return {
    ...m,
    reactions: m.reactions ?? [],
    replies: (m.replies ?? []).map(normaliseMessage),
  };
}

/**
 * Multipart sibling of `request`, for attachments.
 *
 * Deliberately does NOT set Content-Type: the browser has to write it itself so
 * it can append the multipart boundary. Setting it by hand produces a request
 * the server can't parse.
 */
async function upload<T>(
  path: string,
  fields: { file: File } & Record<string, string | File | undefined>
): Promise<T> {
  const fd = new FormData();
  for (const [key, value] of Object.entries(fields)) {
    if (value === undefined || value === "") continue;
    fd.append(key, value as string | File);
  }

  const token = getToken();
  const res = await send(`${API_URL}${path}`, {
    method: "POST",
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: fd,
  });
  if (!res.ok) throw await toApiError(res, token);
  return (await res.json()) as T;
}

export const api = {
  // Auth
  register: (email: string, password: string, name: string) =>
    request<{ token: string; user: User }>("/api/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, name }),
    }),

  login: (email: string, password: string) =>
    request<{ token: string; user: User }>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  me: () => request<{ user: User }>("/api/auth/me"),

  forgotPassword: (email: string) =>
    request<{ ok: boolean }>("/api/auth/forgot-password", {
      method: "POST",
      body: JSON.stringify({ email }),
    }),

  resetPassword: (token: string, password: string) =>
    request<{ ok: boolean }>("/api/auth/reset-password", {
      method: "POST",
      body: JSON.stringify({ token, password }),
    }),

  setup: (targetLanguage: string, dailyGoalXp: number, reason: string) =>
    request<{ user: User }>("/api/auth/setup", {
      method: "POST",
      body: JSON.stringify({ targetLanguage, dailyGoalXp, reason }),
    }),

  updateProfile: (body: {
    name?: string;
    avatarColor?: string;
    dailyGoalXp?: number;
  }) =>
    request<{ user: User }>("/api/auth/profile", {
      method: "PATCH",
      body: JSON.stringify(body),
    }),

  uploadAvatar: async (file: File): Promise<{ user: User }> => {
    const fd = new FormData();
    fd.append("file", file);
    const token = getToken();
    const res = await send(`${API_URL}/api/auth/avatar`, {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: fd,
    });
    if (!res.ok) throw await toApiError(res, token);
    return res.json();
  },

  removeAvatar: () =>
    request<{ user: User }>("/api/auth/avatar", { method: "DELETE" }),

  changePassword: (currentPassword: string, newPassword: string) =>
    request<{ ok: boolean }>("/api/auth/password", {
      method: "POST",
      body: JSON.stringify({ currentPassword, newPassword }),
    }),

  deleteAccount: (password: string) =>
    request<{ ok: boolean }>("/api/auth/account", {
      method: "DELETE",
      body: JSON.stringify({ password }),
    }),

  // Content & progress
  home: () => request<HomeData>("/api/home"),

  skills: () => request<{ skills: Skill[] }>("/api/skills"),

  lesson: (id: number | string) =>
    request<{ lesson: Lesson }>(`/api/lessons/${id}`),

  /** Second opinion on a typed answer marked wrong (alternative wordings). */
  checkAnswer: (input: {
    lessonId: number;
    question: string;
    expected: string;
    answer: string;
  }) =>
    request<{ available: boolean; correct?: boolean; explanation?: string }>(
      "/api/lessons/check-answer",
      { method: "POST", body: JSON.stringify(input) }
    ),

  /** Rule checks + corrections for a free-writing answer. */
  checkWriting: (exerciseId: number, text: string) =>
    request<WritingCheck>(`/api/exercises/${exerciseId}/check-writing`, {
      method: "POST",
      body: JSON.stringify({ text }),
    }),

  completeLesson: (id: number | string, accuracy: number) =>
    request<{
      xpEarned: number;
      accuracy: number;
      user: User;
      firstClear: boolean;
    }>(`/api/lessons/${id}/complete`, {
      method: "POST",
      body: JSON.stringify({ accuracy }),
    }),

  listeningSessions: () =>
    request<{ sessions: ListeningSession[] }>("/api/listening"),

  listeningSession: (id: number | string) =>
    request<{ session: ListeningSession }>(`/api/listening/${id}`),

  completeListening: (id: number | string) =>
    request<{ xpEarned: number; user: User }>(`/api/listening/${id}/complete`, {
      method: "POST",
    }),

  readingSessions: () =>
    request<{ sessions: ReadingSession[] }>("/api/reading"),

  readingSession: (id: number | string) =>
    request<{ session: ReadingSession }>(`/api/reading/${id}`),

  completeReading: (id: number | string) =>
    request<{ xpEarned: number; user: User }>(`/api/reading/${id}/complete`, {
      method: "POST",
    }),

  enrollments: () =>
    request<{ languages: string[]; active: string }>("/api/enrollments"),

  enrollLanguage: (language: string) =>
    request<{ languages: string[]; active: string; user: User }>(
      "/api/enrollments",
      { method: "POST", body: JSON.stringify({ language }) }
    ),

  switchLanguage: (language: string) =>
    request<{ languages: string[]; active: string; user: User }>(
      "/api/enrollments/active",
      { method: "POST", body: JSON.stringify({ language }) }
    ),

  /** Takes a language off the user's courses (progress is kept server-side). */
  removeLanguage: (language: string) =>
    request<{ languages: string[]; active: string; user: User }>(
      `/api/enrollments/${encodeURIComponent(language)}`,
      { method: "DELETE" }
    ),

  practice: () =>
    request<{
      vocab: VocabItem[];
      mistakes: Mistake[];
      listeningCount: number;
      readingCount: number;
    }>("/api/practice"),

  practiceListening: () =>
    request<{ sessions: ListeningSession[] }>("/api/practice/listening"),

  practiceReading: () =>
    request<{ sessions: ReadingSession[] }>("/api/practice/reading"),

  recordMistake: (m: { prompt: string; question: string; correctAnswer: string }) =>
    request<{ ok: boolean }>("/api/mistakes", {
      method: "POST",
      body: JSON.stringify(m),
    }),

  resolveMistakes: (ids: number[]) =>
    request<{ ok: boolean }>("/api/mistakes/resolve", {
      method: "POST",
      body: JSON.stringify({ ids }),
    }),

  completePractice: (xp: number) =>
    request<{ xpEarned: number; user: User }>("/api/practice/complete", {
      method: "POST",
      body: JSON.stringify({ xp }),
    }),

  notifications: () =>
    request<{ notifications: AppNotification[]; unread: number }>(
      "/api/notifications"
    ),

  markNotificationsRead: () =>
    request<{ ok: boolean }>("/api/notifications/read", { method: "POST" }),

  markNotificationRead: (id: number | string) =>
    request<{ ok: boolean; unread: number }>(
      `/api/notifications/${id}/read`,
      { method: "POST" }
    ),

  deleteNotification: (id: number | string) =>
    request<{ ok: boolean; unread: number }>(`/api/notifications/${id}`, {
      method: "DELETE",
    }),

  submitExam: (body: {
    language: string;
    level: string;
    listening: number;
    reading: number;
    writing: number;
    speaking: number;
  }) =>
    request<ExamResult>("/api/exam/submit", {
      method: "POST",
      body: JSON.stringify(body),
    }),

  // Payments (Paystack)
  paymentStatus: () => request<PaymentStatus>("/api/payments/status"),

  initializePayment: (level: string) =>
    request<{ authorizationUrl?: string; reference?: string }>(
      "/api/payments/initialize",
      { method: "POST", body: JSON.stringify({ level }) }
    ),

  verifyPayment: (reference: string) =>
    request<{
      status: string;
      success: boolean;
      level: string;
      product: string;
    }>(`/api/payments/verify?reference=${encodeURIComponent(reference)}`),

  examMeta: () => request<ExamMeta>("/api/exam/meta"),

  startExam: (level: string, language: string) =>
    request<{ ok: boolean }>("/api/exam/start", {
      method: "POST",
      body: JSON.stringify({ level, language }),
    }),

  // Hearts
  heartsStatus: () => request<HeartsStatus>("/api/hearts"),
  loseHeart: () => request<HeartsStatus>("/api/hearts/lose", { method: "POST" }),
  buyHearts: () =>
    request<{ authorizationUrl?: string; reference?: string }>(
      "/api/payments/initialize",
      { method: "POST", body: JSON.stringify({ product: "hearts" }) }
    ),

  examPaper: (level: string) =>
    request<ExamPaper>(`/api/exam/paper?level=${encodeURIComponent(level)}`),

  certificates: () =>
    request<{ certificates: Certificate[] }>("/api/certificates"),

  certificate: (id: number | string) =>
    request<{ certificate: Certificate }>(`/api/certificates/${id}`),

  deleteCertificate: (id: number | string) =>
    request<{ ok: boolean }>(`/api/certificates/${id}`, { method: "DELETE" }),

  // Public — no auth header needed; anyone can verify a certificate by serial.
  verifyCertificate: async (serial: string): Promise<CertVerification> => {
    try {
      const res = await fetch(
        `${API_URL}/api/verify/${encodeURIComponent(serial)}`
      );
      if (!res.ok) return { valid: false };
      return (await res.json()) as CertVerification;
    } catch {
      return { valid: false };
    }
  },

  chatContacts: () =>
    request<{ contacts: ChatUser[] }>("/api/chat/contacts"),

  chatThreads: () => request<{ threads: ChatThread[] }>("/api/chat/threads"),

  chatUnread: () => request<{ count: number }>("/api/chat/unread"),

  chatMessages: (id: number | string) =>
    request<{ messages: ChatMessage[]; user: ChatUser }>(`/api/chat/with/${id}`),

  sendChatMessage: (id: number | string, body: string) =>
    request<{ message: ChatMessage }>(`/api/chat/with/${id}`, {
      method: "POST",
      body: JSON.stringify({ body }),
    }),

  quests: () => request<{ quests: UserQuest[] }>("/api/quests/daily"),

  characters: () =>
    request<{ characters: CharacterWithFriendship[] }>("/api/characters"),

  leaderboard: () =>
    request<{ league: string; rows: LeaderRow[]; userRank: number }>(
      "/api/leaderboard"
    ),

  // --- the weekly league ---
  league: () => request<LeagueStandings>("/api/league"),

  /** The most recent settled season, if its ceremony hasn't been played yet. */
  leagueResult: () =>
    request<{ result: LeagueResult | null }>("/api/league/result"),

  markLeagueResultSeen: (seasonId: string) =>
    request<{ ok: boolean }>(
      `/api/league/result/seen?season=${encodeURIComponent(seasonId)}`,
      { method: "POST" }
    ),

  leagueHistory: () =>
    request<{ history: LeagueHistoryEntry[] }>("/api/league/history"),

  setLeagueCasual: (enabled: boolean) =>
    request<{ ok: boolean; casual: boolean }>("/api/league/casual", {
      method: "POST",
      body: JSON.stringify({ enabled }),
    }),

  reportLeagueMember: (id: number, reason: string) =>
    request<{ ok: boolean; alreadyReported: boolean }>(
      `/api/league/report/${id}`,
      { method: "POST", body: JSON.stringify({ reason }) }
    ),

  // --- direct messages: edit, delete, photos ---
  sendChatImage: (id: number | string, file: File, caption: string) =>
    upload<{ message: ChatMessage }>(`/api/chat/with/${id}/image`, {
      file,
      body: caption,
    }),

  editChatMessage: (messageId: number, body: string) =>
    request<{ message: ChatMessage }>(`/api/chat/messages/${messageId}`, {
      method: "PATCH",
      body: JSON.stringify({ body }),
    }),

  /** Retry translating a message the background pass didn't cover. */
  translateChatMessage: (messageId: number) =>
    request<{ translation: MessageTranslation | null }>(
      `/api/chat/messages/${messageId}/translate`,
      { method: "POST" }
    ),

  translateIdeaMessage: (messageId: number) =>
    request<{ translation: MessageTranslation | null }>(
      `/api/ideas/messages/${messageId}/translate`,
      { method: "POST" }
    ),

  deleteChatMessage: (messageId: number) =>
    request<{ ok: boolean; message: ChatMessage }>(
      `/api/chat/messages/${messageId}`,
      { method: "DELETE" }
    ),

  // --- the ideas workspace ---
  ideas: (params: {
    status?: string;
    tag?: string;
    sort?: string;
    q?: string;
  } = {}) => {
    const qs = new URLSearchParams(
      Object.entries(params).filter(([, v]) => v) as [string, string][]
    ).toString();
    return request<IdeaBoard>(`/api/ideas${qs ? `?${qs}` : ""}`);
  },

  createIdea: (input: {
    title: string;
    description?: string;
    tags?: string[];
    /** "under_review" posts straight into review instead of as a draft. */
    status?: IdeaStatus;
  }) =>
    request<{ idea: Idea }>("/api/ideas", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  idea: (id: number) => request<IdeaDetail>(`/api/ideas/${id}`),

  updateIdea: (
    id: number,
    input: {
      title?: string;
      description?: string;
      status?: IdeaStatus;
      tags?: string[];
      /** Recorded in the history alongside a status change. */
      note?: string;
    }
  ) =>
    request<{ idea: Idea }>(`/api/ideas/${id}`, {
      method: "PATCH",
      body: JSON.stringify(input),
    }),

  deleteIdea: (id: number) =>
    request<{ ok: boolean }>(`/api/ideas/${id}`, { method: "DELETE" }),

  voteIdea: (id: number, value: number) =>
    request<{ idea: Idea }>(`/api/ideas/${id}/vote`, {
      method: "POST",
      body: JSON.stringify({ value }),
    }),

  starIdea: (id: number) =>
    request<{ starred: boolean }>(`/api/ideas/${id}/star`, { method: "POST" }),

  archiveIdea: (id: number, reason: string) =>
    request<{ idea: Idea }>(`/api/ideas/${id}/archive`, {
      method: "POST",
      body: JSON.stringify({ reason }),
    }),

  restoreIdea: (id: number) =>
    request<{ idea: Idea }>(`/api/ideas/${id}/restore`, { method: "POST" }),

  mergeIdea: (id: number, targetId: number) =>
    request<{ idea: Idea; target: Idea }>(`/api/ideas/${id}/merge`, {
      method: "POST",
      body: JSON.stringify({ targetId }),
    }),

  /** Live duplicate check while composing — deliberately conservative. */
  similarIdeas: (q: string) =>
    request<{ similar: SimilarIdea[] }>(
      `/api/ideas/similar?q=${encodeURIComponent(q)}`
    ),

  ideaSummary: (id: number) => request<ThreadSummary>(`/api/ideas/${id}/summary`),

  createIdeaTask: (id: number, input: { title?: string; sprint?: string }) =>
    request<{ task: IdeaTask; idea: Idea }>(`/api/ideas/${id}/tasks`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  updateIdeaTask: (taskId: number, input: { status?: string; sprint?: string }) =>
    request<{ task: IdeaTask }>(`/api/ideas/tasks/${taskId}`, {
      method: "PATCH",
      body: JSON.stringify(input),
    }),

  // --- the thread on an idea ---
  ideaMessages: async (id: number) => {
    const thread = await request<IdeaThread>(`/api/ideas/${id}/messages`);
    return { ...thread, messages: (thread.messages || []).map(normaliseMessage) };
  },

  postIdeaMessage: async (
    id: number,
    input: { body: string; parentId?: number | null; kind?: string }
  ) => {
    const r = await request<{ message: IdeaMessage }>(
      `/api/ideas/${id}/messages`,
      { method: "POST", body: JSON.stringify(input) }
    );
    return { message: normaliseMessage(r.message) };
  },

  postIdeaAttachment: async (
    id: number,
    file: File,
    opts: {
      body?: string;
      parentId?: number | null;
      kind: "image" | "voice";
      duration?: number;
    }
  ) => {
    const r = await upload<{ message: IdeaMessage }>(
      `/api/ideas/${id}/messages`,
      {
        file,
        body: opts.body ?? "",
        kind: opts.kind,
        parentId: opts.parentId ? String(opts.parentId) : "",
        duration: opts.duration ? String(opts.duration) : "",
      }
    );
    return { message: normaliseMessage(r.message) };
  },

  editIdeaMessage: async (messageId: number, body: string) => {
    const r = await request<{ message: IdeaMessage }>(
      `/api/ideas/messages/${messageId}`,
      { method: "PATCH", body: JSON.stringify({ body }) }
    );
    return { message: normaliseMessage(r.message) };
  },

  deleteIdeaMessage: async (messageId: number) => {
    const r = await request<{ ok: boolean; message: IdeaMessage }>(
      `/api/ideas/messages/${messageId}`,
      { method: "DELETE" }
    );
    return { ...r, message: normaliseMessage(r.message) };
  },

  reactToIdeaMessage: async (messageId: number, emoji: string) => {
    const r = await request<{ message: IdeaMessage }>(
      `/api/ideas/messages/${messageId}/react`,
      { method: "POST", body: JSON.stringify({ emoji }) }
    );
    return { message: normaliseMessage(r.message) };
  },

  startBrainstorm: (id: number, minutes: number, topic: string) =>
    request<{ brainstorm: BrainstormSession }>(`/api/ideas/${id}/brainstorm`, {
      method: "POST",
      body: JSON.stringify({ minutes, topic }),
    }),

  stopBrainstorm: (id: number) =>
    request<{ brainstorm: null }>(`/api/ideas/${id}/brainstorm`, {
      method: "DELETE",
    }),
};

export { ApiError };
