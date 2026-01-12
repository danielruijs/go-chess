interface Piece {
    type: string;
    color: "white" | "black";
}

interface Position {
    board: Piece[][];
}

export type { Position };