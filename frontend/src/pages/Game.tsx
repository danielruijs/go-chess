import BoardComponent from "../components/Board";
import MatchInfoPanel from "../components/MatchInfoPanel.tsx";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Game() {
  const { sendMessage, board, legalMoves, pgn, color, whitePlayerName, blackPlayerName, matchResult } = useWebSocket();

  return color && whitePlayerName && blackPlayerName ? (
    <div style={{
      display: "flex",
      alignItems: "center",
      height: "100vh",
    }}>
      <div style={{ display: "flex", gap: "15px", alignItems: "center", justifyContent: "center", width: "100%" }}>
        <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
          <div style={{ fontSize: "18px", fontWeight: "bold" }}>
            {color === "white" ? blackPlayerName : whitePlayerName}
          </div>
          <BoardComponent
            color={color}
            board={board}
            legalMoves={legalMoves}
            sendMessage={sendMessage}
            matchResult={matchResult}
          />
          <div style={{ fontSize: "18px", fontWeight: "bold" }}>
            {color === "white" ? whitePlayerName : blackPlayerName}
          </div>
        </div>
        <MatchInfoPanel pgn={pgn}></MatchInfoPanel>
      </div>
    </div>
  ) : "Starting match...";
}

export default Game;