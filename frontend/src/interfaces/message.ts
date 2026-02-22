import type { Board, PieceType } from "./chess";

const MessageTypeBoard = "board";
const MessageTypeMove = "move";
const MessageTypeJoinMatch = "join_match";
const MessageTypeError = "error";

type LegalMove = {
  to: string;
  promotion: PieceType | null;
};

type BoardData = {
  board: Board;
  legalMoves: Record<string, LegalMove[]>;
};

type MoveData = {
  from: string;
  to: string;
  promotion: PieceType | null;
};

type JoinMatchData = {
  playerName: string;
};

type ErrorData = {
  message: string;
};

type WSMessage =
  | {
    type: typeof MessageTypeBoard;
    data: BoardData;
  }
  | {
    type: typeof MessageTypeMove;
    data: MoveData;
  }
  | {
    type: typeof MessageTypeJoinMatch;
    data: JoinMatchData;
  }
  | {
    type: typeof MessageTypeError;
    data: ErrorData;
  };

export type { WSMessage, MoveData, BoardData, JoinMatchData, ErrorData, LegalMove };
export {
  MessageTypeBoard,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeError,
};
