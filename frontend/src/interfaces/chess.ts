interface Piece {
    type: string;
    color: "white" | "black";
}

interface Position {
    board: Piece[][];
}

type Square = string

interface Move {
    from: Square;
    to: Square;
}

export type { Position, Move, Square };