import type { Square, Color, TimeFormat } from "../interfaces/chess";
import type { QueueData } from "../interfaces/message";

function coordsToString(square: Square): string {
  const fileChar = String.fromCharCode("a".charCodeAt(0) + square.file);
  const rankChar = (square.rank + 1).toString();
  return fileChar + rankChar;
}

function pgnToMoves(pgn: string): string[] {
  const moves: string[] = [];
  const tokens = pgn.split(/\s+/); // split on whitespace
  for (const token of tokens) {
    if (token.trim() === "") continue; // skip empty tokens
    if (token.includes("-") || token.includes("*")) continue; // skip result tokens
    if (token.includes(".")) {
      const parts = token.split(".");
      if (parts.length === 2) {
        moves.push(parts[1]);
      }
    } else {
      moves.push(token);
    }
  }
  return moves;
}

function displayIndexToSquare(index: number, color: Color): Square {
  const displayFile = index % 8;
  const displayRow = Math.floor(index / 8);
  const file = color === "white" ? displayFile : 7 - displayFile;
  const rank = color === "white" ? 7 - displayRow : displayRow;
  return { file, rank };
}

function getQueueData(queues: QueueData[] | null, timeFormat: TimeFormat): QueueData | undefined {
  return queues?.find(
    (q) =>
      q.timeFormat.initialMs === timeFormat.initialMs &&
      q.timeFormat.incrementMs === timeFormat.incrementMs
  );
}

export { coordsToString, pgnToMoves, displayIndexToSquare, getQueueData };
