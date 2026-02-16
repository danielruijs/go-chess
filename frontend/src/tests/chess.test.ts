import { coordsToSquare } from "../utils/chess";

describe("coordsToSquare", () => {
    it("converts (0,0) to a1", () => {
        expect(coordsToSquare(0, 0)).toBe("a1");
    });
    it("converts (7,7) to h8", () => {
        expect(coordsToSquare(7, 7)).toBe("h8");
    });
    it("converts (5,2) to c3", () => {
        expect(coordsToSquare(5, 2)).toBe("f3");
    });
});