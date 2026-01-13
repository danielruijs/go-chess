import type { Position } from "./chess";

export type WSMessage =
  | {
      type: typeof MessageTypePosition;
      data: Position;
    }
  | {
      type: typeof MessageTypeJoinMatch;
      data: {};
    };

export const MessageTypePosition = "position";
export const MessageTypeJoinMatch = "join_match";