import BoardComponent from "../components/Board";
import MatchClock from "../components/MatchClock.tsx";
import MatchInfoPanel from "../components/MatchInfoPanel.tsx";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Game() {
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
  } = useWebSocket();
  if (!playerColor || !whitePlayerName || !blackPlayerName || !clock) {
    return "Starting match...";
  }

  const ownTimeRemainingMs = playerColor === "white" ? clock.whiteTimeMs : clock.blackTimeMs;
  const opponentTimeRemainingMs = playerColor === "white" ? clock.blackTimeMs : clock.whiteTimeMs;

  return (
    <div className="flex items-center h-screen">
      <div className="flex gap-4 items-center justify-center w-full">
        <div className="flex flex-col gap-2.5">
          <div className="text-lg font-bold">
            {playerColor === "white" ? blackPlayerName : whitePlayerName}
          </div>
          <div className="flex gap-4">
            <BoardComponent
              color={playerColor}
              board={board}
              legalMoves={legalMoves}
              sendMessage={sendMessage}
              matchResult={matchResult}
            />
            <div className="flex flex-col justify-between">
              <MatchClock timeMs={opponentTimeRemainingMs} isActive={playerColor !== activeColor} />
              <MatchInfoPanel />
              <MatchClock timeMs={ownTimeRemainingMs} isActive={playerColor === activeColor} />
            </div>
          </div>
          <div className="text-lg font-bold">
            {playerColor === "white" ? whitePlayerName : blackPlayerName}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Game;
