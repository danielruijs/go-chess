import { Button } from "@mui/material";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import { useAuth } from "../contexts/AuthContext.ts";
import MatchmakingComponent from "../components/Matchmaking.tsx";
import { useNavigate } from "react-router-dom";
import { timeFormats } from "../types/chess.ts";
import { getQueueData } from "../utils/chess.ts";
import { MessageTypeJoinMatch, MessageTypeLeaveMatch, type QueueData } from "../types/message.ts";
import type { TimeFormat } from "../types/chess.ts";

function Home() {
  const { playerInfo } = useAuth();
  const { isConnected, queues, inMatch, sendMessage } = useWebSocket();
  const navigate = useNavigate();

  function handleToggleQueue(timeFormat: TimeFormat, queueData: QueueData | undefined) {
    if (!playerInfo?.displayName || !isConnected || inMatch) {
      return;
    }

    sendMessage({
      type: queueData?.inQueue ? MessageTypeLeaveMatch : MessageTypeJoinMatch,
      data: { timeFormat },
    });
  }

  return (
    <div className="flex flex-col items-center p-8">
      {inMatch ? (
        <Button
          variant="contained"
          disabled={!isConnected}
          onClick={() => {
            navigate("/game");
          }}
          size="large"
        >
          Rejoin Match
        </Button>
      ) : (
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
      )}
    </div>
  );
}

export default Home;
