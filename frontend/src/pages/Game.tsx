import BoardComponent from "../components/Board";
import MatchInfoPanel from "../components/MatchInfoPanel.tsx";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Game() {
  const { sendMessage, board, legalMoves, whitePlayerName, blackPlayerName, matchResult } = useWebSocket();

  return (
    <div style={{
      display: "flex",
      alignItems: "center",
      height: "100vh",
    }}>
      <div style={{ display: "flex", gap: "20px", alignItems: "center", justifyContent: "center", width: "100%" }}>
        <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
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
        <MatchInfoPanel></MatchInfoPanel>
      </div>
    </div>
  );
}

export default Game;