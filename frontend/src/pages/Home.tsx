import { Button, Stack, TextField } from "@mui/material";
import { MessageTypeJoinMatch } from "../interfaces/message";
import type { WSMessage, JoinMatchData } from "../interfaces/message";
import { useState } from "react";
import { useWebSocket } from "../contexts/WebSocketContext.ts";

function Home() {
    const [playerName, setPlayerName] = useState<string>("");
    const { sendMessage, isConnected } = useWebSocket();

    return (
        <div>
            <Stack spacing={2} padding={2} width={"300px"}>
                <TextField label="Name" variant="outlined" onChange={(e) => { setPlayerName(e.target.value) }} />
                <Button
                    variant="contained"
                    disabled={!playerName || !isConnected}
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
            </Stack>
        </div>
    );
}

export default Home;
