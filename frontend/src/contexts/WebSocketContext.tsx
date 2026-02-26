import { useEffect, useRef, useState, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import type { WSMessage, BoardData, LegalMove, StartMatchData, MatchmakingUpdateData, EndMatchData } from "../interfaces/message";
import { MessageTypeBoard, MessageTypeMatchmakingUpdate, MessageTypeStartMatch, MessageTypeEndMatch } from "../interfaces/message";
import type { Result } from "../interfaces/result";
import type { Board } from "../interfaces/chess";
import { WebSocketContext } from "./WebSocketContext";

const WS_URL = import.meta.env.VITE_WS_URL;

export function WebSocketProvider({ children }: { children: ReactNode }) {
    const socketRef = useRef<WebSocket | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    const [board, setBoard] = useState<Board | null>(null);
    const [legalMoves, setLegalMoves] = useState<Record<string, LegalMove[]> | null>(null);
    const [whitePlayerName, setWhitePlayerName] = useState<string>("");
    const [blackPlayerName, setBlackPlayerName] = useState<string>("");
    const [queueLength, setQueueLength] = useState<number | null>(null);
    const [inQueue, setInQueue] = useState<boolean>(false);
    const [matchResult, setMatchResult] = useState<Result | null>(null);
    const navigate = useNavigate();
    const navigateRef = useRef(navigate);

    useEffect(() => {
        // Initialize WebSocket once
        socketRef.current = new WebSocket(WS_URL);

        socketRef.current.addEventListener("open", () => {
            console.log("WebSocket connected");
            setIsConnected(true);
        });

        socketRef.current.addEventListener("close", () => {
            console.log("WebSocket disconnected");
            setIsConnected(false);
        });

        socketRef.current.addEventListener("error", (e) => {
            console.error("WebSocket error:", e);
        });

        socketRef.current.addEventListener("message", (event) => {
            const message: WSMessage = JSON.parse(event.data);

            switch (message.type) {
                case MessageTypeBoard: {
                    const boardData: BoardData = message.data;
                    setBoard(boardData.board);
                    setLegalMoves(boardData.legalMoves);
                    break;
                }
                case MessageTypeStartMatch: {
                    const startMatchData: StartMatchData = message.data;
                    setWhitePlayerName(startMatchData.whitePlayerName);
                    setBlackPlayerName(startMatchData.blackPlayerName);
                    navigateRef.current("/game");
                    break;
                }
                case MessageTypeMatchmakingUpdate: {
                    const matchmakingUpdateData: MatchmakingUpdateData = message.data;
                    setQueueLength(matchmakingUpdateData.queueLength);
                    setInQueue(matchmakingUpdateData.inQueue);
                    break;
                }
                case MessageTypeEndMatch: {
                    const endMatchData: EndMatchData = message.data;
                    setMatchResult(endMatchData.result);
                    break;
                }
            }
        });

        // Cleanup on unmount
        return () => {
            socketRef.current?.close();
        };
    }, []);

    const sendMessage = (message: WSMessage) => {
        socketRef.current?.send(JSON.stringify(message));
    };

    return (
        <WebSocketContext.Provider value={{
            isConnected,
            sendMessage,
            board,
            legalMoves,
            whitePlayerName,
            blackPlayerName,
            queueLength,
            inQueue,
            matchResult,
        }}>
            {children}
        </WebSocketContext.Provider>
    );
}


