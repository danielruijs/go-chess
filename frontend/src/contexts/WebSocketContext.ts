import { createContext, useContext  } from "react";
import type { WSMessage, LegalMove,  } from "../interfaces/message";
import type { Board } from "../interfaces/chess";


export interface WebSocketContextType {
    isConnected: boolean;
    sendMessage: (message: WSMessage) => void;
    board: Board | null;
    legalMoves: Record<string, LegalMove[]> | null;
    whitePlayerName: string;
    blackPlayerName: string;
    queueLength: number | null;
    inQueue: boolean;
}

export const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export function useWebSocket() {
    const context = useContext(WebSocketContext);
    if (!context) {
        throw new Error("useWebSocket must be used within a WebSocketProvider");
    }
    return context;
}