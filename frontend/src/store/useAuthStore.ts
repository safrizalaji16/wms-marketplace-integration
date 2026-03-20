import { create } from "zustand";

const authStorageKey = "wms-auth";

export interface AuthState {
  isAuthenticated: boolean;
  userName: string;
  token: string;
  role: string;
  login: (payload: { userName: string; token: string }) => void;
  logout: () => void;
}

function parseJwtRole(token: string) {
  try {
    const [, payload] = token.split(".");

    if (!payload) {
      return "";
    }

    const decodedPayload = JSON.parse(atob(payload)) as { role?: string };
    return decodedPayload.role ?? "";
  } catch {
    return "";
  }
}

function getInitialAuthState() {
  if (typeof window === "undefined") {
    return {
      isAuthenticated: false,
      userName: "",
      token: "",
      role: ""
    };
  }

  const storedValue = window.localStorage.getItem(authStorageKey);

  if (!storedValue) {
    return {
      isAuthenticated: false,
      userName: "",
      token: "",
      role: ""
    };
  }

  try {
    const parsed = JSON.parse(storedValue) as {
      userName?: string;
      token?: string;
      role?: string;
    };

    return {
      isAuthenticated: Boolean(parsed.token),
      userName: parsed.userName ?? "",
      token: parsed.token ?? "",
      role: parsed.role ?? parseJwtRole(parsed.token ?? "")
    };
  } catch {
    return {
      isAuthenticated: false,
      userName: "",
      token: "",
      role: ""
    };
  }
}

export const useAuthStore = create<AuthState>((set) => ({
  ...getInitialAuthState(),
  login: ({ userName, token }) => {
    const role = parseJwtRole(token);

    if (typeof window !== "undefined") {
      window.localStorage.setItem(
        authStorageKey,
        JSON.stringify({ userName, token, role })
      );
    }

    set({ isAuthenticated: true, userName, token, role });
  },
  logout: () => {
    if (typeof window !== "undefined") {
      window.localStorage.removeItem(authStorageKey);
    }

    set({ isAuthenticated: false, userName: "", token: "", role: "" });
  }
}));
