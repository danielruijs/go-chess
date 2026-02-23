import { createContext, useContext, useEffect, useRef, useState, type ReactNode } from "react";
import type { WSMessage } from "../interfaces/message";

const WS_URL = import.meta.env.VITE_WS_URL;

interface WebSocketContextType {
    socket: WebSocket | null;
    isConnected: boolean;
    sendMessage: (message: WSMessage) => void;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(undefined);

export function WebSocketProvider({ children }: { children: ReactNode }) {
    const socketRef = useRef<WebSocket | null>(null);
    const [isConnected, setIsConnected] = useState(false);

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

        // Cleanup on unmount
        return () => {
            socketRef.current?.close();
        };
    }, []);

    const sendMessage = (message: WSMessage) => {
        socketRef.current?.send(JSON.stringify(message));
    };

    return (
        <WebSocketContext.Provider value={{ socket: socketRef.current, isConnected, sendMessage }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export function useWebSocket() {
    const context = useContext(WebSocketContext);
    if (!context) {
        throw new Error("useWebSocket must be used within a WebSocketProvider");
    }
    return context;
}
