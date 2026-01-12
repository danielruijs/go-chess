import { useEffect, useRef, useState } from "react";
import type { WSMessage } from "../interfaces/message";
import type { Position } from "../interfaces/chess";
import Board from "../components/Board";

const WS_URL = import.meta.env.VITE_WS_URL;

function Game() {
  const socket = useRef<WebSocket | null>(null);
  const [position, setPosition] = useState<Position | null>(null)

  useEffect(() => {
    // Initialize WebSocket once
    socket.current = new WebSocket(WS_URL);

    socket.current.addEventListener("open", () => {
      console.log("WebSocket connected");
      socket.current?.send("Client connection established");
    });

    socket.current.addEventListener("error", (e) => {
      console.error("WebSocket error:", e);
    });

    socket.current.addEventListener("message", (event) => {
      const message: WSMessage = JSON.parse(event.data);

      if (message.type === "position") {
        const newPosition: Position = message.data;
        setPosition(newPosition);
      }
    });

    // Cleanup on unmount
    return () => {
      socket.current?.close();
    };
  }, []);

  return <div>
    Game
    <button onClick={() => {
      socket.current?.send("Requesting new position");
    }}
    >
      New Position
    </button>
    <Board position={position} />
  </div>
}

export default Game;