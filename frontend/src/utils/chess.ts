import { type Square } from "../interfaces/chess";

function coordsToSquare(rank: number, file: number): Square {
    const fileChar = String.fromCharCode("a".charCodeAt(0) + file);
    const rankChar = (8 - rank).toString();
    return fileChar + rankChar;
}

export { coordsToSquare };