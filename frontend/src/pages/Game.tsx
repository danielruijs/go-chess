import BoardComponent from "../components/Board";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Game() {
  const { sendMessage, board, legalMoves, whitePlayerName, blackPlayerName, matchResult } = useWebSocket();

  return (
    <div style={{
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      height: "100vh",
    }}>
      <div style={{ display: "flex", flexDirection: "column", gap: "10px", alignItems: "flex-start" }}>
        <div style={{ fontSize: "18px", fontWeight: "bold" }}>
          {blackPlayerName}
        </div>
        <BoardComponent
          board={board}
          legalMoves={legalMoves}
          sendMessage={sendMessage}
          matchResult={matchResult}
        />
        <div style={{ fontSize: "18px", fontWeight: "bold" }}>
          {whitePlayerName}
        </div>
      </div>
    </div>
  );
}

export default Game;