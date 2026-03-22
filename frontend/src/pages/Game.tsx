import BoardComponent from "../components/Board";
import MatchClock from "../components/MatchClock.tsx";
import MatchInfoPanel from "../components/MatchInfoPanel.tsx";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import { useNavigate } from "react-router-dom";
import { getMaterialDiff, pgnToMoves } from "../utils/chess.ts";
import {
  MessageTypeOfferDraw,
  MessageTypeResign,
  type MoveData,
  type WSMessage,
} from "../types/message.ts";
import { MessageTypeMove } from "../types/message.ts";
import type { PieceType } from "../types/chess.ts";
import { useMemo } from "react";
import PlayerInfo from "../components/PlayerInfo.tsx";

function Game() {
  const navigate = useNavigate();
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
    respondToDrawOffer,
    inMatch,
  } = useWebSocket();
  const materialDiff = useMemo(() => getMaterialDiff(board), [board]);

  if (!playerColor || !whitePlayerName || !blackPlayerName || !clock) {
    return "Starting match...";
  }
  if (!board) {
    return "Loading board...";
  }

  const moves = pgnToMoves(pgn);
  const groupedMoves: { moveNumber: number; whiteMove: string; blackMove: string }[] = [];
  for (let i = 0; i < moves.length; i += 2) {
    groupedMoves.push({
      moveNumber: Math.floor(i / 2) + 1,
      whiteMove: moves[i],
      blackMove: moves[i + 1] || "",
    });
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
  }

  function sendMove(from: string, to: string, promotion: PieceType | null) {
    const moveData: MoveData = { from, to, promotion };
    const message: WSMessage = {
      type: MessageTypeMove,
      data: moveData,
    };
    sendMessage(message);
  }

  const ownTimeRemainingMs = playerColor === "white" ? clock.whiteTimeMs : clock.blackTimeMs;
  const opponentTimeRemainingMs = playerColor === "white" ? clock.blackTimeMs : clock.whiteTimeMs;

  return (
    <div className="flex items-center h-screen">
      <div className="flex gap-4 items-center justify-center w-full">
        <div className="flex flex-col gap-2.5">
          <PlayerInfo
            color={playerColor === "white" ? "black" : "white"}
            name={playerColor === "white" ? blackPlayerName : whitePlayerName}
            materialDiff={materialDiff}
          />
          <div className="flex gap-4">
            <BoardComponent
              color={playerColor}
              board={board}
              legalMoves={legalMoves}
              matchResult={matchResult}
              sendMoveMessage={sendMove}
              onBackToMenu={() => navigate("/")}
            />
            <div className="flex flex-col justify-between">
              <MatchClock timeMs={opponentTimeRemainingMs} isActive={playerColor !== activeColor} />
              <MatchInfoPanel
                groupedMoves={groupedMoves}
                isDrawOfferPending={isDrawOfferPending}
                inMatch={inMatch}
                onRespondToDrawOffer={respondToDrawOffer}
                onResign={handleResign}
                onOfferDraw={handleOfferDraw}
              />
              <MatchClock timeMs={ownTimeRemainingMs} isActive={playerColor === activeColor} />
            </div>
          </div>
          <PlayerInfo
            color={playerColor}
            name={playerColor === "white" ? whitePlayerName : blackPlayerName}
            materialDiff={materialDiff}
          />
        </div>
      </div>
    </div>
  );
}

export default Game;
