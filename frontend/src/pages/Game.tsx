import BoardComponent from "../components/Board";
import { useWebSocket } from "../contexts/WebSocketContext";

function Game() {
  const { sendMessage, board, legalMoves, whitePlayerName, blackPlayerName } = useWebSocket();

  return (
    <div>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexGrow: 1,
        }}
      >
        <BoardComponent
          board={board}
          legalMoves={legalMoves}
          sendMessage={sendMessage}
          whiteName={whitePlayerName}
          blackName={blackPlayerName}
        />
      </div>
    </div >
  );
}

export default Game;