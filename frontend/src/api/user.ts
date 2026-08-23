import type { UserProfile } from "../types/user";

const API_URL: string = import.meta.env.VITE_API_URL as string;

export type UserProfileResult = { status: "ok"; data: UserProfile } | { status: "not-found" };

/**
 * Fetches public profile, stats, and matches for a given username.
 */
export async function fetchUserProfile(username: string): Promise<UserProfileResult> {
  const response = await fetch(`${API_URL}/api/users/${encodeURIComponent(username)}`, {
    method: "GET",
    credentials: "include",
  });

  if (response.status === 404) return { status: "not-found" };

  if (!response.ok) {
    const errorText = (await response.text()).trim();
    throw new Error(errorText || `Failed to fetch user profile: ${response.status}`);
  }

  const data = (await response.json()) as UserProfile;
  return { status: "ok", data };
}
