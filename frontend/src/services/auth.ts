import { apiFetch } from "./api";
import type { AuthLoginData, AuthLoginPayload } from "../types/api";

export function login(payload: AuthLoginPayload) {
  return apiFetch<AuthLoginData>("/api/auth/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}
