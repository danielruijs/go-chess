import type { Credentials } from "../types/auth";
import type { PlayerInfoData } from "../types/message";

const API_URL: string = import.meta.env.VITE_API_URL as string;

/** Returns session data if a valid session exists, otherwise null. */
export async function checkSession(): Promise<PlayerInfoData | null> {
  const response = await fetch(`${API_URL}/api/check`, { method: "GET", credentials: "include" });
  if (!response.ok) {
    return null;
  }
  return (await response.json()) as PlayerInfoData;
}

/** Sends logout request. Throws on network error. */
export async function logout(): Promise<void> {
  await fetch(`${API_URL}/api/logout`, { method: "POST", credentials: "include" });
}

/** Sends login request. Returns an Error on failure, or null on success. */
export async function login(credentials: Credentials): Promise<Error | null> {
  try {
    const response = await fetch(`${API_URL}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credentials),
      credentials: "include",
    });
    if (!response.ok) {
      const message = await response.text();
      return new Error(message || "An error occurred during login");
    }
    return null;
  } catch (err) {
    return err instanceof Error ? err : new Error("Failed to connect to the server");
  }
}

/** Sends register request. Returns an Error on failure, or null on success. */
export async function register(credentials: Credentials): Promise<Error | null> {
  try {
    const response = await fetch(`${API_URL}/api/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credentials),
      credentials: "include",
    });
    if (!response.ok) {
      const message = await response.text();
      return new Error(message || "An error occurred during registration");
    }
    return null;
  } catch (err) {
    return err instanceof Error ? err : new Error("Failed to connect to the server");
  }
}
