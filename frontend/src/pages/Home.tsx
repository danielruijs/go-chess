import { Button, Stack, TextField } from "@mui/material";
import { useState } from "react";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import MatchmakingComponent from "../components/Matchmaking.tsx";
import { useNavigate } from "react-router-dom";
import { timeFormats } from "../types/chess.ts";
import { getQueueData } from "../utils/chess.ts";
import {
  MessageTypeJoinMatch,
  type WSMessage,
  type JoinMatchData,
  type QueueData,
} from "../types/message.ts";
import type { TimeFormat } from "../types/chess.ts";

function Home() {
  const [playerName, setPlayerName] = useState<string>("");
  const { isConnected, queues, inMatch, sendMessage } = useWebSocket();
  const navigate = useNavigate();

  function handleJoinQueue(timeFormat: TimeFormat, queueData: QueueData | undefined) {
    if (!playerName || !isConnected || inMatch || queueData?.inQueue) {
      return;
    }

    const joinMatchData: JoinMatchData = { playerName, timeFormat };
    const message: WSMessage = {
      type: MessageTypeJoinMatch,
      data: joinMatchData,
    };
    sendMessage(message);
  }

  return (
    <div className="min-h-screen flex flex-row gap-12 p-4">
      <Stack spacing={2} width={"300px"}>
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
            {group.map((timeFormat) => {
              const queueData = getQueueData(queues, timeFormat);
              return (
                <MatchmakingComponent
                  key={`${timeFormat.initialMs}-${timeFormat.incrementMs}`}
                  timeFormat={timeFormat}
                  queueData={queueData}
                  playerName={playerName}
                  isConnected={isConnected}
                  inMatch={inMatch}
                  onJoinQueue={() => handleJoinQueue(timeFormat, queueData)}
                />
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}

export default Home;
