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

export type { Board, PieceType, Color, Square };