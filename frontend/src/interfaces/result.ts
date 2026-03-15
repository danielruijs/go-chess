type Outcome = "white_win" | "black_win" | "draw";

type Reason =
  | "checkmate"
  | "stalemate"
  | "threefold_repetition"
  | "fifty_moves_rule"
  | "insufficient_material"
  | "resignation"
  | "agreed_draw"
  | "timeout";

type Result = {
  outcome: Outcome;
  reason: Reason;
};

export type { Result };
