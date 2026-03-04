import type { Board, Color, PieceType } from "./chess";
import type { Result } from "./result";

// inbound
const MessageTypeBoard = "board";
const MessageTypeMatchmakingUpdate = "matchmaking_update";
const MessageTypeStartMatch = "start_match";
const MessageTypeEndMatch = "end_match"
// outbound
const MessageTypeJoinMatch = "join_match";
const MessageTypeMove = "move";
const MessageTypeResign = "resign";
const MessageTypeOfferDraw = "offer_draw";
const MessageTypeRespondDraw = "respond_draw";

// inbound
type LegalMove = {
  to: string;
  promotion: PieceType | null;
};

type BoardData = {
  board: Board;
  legalMoves: Record<string, LegalMove[]>;
  pgn: string;
};

type MatchmakingUpdateData = {
  queueLength: number;
  inQueue: boolean;
};

type StartMatchData = {
  color: Color;
  whitePlayerName: string;
  blackPlayerName: string;
};

type EndMatchData = {
  result: Result;
};

// outbound
type JoinMatchData = {
  playerName: string;
};

type MoveData = {
  from: string;
  to: string;
  promotion: PieceType | null;
};

type RespondDrawData = {
  accept: boolean;
};

type WSMessage =
  // inbound
  | {
    type: typeof MessageTypeBoard;
    data: BoardData;
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
  }
  // outbound
  | {
    type: typeof MessageTypeJoinMatch;
    data: JoinMatchData;
  }
  | {
    type: typeof MessageTypeMove;
    data: MoveData;
  }
  | {
    type: typeof MessageTypeResign;
  }
  | {
    type: typeof MessageTypeOfferDraw;
  }
  | {
    type: typeof MessageTypeRespondDraw;
    data: RespondDrawData;
  };

export type { WSMessage, MoveData, BoardData, JoinMatchData, LegalMove, StartMatchData, MatchmakingUpdateData, EndMatchData };
export {
  MessageTypeBoard,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeStartMatch,
  MessageTypeMatchmakingUpdate,
  MessageTypeEndMatch,
  MessageTypeResign,
  MessageTypeOfferDraw,
  MessageTypeRespondDraw
};
