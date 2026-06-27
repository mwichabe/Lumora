"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import {
  api,
  ApiError,
  getToken,
  setToken,
  clearSession,
  getStoredUser,
  setStoredUser,
} from "./api";
import type { User } from "./types";

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  setUser: (u: User) => void;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
  refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUserInternal] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  // Single source of truth: every user update also writes through to the cache,
  // so a reload can rehydrate immediately.
  const setUser = useCallback((u: User | null) => {
    setUserInternal(u);
    setStoredUser(u);
  }, []);

  const refresh = useCallback(async () => {
    if (!getToken()) {
      setUser(null);
      setLoading(false);
      return;
    }
    try {
      const { user } = await api.me();
      setUser(user);
    } catch (e) {
      // Only a genuine 401 means the session is invalid — clear it. Any other
      // failure (server down, offline, timeout) keeps the cached user so the
      // app stays usable and never bounces the user to onboarding.
      if (e instanceof ApiError && e.status === 401) {
        clearSession();
        setUserInternal(null);
      }
    } finally {
      setLoading(false);
    }
  }, [setUser]);

  // Rehydrate on mount. We seed from the cached user first (instant, no flash),
  // then revalidate against the server in the background.
  useEffect(() => {
    const cached = getStoredUser();
    if (cached) setUserInternal(cached);
    refresh();
  }, [refresh]);

  const login = async (email: string, password: string) => {
    const { token, user } = await api.login(email, password);
    setToken(token);
    setUser(user);
  };

  const register = async (email: string, password: string, name: string) => {
    const { token, user } = await api.register(email, password, name);
    setToken(token);
    setUser(user);
  };

  const logout = () => {
    clearSession();
    setUserInternal(null);
    router.push("/");
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        setUser,
        login,
        register,
        logout,
        refresh,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
