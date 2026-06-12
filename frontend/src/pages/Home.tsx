import { Button, Stack } from "@mui/material";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import { useAuth } from "../contexts/AuthContext.ts";
import MatchmakingComponent from "../components/Matchmaking.tsx";
import { useNavigate } from "react-router-dom";
import { timeFormats } from "../types/chess.ts";
import { getQueueData } from "../utils/chess.ts";
import {
  MessageTypeJoinMatch,
  MessageTypeLeaveMatch,
  type WSMessage,
  type JoinMatchData,
  type QueueData,
} from "../types/message.ts";
import type { TimeFormat } from "../types/chess.ts";

function Home() {
  const { playerInfo } = useAuth();
  const { isConnected, queues, inMatch, sendMessage } = useWebSocket();
  const navigate = useNavigate();

  function handleToggleQueue(timeFormat: TimeFormat, queueData: QueueData | undefined) {
    if (!playerInfo?.displayName || !isConnected || inMatch) {
      return;
    }

    if (queueData?.inQueue) {
      const leaveData = { timeFormat };
      const message: WSMessage = { type: MessageTypeLeaveMatch, data: leaveData };
      sendMessage(message);
      return;
    }

    const joinMatchData: JoinMatchData = { timeFormat };
    const message: WSMessage = {
      type: MessageTypeJoinMatch,
      data: joinMatchData,
    };
    sendMessage(message);
  }

  return (
    <div className="flex flex-row gap-12 p-4">
      <Stack spacing={2} className="w-75">
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
                  displayName={playerInfo?.displayName ?? ""}
                  isConnected={isConnected}
                  inMatch={inMatch}
                  onToggleQueue={() => handleToggleQueue(timeFormat, queueData)}
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
