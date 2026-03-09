type PieceType = "pawn" | "knight" | "bishop" | "rook" | "queen" | "king";
type Color = "white" | "black";

interface Piece {
    type: PieceType;
    color: Color;
}

type Board = (Piece | null)[][];

/**
 * (0,0) is a1, (1,0) is b1, (7,7) is h8
 */
type Square = { file: number; rank: number };

type TimeFormat = {
    initialMs: number;
    incrementMs: number;
};

const timeFormats: TimeFormat[] = [
    { initialMs: 1 * 60 * 1000, incrementMs: 0 }, // 1|0
    { initialMs: 1 * 60 * 1000, incrementMs: 1 * 1000 }, // 1|1
    { initialMs: 3 * 60 * 1000, incrementMs: 0 }, // 3|0
    { initialMs: 3 * 60 * 1000, incrementMs: 2 * 1000 }, // 3|2
    { initialMs: 10 * 60 * 1000, incrementMs: 0 }, // 10|0
    { initialMs: 15 * 60 * 1000, incrementMs: 10 * 1000 }, // 15|10
];

export type { Board, PieceType, Color, Square, TimeFormat };
export { timeFormats };