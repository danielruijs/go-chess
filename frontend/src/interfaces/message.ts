import type { Board, Move } from "./chess";

const MessageTypeBoard = "board";
const MessageTypeMove = "move";
const MessageTypeJoinMatch = "join_match";
const MessageTypeError = "error";

type BoardData = {
  board: Board;
};

type MoveData = {
  move: Move;
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

export type { WSMessage, MoveData, BoardData, JoinMatchData, ErrorData };
export {
  MessageTypeBoard,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeError,
};
