import type {
  User,
  Lesson,
  Skill,
  UserQuest,
  CharacterWithFriendship,
  LeaderRow,
  HomeData,
  ListeningSession,
} from "./types";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

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

  quests: () => request<{ quests: UserQuest[] }>("/api/quests/daily"),

  characters: () =>
    request<{ characters: CharacterWithFriendship[] }>("/api/characters"),

  leaderboard: () =>
    request<{ league: string; rows: LeaderRow[] }>("/api/leaderboard"),
};

export { ApiError };
