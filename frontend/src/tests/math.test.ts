import { getPct } from "../utils/math";

describe("getPct", () => {
  it("calculates percentage correctly", () => {
    expect(getPct(4, 10)).toBe(40);
    expect(getPct(1, 3)).toBeCloseTo(33.333, 3);
  });

  it("returns 0 when denominator is 0 without NaN", () => {
    expect(getPct(0, 0)).toBe(0);
    expect(getPct(5, 0)).toBe(0);
  });
});
