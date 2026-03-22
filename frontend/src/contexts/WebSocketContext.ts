import { createContext, useContext } from "react";
import type { WSMessage, LegalMove, QueueData, ClockData } from "../types/message";
import type { Board, Color } from "../types/chess";
import type { Result } from "../types/result";

export interface WebSocketContextType {
  isConnected: boolean;
  sendMessage: (message: WSMessage) => void;
  board: Board | null;
  legalMoves: Record<string, LegalMove[]> | null;
  pgn: string;
  activeColor: Color | null;
  clock: ClockData | null;
  playerColor: Color | null;
  whitePlayerName: string;
  blackPlayerName: string;
  queues: QueueData[] | null;
  inMatch: boolean;
  matchResult: Result | null;
  isDrawOfferPending: boolean;
  isDrawDeclinedNoticeOpen: boolean;
  respondToDrawOffer: (accept: boolean) => void;
  closeDrawDeclinedNotice: () => void;
}

export const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
}
