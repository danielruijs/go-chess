import type { Board, PieceType } from "./chess";

const MessageTypeBoard = "board";
const MessageTypeMove = "move";
const MessageTypeJoinMatch = "join_match";
const MessageTypeStartMatch = "start_match";

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

type StartMatchData = {
	whitePlayerName: string;
	blackPlayerName: string;
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
  };

export type { WSMessage, MoveData, BoardData, JoinMatchData, LegalMove, StartMatchData };
export {
  MessageTypeBoard,
  MessageTypeMove,
  MessageTypeJoinMatch,
  MessageTypeStartMatch,
};
