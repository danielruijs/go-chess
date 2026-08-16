import {
  type Square,
  type Color,
  type TimeFormat,
  type Board,
  type PieceType,
  type MaterialDiff,
  pieceValues,
} from "../types/chess";
import type { QueueData } from "../types/message";

function coordsToString(square: Square): string {
  const fileChar = String.fromCharCode("a".charCodeAt(0) + square.file);
  const rankChar = (square.rank + 1).toString();
  return fileChar + rankChar;
}

function pgnToMoves(pgn: string): string[] {
  const moves: string[] = [];
  const tokens = pgn.trim().split(/\s+/);
  for (const token of tokens) {
    // Remove move numbers
    const moveCandidate = token.replace(/^\d+\./, "").trim();
    if (!moveCandidate) continue;

    if (
      moveCandidate === "1-0" ||
      moveCandidate === "0-1" ||
      moveCandidate === "1/2-1/2" ||
      moveCandidate === "*"
    ) {
      continue;
    }

    moves.push(moveCandidate);
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

function getMaterialDiff(board: Board | null): MaterialDiff {
  const materialDiff: MaterialDiff = {
    score: 0,
    extraPieces: { white: {}, black: {} },
  };
  if (!board) {
    return materialDiff;
  }

  const counts: Record<Color, Record<PieceType, number>> = {
    white: { pawn: 0, knight: 0, bishop: 0, rook: 0, queen: 0, king: 0 },
    black: { pawn: 0, knight: 0, bishop: 0, rook: 0, queen: 0, king: 0 },
  };
  board.flat().forEach((piece) => {
    if (piece) counts[piece.color][piece.type]++;
  });

  for (const [type, value] of Object.entries(pieceValues) as [PieceType, number][]) {
    const diff = counts.white[type] - counts.black[type];

    materialDiff.score += diff * value;

    if (diff > 0) materialDiff.extraPieces.white[type] = diff;
    else if (diff < 0) materialDiff.extraPieces.black[type] = Math.abs(diff);
  }

  return materialDiff;
}

type MoveItem = {
  san: string;
  positionIndex: number;
};

type MoveRow = {
  moveNumber: number;
  white: MoveItem;
  black: MoveItem | null;
};

function buildMoveRows(moves: MoveItem[]): MoveRow[] {
  const rows: MoveRow[] = [];
  for (let i = 0; i < moves.length; i += 2) {
    rows.push({
      moveNumber: Math.floor(i / 2) + 1,
      white: moves[i],
      black: moves[i + 1] ?? null,
    });
  }
  return rows;
}

export {
  coordsToString,
  pgnToMoves,
  displayIndexToSquare,
  getQueueData,
  getMaterialDiff,
  buildMoveRows,
};

export type { MoveItem, MoveRow };
