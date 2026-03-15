import { Button, Stack, TextField } from "@mui/material";
import { useState } from "react";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import MatchmakingComponent from "../components/Matchmaking.tsx";
import { useNavigate } from "react-router-dom";
import { timeFormats } from "../interfaces/chess.ts";
import { getQueueData } from "../utils/chess.ts";

function Home() {
  const [playerName, setPlayerName] = useState<string>("");
  const { isConnected, queues, inMatch } = useWebSocket();
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex flex-row gap-12 p-2">
      <Stack spacing={2} padding={2} width={"300px"}>
        <TextField
          label="Name"
          variant="outlined"
          onChange={(e) => {
            setPlayerName(e.target.value);
          }}
        />

        <Button
          variant="contained"
          disabled={!isConnected || !inMatch}
          onClick={() => {
            navigate("/game");
          }}
        >
          Rejoin Match
        </Button>
      </Stack>
      <div className="flex flex-col gap-6">
        {timeFormats.map((group, groupIdx) => (
          <div key={groupIdx} className="flex flex-wrap gap-4">
            {group.map((timeFormat) => (
              <MatchmakingComponent
                timeFormat={timeFormat}
                queueData={getQueueData(queues, timeFormat)}
                playerName={playerName}
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}

export default Home;
