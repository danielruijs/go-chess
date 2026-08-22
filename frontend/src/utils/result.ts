import type { Color } from "../types/chess";
import type { PlayerResult, Result } from "../types/result";

function getPlayerResult(result: Result, playedColor: Color): PlayerResult {
  if (result.outcome === "draw") {
    return "draw";
  }

  const isWin =
    (playedColor === "white" && result.outcome === "white_win") ||
    (playedColor === "black" && result.outcome === "black_win");

  return isWin ? "win" : "loss";
}

export { getPlayerResult };
