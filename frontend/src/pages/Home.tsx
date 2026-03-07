import { Button, Stack, TextField } from "@mui/material";
import { MessageTypeJoinMatch } from "../interfaces/message";
import type { WSMessage, JoinMatchData } from "../interfaces/message";
import { useState } from "react";
import { useWebSocket } from "../contexts/WebSocketContext.ts";
import MatchmakingComponent from "../components/Matchmaking.tsx";
import { useNavigate } from "react-router-dom";

function Home() {
    const [playerName, setPlayerName] = useState<string>("");
    const { sendMessage, isConnected, queueLength, inQueue, inMatch } = useWebSocket();
    const navigate = useNavigate();

    return (
        <div className="flex flex-row gap-12 p-2">
            <Stack spacing={2} padding={2} width={"300px"}>
                <TextField label="Name" variant="outlined" onChange={(e) => { setPlayerName(e.target.value) }} />
                <Button
                    variant="contained"
                    disabled={!playerName || !isConnected || inQueue || inMatch}
                    onClick={() => {
                        const joinMatchData: JoinMatchData = { playerName };
                        const message: WSMessage = {
                            type: MessageTypeJoinMatch,
                            data: joinMatchData,
                        };
                        sendMessage(message);
                    }}>
                    Join Matchmaking Queue
                </Button>
                <Button
                    variant="contained"
                    disabled={!isConnected || !inMatch}
                    onClick={() => {
                        navigate("/game");
                    }}>
                    Rejoin Match
                </Button>
            </Stack>
            <div>
                <MatchmakingComponent queueLength={queueLength} inQueue={inQueue} />
            </div>
        </div>
    );
}

export default Home;
