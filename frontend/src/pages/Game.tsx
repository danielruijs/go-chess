import BoardComponent from "../components/Board";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Game() {
  const { sendMessage, board, legalMoves, whitePlayerName, blackPlayerName, matchResult } = useWebSocket();

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
      }}
    >
      <BoardComponent
        board={board}
        legalMoves={legalMoves}
        sendMessage={sendMessage}
        whiteName={whitePlayerName}
        blackName={blackPlayerName}
        matchResult={matchResult}
      />
    </div>
  );
}

export default Game;