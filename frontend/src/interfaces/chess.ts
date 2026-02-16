type PieceType = "pawn" | "knight" | "bishop" | "rook" | "queen" | "king";
type Color = "white" | "black";

interface Piece {
    type: PieceType;
    color: Color;
}

type Board = Piece[][];

export type { Board, PieceType };