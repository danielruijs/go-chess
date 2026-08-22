import { getPlayerResult } from "../utils/result";

describe("getPlayerResult", () => {
  it("identifies white win as win for white player", () => {
    const outcome = getPlayerResult(
      {
        outcome: "white_win",
        reason: "checkmate",
      },
      "white"
    );
    expect(outcome).toBe("win");
  });

  it("identifies white win as loss for black player", () => {
    const outcome = getPlayerResult(
      {
        outcome: "white_win",
        reason: "checkmate",
      },
      "black"
    );
    expect(outcome).toBe("loss");
  });

  it("identifies black win as win for black player", () => {
    const outcome = getPlayerResult(
      {
        outcome: "black_win",
        reason: "resignation",
      },
      "black"
    );
    expect(outcome).toBe("win");
  });

  it("identifies draw as draw for either player", () => {
    const outcome = getPlayerResult(
      {
        outcome: "draw",
        reason: "stalemate",
      },
      "white"
    );
    expect(outcome).toBe("draw");
  });
});
