import BoardComponent from "../components/Board";
import MatchClock from "../components/MatchClock.tsx";
import MatchInfoPanel from "../components/MatchInfoPanel.tsx";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import { useNavigate, useParams } from "react-router-dom";
import { buildMoveRows, getMaterialDiff, pgnToMoves, type MoveRow } from "../utils/chess.ts";
import {
  MessageTypeOfferDraw,
  MessageTypeResign,
  MessageTypeMove,
  type MoveData,
  type WSMessage,
} from "../types/message.ts";
import type { PieceType } from "../types/chess.ts";
import LoadingState from "../components/LoadingState.tsx";
import { useEffect, useMemo } from "react";
import PlayerInfo from "../components/PlayerInfo.tsx";
import useInterpolatedClock from "../hooks/useInterpolatedClock.ts";

function Game() {
  const navigate = useNavigate();
  const { publicId: urlPublicId } = useParams<{ publicId: string }>();
  const {
    sendMessage,
    board,
    legalMoves,
    playerColor,
    activeColor,
    clock,
    whitePlayerName,
    blackPlayerName,
    matchResult,
    pgn,
    isDrawOfferPending,
    isDrawOfferSent,
    setIsDrawOfferSent,
    respondToDrawOffer,
    inMatch,
    publicId: wsPublicId,
    isConnected,
  } = useWebSocket();

  useEffect(() => {
    if (isConnected && wsPublicId !== urlPublicId) {
      navigate("/", { replace: true });
    }
  }, [isConnected, wsPublicId, urlPublicId, navigate]);

  const materialDiff = useMemo(() => getMaterialDiff(board), [board]);
  const interpolatedClock = useInterpolatedClock({ clock, activeColor, inMatch, intervalMs: 100 });
  const groupedMoves = useMemo<MoveRow[]>(() => {
    const moves = pgnToMoves(pgn).map((san, idx) => ({
      san,
      positionIndex: idx + 1,
    }));
    return buildMoveRows(moves);
  }, [pgn]);

  if (!playerColor || !whitePlayerName || !blackPlayerName || !clock) {
    return <LoadingState message="Starting match..." />;
  }
  if (!board) {
    return <LoadingState message="Loading board..." />;
  }

  function handleResign() {
    const message: WSMessage = {
      type: MessageTypeResign,
    };
    sendMessage(message);
  }

  function handleOfferDraw() {
    const message: WSMessage = {
      type: MessageTypeOfferDraw,
    };
    sendMessage(message);
    setIsDrawOfferSent(true);
  }

  function sendMove(from: string, to: string, promotion: PieceType | null) {
    const moveData: MoveData = { from, to, promotion };
    const message: WSMessage = {
      type: MessageTypeMove,
      data: moveData,
    };
    sendMessage(message);
  }

  const displayedClock = interpolatedClock ?? clock;
  const ownTimeRemainingMs =
    playerColor === "white" ? displayedClock.whiteTimeMs : displayedClock.blackTimeMs;
  const opponentTimeRemainingMs =
    playerColor === "white" ? displayedClock.blackTimeMs : displayedClock.whiteTimeMs;

  return (
    <div className="flex-1 flex flex-col items-center justify-center p-2 w-full">
      <div className="flex flex-col gap-1 w-full md:w-auto items-center">
        {/* Opponent Info (Top) */}
        <PlayerInfo
          color={playerColor === "white" ? "black" : "white"}
          name={playerColor === "white" ? blackPlayerName : whitePlayerName}
          materialDiff={materialDiff}
          timeMs={opponentTimeRemainingMs}
          isActive={playerColor !== activeColor}
        />

        {/* Board and Desktop Side Column */}
        <div className="flex gap-4 items-stretch justify-center w-full">
          <div className="chessboard-container flex items-center justify-center">
            <BoardComponent
              color={playerColor}
              board={board}
              legalMoves={legalMoves}
              matchResult={matchResult}
              sendMoveMessage={sendMove}
              onBackToMenu={() => navigate("/")}
              onAnalyze={
                matchResult && urlPublicId ? () => navigate(`/analysis/${urlPublicId}`) : undefined
              }
            />
          </div>

          <div className="hidden md:flex flex-col justify-between w-48 shrink-0">
            <MatchClock timeMs={opponentTimeRemainingMs} isActive={playerColor !== activeColor} />
            <MatchInfoPanel
              groupedMoves={groupedMoves}
              isDrawOfferPending={isDrawOfferPending}
              isDrawOfferSent={isDrawOfferSent}
              inMatch={inMatch}
              onRespondToDrawOffer={respondToDrawOffer}
              onResign={handleResign}
              onOfferDraw={handleOfferDraw}
            />
            <MatchClock timeMs={ownTimeRemainingMs} isActive={playerColor === activeColor} />
          </div>
        </div>

        {/* Own Info (Bottom) */}
        <PlayerInfo
          color={playerColor}
          name={playerColor === "white" ? whitePlayerName : blackPlayerName}
          materialDiff={materialDiff}
          timeMs={ownTimeRemainingMs}
          isActive={playerColor === activeColor}
        />

        {/* Mobile Info Panel (below board) */}
        <div className="md:hidden w-full mt-1">
          <MatchInfoPanel
            groupedMoves={groupedMoves}
            isDrawOfferPending={isDrawOfferPending}
            isDrawOfferSent={isDrawOfferSent}
            inMatch={inMatch}
            onRespondToDrawOffer={respondToDrawOffer}
            onResign={handleResign}
            onOfferDraw={handleOfferDraw}
          />
        </div>
      </div>
    </div>
  );
}

export default Game;
