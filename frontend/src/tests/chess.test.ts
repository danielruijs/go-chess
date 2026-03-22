import type { Board } from "../types/chess";
import type { QueueData } from "../types/message";
import {
  coordsToString,
  pgnToMoves,
  displayIndexToSquare,
  getQueueData,
  getMaterialDiff,
} from "../utils/chess";

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

describe("pgnToMoves", () => {
  it("parses simple PGN", () => {
    const pgn = "1.d4 d5 2.c4 c5 3.b4 b5 4.bxc5 Bd7 5.c6 Bxc6 ";
    const moves = pgnToMoves(pgn);
    expect(moves).toEqual(["d4", "d5", "c4", "c5", "b4", "b5", "bxc5", "Bd7", "c6", "Bxc6"]);
  });
  it("parses PGN with result", () => {
    const pgn = "1.d4 e5 0-1";
    const moves = pgnToMoves(pgn);
    expect(moves).toEqual(["d4", "e5"]);
  });
});

describe("displayIndexToSquare", () => {
  it("index 0 white -> a8", () => {
    const sq = displayIndexToSquare(0, "white");
    expect(coordsToString(sq)).toBe("a8");
  });

  it("index 1 white -> b8", () => {
    const sq = displayIndexToSquare(1, "white");
    expect(coordsToString(sq)).toBe("b8");
  });

  it("index 63 white -> h1", () => {
    const sq = displayIndexToSquare(63, "white");
    expect(coordsToString(sq)).toBe("h1");
  });

  it("index 0 black -> h1", () => {
    const sq = displayIndexToSquare(0, "black");
    expect(coordsToString(sq)).toBe("h1");
  });

  it("index 1 black -> g1", () => {
    const sq = displayIndexToSquare(1, "black");
    expect(coordsToString(sq)).toBe("g1");
  });

  it("index 63 black -> a8", () => {
    const sq = displayIndexToSquare(63, "black");
    expect(coordsToString(sq)).toBe("a8");
  });
});

describe("getQueueData", () => {
  it("returns the matching queue data", () => {
    const queues: QueueData[] = [
      { timeFormat: { initialMs: 60000, incrementMs: 0 } } as QueueData,
      { timeFormat: { initialMs: 300000, incrementMs: 3000 } } as QueueData,
    ];
    const result = getQueueData(queues, {
      initialMs: 300000,
      incrementMs: 3000,
    });
    expect(result).toEqual(queues[1]);
  });

  it("returns undefined when no queue matches", () => {
    const queues: QueueData[] = [{ timeFormat: { initialMs: 60000, incrementMs: 0 } } as QueueData];
    const result = getQueueData(queues, {
      initialMs: 300000,
      incrementMs: 3000,
    });
    expect(result).toBeUndefined();
  });

  it("returns undefined when queues is null", () => {
    const result = getQueueData(null, { initialMs: 60000, incrementMs: 0 });
    expect(result).toBeUndefined();
  });
});

describe("getMaterialDiff", () => {
  it("returns zero score for empty board", () => {
    const board: Board = Array(8)
      .fill(null)
      .map(() => Array(8).fill(null));
    const result = getMaterialDiff(board);
    expect(result.score).toBe(0);
    expect(result.extraPieces.white).toEqual({});
    expect(result.extraPieces.black).toEqual({});
  });

  it("calculates score when white has extra pawn", () => {
    const board: Board = Array(8)
      .fill(null)
      .map(() => Array(8).fill(null));
    board[0][0] = { color: "white", type: "pawn" };
    const result = getMaterialDiff(board);
    expect(result.score).toBe(1);
    expect(result.extraPieces.white).toEqual({ pawn: 1 });
    expect(result.extraPieces.black).toEqual({});
  });

  it("calculates score when black has extra queen", () => {
    const board: Board = Array(8)
      .fill(null)
      .map(() => Array(8).fill(null));
    board[0][0] = { color: "black", type: "queen" };
    const result = getMaterialDiff(board);
    expect(result.score).toBe(-9);
    expect(result.extraPieces.white).toEqual({});
    expect(result.extraPieces.black).toEqual({ queen: 1 });
  });

  it("calculates combined material difference", () => {
    const board: Board = Array(8)
      .fill(null)
      .map(() => Array(8).fill(null));
    board[0][0] = { color: "white", type: "rook" };
    board[1][0] = { color: "white", type: "knight" };
    board[2][0] = { color: "black", type: "bishop" };
    const result = getMaterialDiff(board);
    expect(result.score).toBe(5); // 5 + 3 - 3
    expect(result.extraPieces.white).toEqual({ rook: 1, knight: 1 });
    expect(result.extraPieces.black).toEqual({ bishop: 1 });
  });

  it("handles multiple pieces of same type", () => {
    const board: Board = Array(8)
      .fill(null)
      .map(() => Array(8).fill(null));
    board[0][0] = { color: "white", type: "pawn" };
    board[0][1] = { color: "white", type: "pawn" };
    board[1][0] = { color: "black", type: "pawn" };
    const result = getMaterialDiff(board);
    expect(result.score).toBe(1);
    expect(result.extraPieces.white).toEqual({ pawn: 1 });
  });

  it("handles material difference but different pieces", () => {
    const board: Board = Array(8)
      .fill(null)
      .map(() => Array(8).fill(null));
    board[1][0] = { color: "white", type: "knight" };
    board[2][0] = { color: "black", type: "bishop" };
    const result = getMaterialDiff(board);
    expect(result.score).toBe(0);
    expect(result.extraPieces.white).toEqual({ knight: 1 });
    expect(result.extraPieces.black).toEqual({ bishop: 1 });
  });

  it("handles null board", () => {
    const result = getMaterialDiff(null);
    expect(result.score).toBe(0);
    expect(result.extraPieces.white).toEqual({});
    expect(result.extraPieces.black).toEqual({});
  });
});
