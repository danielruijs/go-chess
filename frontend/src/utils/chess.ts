import { type Square } from "../interfaces/chess";

function coordsToSquare(file: number, rank: number): Square {
    const fileChar = String.fromCharCode("a".charCodeAt(0) + file);
    const rankChar = (rank + 1).toString();
    return fileChar + rankChar;
}

export { coordsToSquare };