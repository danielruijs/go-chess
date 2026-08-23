import { getPct } from "./math";
import type { DerivedStats, UserStats } from "../types/user";

function computeDerivedStats(stats: UserStats): DerivedStats {
  const wins = stats.white.wins + stats.black.wins;
  const draws = stats.white.draws + stats.black.draws;
  const losses = stats.white.losses + stats.black.losses;
  const totalGames = wins + draws + losses;

  return {
    wins,
    draws,
    losses,
    totalGames,
    winPct: getPct(wins, totalGames),
    drawPct: getPct(draws, totalGames),
    lossPct: getPct(losses, totalGames),
  };
}

export { computeDerivedStats };
