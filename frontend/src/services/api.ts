import { useAuthStore } from "../store/useAuthStore";
import type { ApiErrorPayload, ApiResponse } from "../types/api";

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? "/api";

export class ApiError extends Error {
  status: number;
  data?: unknown;

  constructor(message: string, status: number, data?: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
  }
}

function buildUrl(path: string) {
  if (path.startsWith("http")) {
    return path;
  }

  if (path.startsWith(apiBaseUrl)) {
    return path;
  }

  return `${apiBaseUrl}${path}`;
}

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {}
): Promise<T> {
  const token = useAuthStore.getState().token;
  const headers = new Headers(init.headers);

  headers.set("Accept", "application/json");

  if (!(init.body instanceof FormData) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(buildUrl(path), {
    ...init,
    headers
  });

  const contentType = response.headers.get("content-type") ?? "";
  let payload: ApiResponse<T> | ApiErrorPayload | null = null;

  try {
    if (contentType.includes("application/json")) {
      payload = (await response.json()) as ApiResponse<T> | ApiErrorPayload;
    } else {
      const text = await response.text();
      payload = text ? { message: text } : null;
    }
  } catch {
    payload = null;
  }

  if (!response.ok) {
    if (response.status === 401) {
      useAuthStore.getState().logout();
    }

    const message =
      payload?.message ||
      response.statusText ||
      `Request failed with status ${response.status}`;

    throw new ApiError(message, response.status, payload?.data);
  }

  if (!payload) {
    throw new ApiError("Invalid API response", response.status);
  }

  return (payload as ApiResponse<T>).data;
}
