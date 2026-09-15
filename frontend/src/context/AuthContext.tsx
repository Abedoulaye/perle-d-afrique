import { createContext, useContext, useState, type ReactNode } from "react";
import type { User } from "../types";
import { apiFetch } from "../services/api";

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>; // ← add Promise<void> here
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => {
    const savedUser = localStorage.getItem("user");
    return savedUser ? JSON.parse(savedUser) : null;
  });

  const [token, setToken] = useState<string | null>(
    localStorage.getItem("token"),
  );

  const login = async (email: string, password: string) => {
    const data = await apiFetch("/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    localStorage.setItem("token", data.access_token);
    localStorage.setItem("refresh_token", data.refresh_token);
    localStorage.setItem("user", JSON.stringify(data.user));
    setToken(data.access_token);
    setUser(data.user);
  };

  const register = async (email: string, password: string) => {
    await apiFetch("/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    await login(email, password);
  };

  const logout = async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (refreshToken) {
      try {
        await apiFetch("/logout", {
          method: "POST",
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      } catch {
        // Ignore errors, just clear local state
      }
    }
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    setToken(null);
    setUser(null);
  };

  const isAdmin = user?.role === "admin";

  return (
    <AuthContext.Provider
      value={{ user, token, login, register, logout, isAdmin }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
