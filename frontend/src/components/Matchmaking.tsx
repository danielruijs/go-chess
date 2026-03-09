import Button from "@mui/material/Button";
import { MessageTypeJoinMatch } from "../interfaces/message";
import type { WSMessage, JoinMatchData } from "../interfaces/message";
import type { QueueData } from "../interfaces/message";
import type { TimeFormat } from "../interfaces/chess";
import { useWebSocket } from "../contexts/WebSocketContext";

function MatchmakingComponent({ timeFormat, queueData, playerName }: { timeFormat: TimeFormat, queueData: QueueData | undefined; playerName: string }) {
    const { isConnected, sendMessage, inMatch } = useWebSocket();
    const queueLength = queueData ? queueData.queueLength : null;
    const inQueue = queueData ? queueData.inQueue : false;

    return (
        <div className="p-2">
            <p className="text-2xl font-bold mb-4">{`${timeFormat.initialMs / 1000}+${timeFormat.incrementMs / 1000}`}</p>
            <Button
                variant="contained"
                disabled={!playerName || !isConnected || inQueue || inMatch}
                onClick={() => {
                    const joinMatchData: JoinMatchData = { playerName, timeFormat };
                    const message: WSMessage = {
                        type: MessageTypeJoinMatch,
                        data: joinMatchData,
                    };
                    sendMessage(message);
                }}>
                Join Matchmaking Queue
            </Button>
            <div className="text-xl font-semibold mb-2">Matchmaking Queue</div>
            {queueLength !== null ? (
                <p>Current queue length: {queueLength}</p>
            ) : (
                <p>Loading queue length...</p>
            )}
            {inQueue && (
                <p className="mt-2 text-gray-600">You are in the matchmaking queue.</p>
            )}
        </div>
    );
}

export default MatchmakingComponent;