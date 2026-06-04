import type { Board, Color, PieceType, TimeFormat } from "./chess";
import type { Result } from "./result";

// inbound
const MessageTypeBoard = "board";
const MessageTypeMatchmakingUpdate = "matchmaking_update";
const MessageTypeStartMatch = "start_match";
const MessageTypeEndMatch = "end_match";
const MessageTypeDrawOffered = "draw_offered";
const MessageTypeDrawDeclined = "draw_declined";
const MessageTypePlayerInfo = "player_info";
// outbound
const MessageTypeJoinMatch = "join_match";
const MessageTypeLeaveMatch = "leave_match";
const MessageTypeMove = "move";
const MessageTypeResign = "resign";
const MessageTypeOfferDraw = "offer_draw";
const MessageTypeRespondDraw = "respond_draw";

// inbound
type LegalMove = {
  to: string;
  promotion: PieceType | null;
};

type ClockData = {
  whiteTimeMs: number;
  blackTimeMs: number;
  incrementMs: number;
};

type BoardData = {
  board: Board;
  legalMoves: Record<string, LegalMove[]>;
  pgn: string;
  activeColor: Color;
  clock: ClockData;
};

type QueueData = {
  timeFormat: TimeFormat;
  queueLength: number;
  inQueue: boolean;
};

type MatchmakingUpdateData = {
  queues: QueueData[];
};

type StartMatchData = {
  color: Color;
  whitePlayerName: string;
  blackPlayerName: string;
  clock: ClockData;
};

type EndMatchData = {
  result: Result;
};

type PlayerInfoData = {
  displayName: string;
  username: string;
  isAuthenticated: boolean;
};

// outbound
type JoinMatchData = {
  timeFormat: TimeFormat;
};

type LeaveMatchData = {
  timeFormat: TimeFormat;
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
  | {
      type: typeof MessageTypeDrawOffered;
    }
  | {
      type: typeof MessageTypeDrawDeclined;
    }
  | {
      type: typeof MessageTypePlayerInfo;
      data: PlayerInfoData;
    }
  // outbound
  | {
      type: typeof MessageTypeJoinMatch;
      data: JoinMatchData;
    }
  | {
      type: typeof MessageTypeLeaveMatch;
      data: LeaveMatchData;
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

export type {
  WSMessage,
  MoveData,
  BoardData,
  JoinMatchData,
  LegalMove,
  StartMatchData,
  MatchmakingUpdateData,
  EndMatchData,
  ClockData,
  QueueData,
  PlayerInfoData,
};
export {
  MessageTypeBoard,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeLeaveMatch,
  MessageTypeStartMatch,
  MessageTypeMatchmakingUpdate,
  MessageTypeEndMatch,
  MessageTypeDrawOffered,
  MessageTypeDrawDeclined,
  MessageTypeResign,
  MessageTypeOfferDraw,
  MessageTypeRespondDraw,
  MessageTypePlayerInfo,
};
