import { createContext, useContext } from "react";
import type { WSMessage, LegalMove, } from "../interfaces/message";
import type { Board, Color } from "../interfaces/chess";
import type { Result } from "../interfaces/result";


export interface WebSocketContextType {
    isConnected: boolean;
    sendMessage: (message: WSMessage) => void;
    board: Board | null;
    legalMoves: Record<string, LegalMove[]> | null;
    pgn: string;
    color: Color | null;
    whitePlayerName: string;
    blackPlayerName: string;
    queueLength: number | null;
    inQueue: boolean;
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