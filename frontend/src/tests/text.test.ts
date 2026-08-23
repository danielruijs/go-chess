import { pluralize } from "../utils/text";

describe("pluralize", () => {
  it("uses singular form when count is 1", () => {
    expect(pluralize(1, "game", "games", true)).toBe("1 game");
    expect(pluralize(1, "move", "moves", true)).toBe("1 move");
  });

  it("uses plural form when count is 0", () => {
    expect(pluralize(0, "game", "games", true)).toBe("0 games");
    expect(pluralize(0, "move", "moves", true)).toBe("0 moves");
  });

  it("uses plural form when count is greater than 1", () => {
    expect(pluralize(2, "game", "games", true)).toBe("2 games");
    expect(pluralize(42, "move", "moves", true)).toBe("42 moves");
  });

  it("omits count when includeCount is false", () => {
    expect(pluralize(1, "Win", "Wins", false)).toBe("Win");
    expect(pluralize(0, "Win", "Wins", false)).toBe("Wins");
    expect(pluralize(3, "Win", "Wins", false)).toBe("Wins");
    expect(pluralize(1, "Loss", "Losses", false)).toBe("Loss");
    expect(pluralize(5, "Loss", "Losses", false)).toBe("Losses");
  });
});
