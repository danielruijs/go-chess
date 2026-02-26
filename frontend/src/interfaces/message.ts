import type { Board, PieceType } from "./chess";
import type { Result } from "./result";

const MessageTypeBoard = "board";
const MessageTypeMove = "move";
const MessageTypeJoinMatch = "join_match";
const MessageTypeMatchmakingUpdate = "matchmaking_update";
const MessageTypeStartMatch = "start_match";
const MessageTypeEndMatch = "end_match"

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

type MatchmakingUpdateData = {
  queueLength: number;
  inQueue: boolean;
};

type StartMatchData = {
  whitePlayerName: string;
  blackPlayerName: string;
};

type EndMatchData = {
  result: Result;
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
    type: typeof MessageTypeStartMatch;
    data: StartMatchData;
  }
  | {
    type: typeof MessageTypeMatchmakingUpdate;
    data: MatchmakingUpdateData;
  }
  | {
    type: typeof MessageTypeEndMatch;
    data: EndMatchData;
  };

export type { WSMessage, MoveData, BoardData, JoinMatchData, LegalMove, StartMatchData, MatchmakingUpdateData, EndMatchData };
export {
  MessageTypeBoard,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeStartMatch,
  MessageTypeMatchmakingUpdate,
  MessageTypeEndMatch,
};
