import { useState, type RefObject } from "react";
import type { Move, Position } from "../interfaces/chess";
import { MessageTypeMove } from "../interfaces/message";
import type { WSMessage, MoveData } from "../interfaces/message";
import { coordsToSquare } from "../utils/chess";

function Board({ position, socketRef }: { position: Position | null, socketRef: RefObject<WebSocket | null> }) {
    const [selected, setSelected] = useState<{ rank: number; file: number } | null>(null);

    if (!position) {
        return <div>Loading board...</div>;
    }

    function handleClick(i: number, j: number) {
        // Nothing selected
        if (!selected) {
            setSelected({ rank: i, file: j });
            return;
        }

        // Deselect if square clicked again
        if (selected.rank === i && selected.file === j) {
            setSelected(null);
            return;
        }

        // Move piece
        const from = coordsToSquare(selected.rank, selected.file)
        const to = coordsToSquare(i, j)
        console.log(`${from} -> ${to}`);
        const move: Move = {
            from: from,
            to: to
        };
        const moveData: MoveData = { move };
        const message: WSMessage = {
            type: MessageTypeMove,
            data: moveData,
        };
        socketRef.current?.send(JSON.stringify(message));

        setSelected(null);
    }

    return <div>
        {position.board.map((rank, i) => (
            <div key={i} style={{ display: "flex" }}>
                {rank.map((piece, j) => {
                    const imgPath = piece ? `/pieces/${piece.color}-${piece.type}.png` : null;
                    const isSelected = selected?.rank === i && selected?.file === j;
                    return (
                        <div
                            onClick={() => imgPath && handleClick(i, j)}
                            style={{
                                width: 80,
                                height: 80,
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                border: isSelected ? "3px solid #1976d2" : "1px solid #000000",
                                cursor: "pointer",
                                boxSizing: "border-box",
                            }}
                        >
                            {imgPath && <img src={imgPath} alt={`${piece.color} ${piece.type}`} style={{ width: "80px", height: "80px" }} />}
                        </div>
                    );
                })}
            </div>
        ))}
    </div>
}

export default Board;