import { coordsToSquare } from "../utils/chess";

describe("coordsToSquare", () => {
    it("converts (0,0) to a8", () => {
        expect(coordsToSquare(0, 0)).toBe("a8");
    });
    it("converts (7,7) to h1", () => {
        expect(coordsToSquare(7, 7)).toBe("h1");
    });
    it("converts (5,2) to c3", () => {
        expect(coordsToSquare(5, 2)).toBe("c3");
    });
});