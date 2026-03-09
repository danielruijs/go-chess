import { Button } from "@mui/material";
import { pgnToMoves } from "../utils/chess";
import { MessageTypeOfferDraw, MessageTypeResign, type WSMessage } from "../interfaces/message";
import { useWebSocket } from "../contexts/WebSocketContext";

function MatchInfoPanel() {
    const { sendMessage, isDrawOfferPending, respondToDrawOffer, inMatch, pgn } = useWebSocket();
    const moves = pgnToMoves(pgn);
    const groupedMoves: { moveNumber: number; whiteMove: string; blackMove: string }[] = [];
    for (let i = 0; i < moves.length; i += 2) {
        groupedMoves.push({
            moveNumber: Math.floor(i / 2) + 1,
            whiteMove: moves[i],
            blackMove: moves[i + 1] || "",
        });
    }

    return (
        <div className="p-2 border border-gray-600 h-72 w-48 flex flex-col box-border">
            <div className="flex-1 overflow-y-auto">
                <table className="w-full border-collapse">
                    <thead>
                        <tr>
                            <th className="text-left w-1/6">#</th>
                            <th className="text-left w-2/5">White</th>
                            <th className="text-left w-2/5">Black</th>
                        </tr>
                    </thead>
                    <tbody>
                        {groupedMoves.map(({ moveNumber, whiteMove, blackMove }, index) => (
                            <tr key={moveNumber} className={index % 2 === 0 ? "bg-gray-200" : ""}>
                                <td>{moveNumber}.</td>
                                <td>{whiteMove}</td>
                                <td>{blackMove}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <div className="mt-2">
                {isDrawOfferPending ? (
                    <div className="flex flex-col items-center gap-2">
                        <div className="text-base font-semibold">Draw offered</div>
                        <div className="flex justify-center gap-2">
                            <Button
                                sx={{ width: "40px", height: "30px", fontSize: "22px" }}
                                variant="contained"
                                color="success"
                                onClick={() => respondToDrawOffer(true)}
                            >
                                ✓
                            </Button>
                            <Button
                                sx={{ width: "40px", height: "30px", fontSize: "22px" }}
                                variant="contained"
                                color="error"
                                onClick={() => respondToDrawOffer(false)}
                            >
                                ✕
                            </Button>
                        </div>
                    </div>
                ) : (
                    <div className="flex justify-center gap-2">
                        <Button
                            sx={{ width: "100px" }}
                            variant="contained"
                            disabled={!inMatch}
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
                            disabled={!inMatch}
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
                )}
            </div>
        </div>
    );
}

export default MatchInfoPanel;