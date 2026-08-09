import type { Board } from "./chess";
import type { Result } from "./result";

type Position = {
  index: number;
  board: Board;
  san?: string; // empty/omitted for starting position (index 0)
  whiteTimeMs: number;
  blackTimeMs: number;
};

type Match = {
  whitePlayerName: string;
  blackPlayerName: string;
  result: Result;
  positions: Position[];
};

export type { Match };
