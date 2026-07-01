import type {
  User,
  Lesson,
  Skill,
  UserQuest,
  CharacterWithFriendship,
  LeaderRow,
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
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
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

  const res = await fetch(`${API_URL}${path}`, { ...options, headers });

  if (res.status === 401) {
    clearSession();
    throw new ApiError("unauthorized", 401);
  }
  if (!res.ok) {
    let msg = `request failed (${res.status})`;
    try {
      const body = await res.json();
      if (body?.error) msg = body.error;
    } catch {
      /* ignore */
    }
    throw new ApiError(msg, res.status);
  }
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
    const res = await fetch(`${API_URL}/api/auth/avatar`, {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: fd,
    });
    if (!res.ok) {
      let msg = "upload failed";
      try {
        msg = (await res.json())?.error || msg;
      } catch {
        /* ignore */
      }
      throw new ApiError(msg, res.status);
    }
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
};

export { ApiError };
