import type { Position, Move } from "./chess";

const MessageTypePosition = "position";
const MessageTypeMove = "move";
const MessageTypeJoinMatch = "join_match";
const MessageTypeError = "error";

type PositionData = {
  position: Position;
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
    type: typeof MessageTypePosition;
    data: PositionData;
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

export type { WSMessage, MoveData, PositionData, JoinMatchData, ErrorData };
export {
  MessageTypePosition,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeError,
};
