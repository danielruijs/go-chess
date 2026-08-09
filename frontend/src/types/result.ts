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

const OUTCOME_TEXT: Record<Outcome, string> = {
  white_win: "White wins",
  black_win: "Black wins",
  draw: "Draw",
} as const;

const REASON_TEXT: Record<Reason, string> = {
  checkmate: "by Checkmate",
  stalemate: "by Stalemate",
  threefold_repetition: "by Threefold Repetition",
  fifty_moves_rule: "by Fifty-Move Rule",
  insufficient_material: "by Insufficient Material",
  resignation: "by Resignation",
  agreed_draw: "by Agreement",
  timeout: "by Timeout",
} as const;

export type { Result };
export { OUTCOME_TEXT, REASON_TEXT };
