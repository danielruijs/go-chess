type PieceType = "pawn" | "knight" | "bishop" | "rook" | "queen" | "king";
type Color = "white" | "black";

interface Piece {
  type: PieceType;
  color: Color;
}

type Board = (Piece | null)[][];

const pieceValues: Record<PieceType, number> = {
  pawn: 1,
  knight: 3,
  bishop: 3,
  rook: 5,
  queen: 9,
  king: 1000,
};

type MaterialDiff = {
  extraPieces: {
    white: Partial<Record<PieceType, number>>;
    black: Partial<Record<PieceType, number>>;
  };
  score: number; // positive if white has more, negative if black has more
};

/**
 * (0,0) is a1, (1,0) is b1, (7,7) is h8
 */
type Square = { file: number; rank: number };

type TimeFormat = {
  initialMs: number;
  incrementMs: number;
};

const timeFormats: TimeFormat[][] = [
  [
    { initialMs: 1 * 60 * 1000, incrementMs: 0 }, // 1|0
    { initialMs: 1 * 60 * 1000, incrementMs: 1 * 1000 }, // 1|1
  ],
  [
    { initialMs: 3 * 60 * 1000, incrementMs: 0 }, // 3|0
    { initialMs: 3 * 60 * 1000, incrementMs: 2 * 1000 }, // 3|2
    { initialMs: 5 * 60 * 1000, incrementMs: 5 * 1000 }, // 5|5
  ],
  [
    { initialMs: 10 * 60 * 1000, incrementMs: 0 }, // 10|0
    { initialMs: 15 * 60 * 1000, incrementMs: 10 * 1000 }, // 15|10
  ],
];

export type { Board, PieceType, Color, Square, TimeFormat, MaterialDiff };
export { timeFormats, pieceValues };
