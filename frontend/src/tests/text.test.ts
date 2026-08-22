import { pluralize } from "../utils/text";

describe("pluralize", () => {
  it("uses singular form when count is 1", () => {
    expect(pluralize(1, "game", "games")).toBe("1 game");
    expect(pluralize(1, "move", "moves")).toBe("1 move");
  });

  it("uses plural form when count is 0", () => {
    expect(pluralize(0, "game", "games")).toBe("0 games");
    expect(pluralize(0, "move", "moves")).toBe("0 moves");
  });

  it("uses plural form when count is greater than 1", () => {
    expect(pluralize(2, "game", "games")).toBe("2 games");
    expect(pluralize(42, "move", "moves")).toBe("42 moves");
  });
});
