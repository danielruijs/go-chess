import { coordsToString } from "../utils/chess";

describe("coordsToString", () => {
    it("converts (0,0) to a1", () => {
        expect(coordsToString(0, 0)).toBe("a1");
    });
    it("converts (7,7) to h8", () => {
        expect(coordsToString(7, 7)).toBe("h8");
    });
    it("converts (5,2) to f3", () => {
        expect(coordsToString(5, 2)).toBe("f3");
    });
});