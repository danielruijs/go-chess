import { Button } from "@mui/material";
import { parsePgn } from "../utils/chess";
import { MessageTypeOfferDraw, MessageTypeResign, type WSMessage } from "../interfaces/message";
import { useWebSocket } from "../contexts/WebSocketContext";

function MatchInfoPanel({ pgn }: { pgn: string }) {
    const { sendMessage } = useWebSocket();
    const moves = parsePgn(pgn);
    const groupedMoves: { moveNumber: number; whiteMove: string; blackMove: string }[] = [];
    for (let i = 0; i < moves.length; i += 2) {
        groupedMoves.push({
            moveNumber: Math.floor(i / 2) + 1,
            whiteMove: moves[i],
            blackMove: moves[i + 1] || "",
        });
    }

    return (
        <div style={{ padding: "10px", border: "1px solid #ccc", height: "300px", width: "200px", display: "flex", flexDirection: "column", boxSizing: "border-box" }}>
            <div style={{ flex: 1, overflowY: "auto" }}>
                <table style={{ width: "100%", borderCollapse: "collapse" }}>
                    <thead>
                        <tr>
                            <th style={{ textAlign: "left", width: "16%" }}>#</th>
                            <th style={{ textAlign: "left", width: "42%" }}>White</th>
                            <th style={{ textAlign: "left", width: "42%" }}>Black</th>
                        </tr>
                    </thead>
                    <tbody>
                        {groupedMoves.map(({ moveNumber, whiteMove, blackMove }, index) => (
                            <tr key={moveNumber} style={{ backgroundColor: index % 2 === 0 ? "#d5d5d5" : "transparent" }}>
                                <td>{moveNumber}.</td>
                                <td>{whiteMove}</td>
                                <td>{blackMove}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <div style={{ display: "flex", justifyContent: "center", gap: "10px", marginTop: "10px" }}>
                <Button
                    sx={{ width: "100px" }}
                    variant="contained"
                    onClick={() => {
                        const message: WSMessage = {
                            type: MessageTypeResign,
                        };
                        sendMessage(message);
                    }}
                >
                    Resign
                </Button>
                <Button
                    sx={{ width: "100px" }}
                    variant="contained"
                    onClick={() => {
                        const message: WSMessage = {
                            type: MessageTypeOfferDraw,
                        };
                        sendMessage(message);
                    }}
                >
                    Draw
                </Button>
            </div>
        </div>
    );
}

export default MatchInfoPanel;