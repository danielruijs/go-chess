import { useEffect, useRef, useState } from "react";
import { MessageTypeBoard, MessageTypeJoinMatch } from "../interfaces/message";
import type { WSMessage, JoinMatchData, BoardData } from "../interfaces/message";
import type { Board } from "../interfaces/chess";
import { Button, Stack, TextField } from "@mui/material";
import BoardComponent from "../components/Board";

const WS_URL = import.meta.env.VITE_WS_URL;

function Game() {
  const socket = useRef<WebSocket | null>(null);
  const [board, setBoard] = useState<Board | null>(null)
  const [playerName, setPlayerName] = useState<string>("");

  useEffect(() => {
    // Initialize WebSocket once
    socket.current = new WebSocket(WS_URL);

    socket.current.addEventListener("open", () => {
      console.log("WebSocket connected");
    });

    socket.current.addEventListener("error", (e) => {
      console.error("WebSocket error:", e);
    });

    socket.current.addEventListener("message", (event) => {
      const message: WSMessage = JSON.parse(event.data);

      if (message.type === MessageTypeBoard) {
        const boardData: BoardData = message.data;
        setBoard(boardData.board);
      }
    });

    // Cleanup on unmount
    return () => {
      socket.current?.close();
    };
  }, []);

  return (
    <div>
      <Stack spacing={2} padding={2} width={"300px"}>
        <TextField label="Name" variant="outlined" onChange={(e) => { setPlayerName(e.target.value) }} />

        <Button
          variant="contained"
          disabled={!playerName}
          onClick={() => {
            const joinMatchData: JoinMatchData = { playerName };
            const message: WSMessage = {
              type: MessageTypeJoinMatch,
              data: joinMatchData,
            };
            socket.current?.send(JSON.stringify(message));
          }}>
          Join Matchmaking Queue
        </Button>
      </Stack>

      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexGrow: 1,
        }}
      >
        <BoardComponent board={board} socketRef={socket} />
      </div>
    </div >
  );
}

export default Game;