import type { Position, Move } from "./chess";

type WSMessage =
  | {
      type: typeof MessageTypePosition;
      data: Position
    }
  | {
      type: typeof MessageTypeJoinMatch;
      data: string
    }
    | {
        type: typeof MessageTypeMovePiece
        data: Move
    };

const MessageTypePosition = "position";
const MessageTypeMovePiece = "move";
const MessageTypeJoinMatch = "join_match";

export type { WSMessage };
export { MessageTypePosition, MessageTypeMovePiece, MessageTypeJoinMatch };