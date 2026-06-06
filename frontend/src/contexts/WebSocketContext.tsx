import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import type {
  WSMessage,
  BoardData,
  LegalMove,
  StartMatchData,
  MatchmakingUpdateData,
  EndMatchData,
  ClockData,
  QueueData,
  PlayerInfoData,
} from "../types/message";
import {
  MessageTypeBoard,
  MessageTypeMatchmakingUpdate,
  MessageTypeStartMatch,
  MessageTypeEndMatch,
  MessageTypeDrawOffered,
  MessageTypeDrawDeclined,
  MessageTypeRespondDraw,
  MessageTypePlayerInfo,
} from "../types/message";
import type { Result } from "../types/result";
import type { Color, Board } from "../types/chess";
import { WebSocketContext } from "./WebSocketContext";
import { useAuth } from "./AuthContext";

const WS_URL: string = import.meta.env.VITE_WS_URL as string;

function WebSocketProvider({ children }: { children: ReactNode }) {
  const socketRef = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [board, setBoard] = useState<Board | null>(null);
  const [legalMoves, setLegalMoves] = useState<Record<string, LegalMove[]> | null>(null);
  const [pgn, setPgn] = useState<string>("");
  const [activeColor, setActiveColor] = useState<Color | null>(null);
  const [clock, setClock] = useState<ClockData | null>(null);
  const [playerColor, setPlayerColor] = useState<Color | null>(null);
  const [whitePlayerName, setWhitePlayerName] = useState<string>("");
  const [blackPlayerName, setBlackPlayerName] = useState<string>("");
  const [queues, setQueues] = useState<QueueData[] | null>(null);
  const [inMatch, setInMatch] = useState<boolean>(false);
  const [matchResult, setMatchResult] = useState<Result | null>(null);
  const [isDrawOfferPending, setIsDrawOfferPending] = useState(false);
  const [isDrawDeclinedNoticeOpen, setIsDrawDeclinedNoticeOpen] = useState(false);

  const { playerInfo, isLoading, setPlayerInfo } = useAuth();

  const navigate = useNavigate();
  const navigateRef = useRef(navigate);

  function respondToDrawOffer(accept: boolean) {
    const responseMessage: WSMessage = {
      type: MessageTypeRespondDraw,
      data: { accept },
    };
    sendMessage(responseMessage);
    setIsDrawOfferPending(false);
  }

  function closeDrawDeclinedNotice() {
    setIsDrawDeclinedNoticeOpen(false);
  }

  function resetMatchState() {
    setBoard(null);
    setLegalMoves(null);
    setPgn("");
    setPlayerColor(null);
    setWhitePlayerName("");
    setBlackPlayerName("");
    setInMatch(false);
    setMatchResult(null);
    setIsDrawOfferPending(false);
    setIsDrawDeclinedNoticeOpen(false);
  }

  function sendMessage(message: WSMessage) {
    socketRef.current?.send(JSON.stringify(message));
  }

  const connect = useCallback(() => {
    if (socketRef.current) {
      socketRef.current.close();
    }

    const socket = new WebSocket(WS_URL);
    socketRef.current = socket;

    socket.addEventListener("open", () => {
      if (socketRef.current !== socket) return;
      console.log("WebSocket connected");
      setIsConnected(true);
    });

    socket.addEventListener("close", () => {
      if (socketRef.current !== socket) return;
      console.log("WebSocket disconnected");
      setIsConnected(false);
    });

    socket.addEventListener("error", (e) => {
      if (socketRef.current !== socket) return;
      console.error("WebSocket error:", e);
    });

    socket.addEventListener("message", (event) => {
      if (socketRef.current !== socket) return;
      const message: WSMessage = JSON.parse(event.data);

      switch (message.type) {
        case MessageTypeBoard: {
          const boardData: BoardData = message.data;
          setBoard(boardData.board);
          setLegalMoves(boardData.legalMoves);
          setPgn(boardData.pgn);
          setActiveColor(boardData.activeColor);
          setClock(boardData.clock);
          break;
        }
        case MessageTypeStartMatch: {
          const startMatchData: StartMatchData = message.data;
          resetMatchState();
          setInMatch(true);
          setPlayerColor(startMatchData.color);
          setWhitePlayerName(startMatchData.whitePlayerName);
          setBlackPlayerName(startMatchData.blackPlayerName);
          setClock(startMatchData.clock);
          navigateRef.current("/game");
          break;
        }
        case MessageTypeMatchmakingUpdate: {
          const matchmakingUpdateData: MatchmakingUpdateData = message.data;
          setQueues(matchmakingUpdateData.queues);
          break;
        }
        case MessageTypeEndMatch: {
          const endMatchData: EndMatchData = message.data;
          setMatchResult(endMatchData.result);
          setInMatch(false);
          break;
        }
        case MessageTypeDrawOffered: {
          setIsDrawOfferPending(true);
          break;
        }
        case MessageTypeDrawDeclined: {
          setIsDrawDeclinedNoticeOpen(true);
          break;
        }
        case MessageTypePlayerInfo: {
          const playerInfoData: PlayerInfoData = message.data;
          setPlayerInfo(playerInfoData);
          break;
        }
      }
    });
  }, [setPlayerInfo]);

  // Connect WS on mount once session loading is complete, and disconnect on unmount
  // Also reconnect whenever auth status changes
  useEffect(() => {
    if (isLoading) return;

    connect();

    return () => {
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, [connect, isLoading, playerInfo?.isAuthenticated]);

  return (
    <WebSocketContext.Provider
      value={{
        isConnected,
        sendMessage,
        board,
        legalMoves,
        playerColor,
        whitePlayerName,
        blackPlayerName,
        queues,
        inMatch,
        matchResult,
        pgn,
        activeColor,
        clock,
        isDrawOfferPending,
        isDrawDeclinedNoticeOpen,
        respondToDrawOffer,
        closeDrawDeclinedNotice,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
}

export default WebSocketProvider;
