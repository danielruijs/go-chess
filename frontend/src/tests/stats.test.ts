import { computeDerivedStats } from "../utils/stats";
import type { UserStats } from "../types/user";

describe("computeDerivedStats", () => {
  it("computes overall totals and percentages correctly", () => {
    const stats: UserStats = {
      white: { wins: 3, draws: 1, losses: 1 },
      black: { wins: 2, draws: 1, losses: 2 },
    };

    const result = computeDerivedStats(stats);
    expect(result).toEqual({
      wins: 5,
      draws: 2,
      losses: 3,
      totalGames: 10,
      winPct: 50,
      drawPct: 20,
      lossPct: 30,
    });
  });

  it("handles zero games without NaN", () => {
    const stats: UserStats = {
      white: { wins: 0, draws: 0, losses: 0 },
      black: { wins: 0, draws: 0, losses: 0 },
    };

    const result = computeDerivedStats(stats);
    expect(result).toEqual({
      wins: 0,
      draws: 0,
      losses: 0,
      totalGames: 0,
      winPct: 0,
      drawPct: 0,
      lossPct: 0,
    });
  });
});
