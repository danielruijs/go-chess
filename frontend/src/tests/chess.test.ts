import { coordsToString, parsePgn } from "../utils/chess";

describe("coordsToString", () => {
    it("converts (0,0) to a1", () => {
        expect(coordsToString({ file: 0, rank: 0 })).toBe("a1");
    });
    it("converts (1,0) to b1", () => {
        expect(coordsToString({ file: 1, rank: 0 })).toBe("b1");
    });
    it("converts (7,7) to h8", () => {
        expect(coordsToString({ file: 7, rank: 7 })).toBe("h8");
    });
    it("converts (5,2) to f3", () => {
        expect(coordsToString({ file: 5, rank: 2 })).toBe("f3");
    });
});

describe("parsePgn", () => {
    it("parses simple PGN", () => {
        const pgn = "1.d4 d5 2.c4 c5 3.b4 b5 4.bxc5 Bd7 5.c6 Bxc6 ";
        const moves = parsePgn(pgn);
        expect(moves).toEqual(["d4", "d5", "c4", "c5", "b4", "b5", "bxc5", "Bd7", "c6", "Bxc6"]);
    });
});