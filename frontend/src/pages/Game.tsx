import BoardComponent from "../components/Board";
import MatchInfoPanel from "../components/MatchInfoPanel.tsx";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Game() {
  const { sendMessage, board, legalMoves, color, whitePlayerName, blackPlayerName, matchResult } = useWebSocket();

  return color && whitePlayerName && blackPlayerName ? (
    <div className="flex items-center h-screen">
      <div className="flex gap-4 items-center justify-center w-full">
        <div className="flex flex-col gap-2.5">
          <div className="text-lg font-bold">
            {color === "white" ? blackPlayerName : whitePlayerName}
          </div>
          <BoardComponent
            color={color}
            board={board}
            legalMoves={legalMoves}
            sendMessage={sendMessage}
            matchResult={matchResult}
          />
          <div className="text-lg font-bold">
            {color === "white" ? whitePlayerName : blackPlayerName}
          </div>
        </div>
        <MatchInfoPanel></MatchInfoPanel>
      </div>
    </div>
  ) : "Starting match...";
}

export default Game;