import type { Match } from "../types/match";

const API_URL: string = import.meta.env.VITE_API_URL as string;

export type MatchResult =
  | { status: "ok"; data: Match }
  | { status: "not-found" }
  | { status: "in-progress" };

/**
 * Fetches the data for a completed match by its public ID.
 * Throws on network or server error.
 */
export async function fetchMatch(publicId: string): Promise<MatchResult> {
  const response = await fetch(`${API_URL}/api/match/${publicId}`, {
    method: "GET",
    credentials: "include",
  });

  if (response.status === 404) return { status: "not-found" };
  if (response.status === 409) return { status: "in-progress" };

  if (!response.ok) {
    throw new Error(`Failed to fetch match analysis: ${response.status}`);
  }

  const data = (await response.json()) as Match;
  return { status: "ok", data };
}
