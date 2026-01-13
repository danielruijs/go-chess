import {type Square } from "../interfaces/chess";

function coordsToSquare(row: number, col: number): Square {
    const file = String.fromCharCode("a".charCodeAt(0) + col);
    const rank = (8 - row).toString();
    return file + rank;
}

export { coordsToSquare };