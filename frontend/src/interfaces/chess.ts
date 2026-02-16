interface Piece {
    type: string;
    color: "white" | "black";
}

type Board = Piece[][];

type Square = string

interface Move {
    from: Square;
    to: Square;
}

export type { Board, Move, Square };