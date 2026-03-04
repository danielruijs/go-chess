import { createContext, useContext } from "react";
import type { WSMessage, LegalMove, } from "../interfaces/message";
import type { Board } from "../interfaces/chess";
import type { Result } from "../interfaces/result";


export interface WebSocketContextType {
    isConnected: boolean;
    sendMessage: (message: WSMessage) => void;
    board: Board | null;
    legalMoves: Record<string, LegalMove[]> | null;
    pgn: string;
    whitePlayerName: string;
    blackPlayerName: string;
    queueLength: number | null;
    inQueue: boolean;
    matchResult: Result | null;
}

export const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export function useWebSocket() {
    const context = useContext(WebSocketContext);
    if (!context) {
        throw new Error("useWebSocket must be used within a WebSocketProvider");
    }
    return context;
}