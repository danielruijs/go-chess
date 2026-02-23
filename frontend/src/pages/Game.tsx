import { useEffect, useState } from "react";
import { MessageTypeBoard } from "../interfaces/message";
import type { WSMessage, BoardData, LegalMove } from "../interfaces/message";
import type { Board } from "../interfaces/chess";
import BoardComponent from "../components/Board";
import { useWebSocket } from "../contexts/WebSocketContext";

function Game() {
  const { socket } = useWebSocket();
  const [board, setBoard] = useState<Board | null>(null)
  const [legalMoves, setLegalMoves] = useState<Record<string, LegalMove[]> | null>(null);
  const [whiteName, setWhiteName] = useState<string>("");
  const [blackName, setBlackName] = useState<string>("");

  useEffect(() => {
    if (!socket) return;

    function handleMessage(event: MessageEvent) {
      const message: WSMessage = JSON.parse(event.data);

      if (message.type === MessageTypeBoard) {
        const boardData: BoardData = message.data;
        setBoard(boardData.board);
        setLegalMoves(boardData.legalMoves);
        setWhiteName(boardData.whiteName);
        setBlackName(boardData.blackName);
      }
    };

    socket.addEventListener("message", handleMessage);

    return () => {
      socket.removeEventListener("message", handleMessage);
    };
  }, [socket]);

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
        <BoardComponent board={board} legalMoves={legalMoves} socket={socket} whiteName={whiteName} blackName={blackName} />
      </div>
    </div >
  );
}

export default Game;