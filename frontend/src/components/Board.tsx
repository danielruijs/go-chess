import { useState, type RefObject } from "react";
import type { Board } from "../interfaces/chess";
import { MessageTypeMove } from "../interfaces/message";
import type { WSMessage, MoveData } from "../interfaces/message";
import { coordsToString } from "../utils/chess";

function BoardComponent({ board, socketRef }: { board: Board | null, socketRef: RefObject<WebSocket | null> }) {
    const [selected, setSelected] = useState<{ file: number; rank: number } | null>(null);

    if (!board) {
        return <div>Loading board...</div>;
    }

    function handleClick(i: number, j: number) {
        // Nothing selected
        if (!selected) {
            setSelected({ file: i, rank: j });
            return;
        }

        // Deselect if square clicked again
        if (selected.file === i && selected.rank === j) {
            setSelected(null);
            return;
        }

        // Move piece
        const from = coordsToString(selected.file, selected.rank)
        const to = coordsToString(i, j)
        console.log(`${from} -> ${to}`);
        const moveData: MoveData = { from: from, to: to, promotion: null };
        const message: WSMessage = {
            type: MessageTypeMove,
            data: moveData,
        };
        socketRef.current?.send(JSON.stringify(message));

        setSelected(null);
    }

    return (
        <div
            style={{
                display: "grid",
                gridTemplateColumns: "repeat(8, 80px)",
                gridTemplateRows: "repeat(8, 80px)",
            }}
        >
            {Array.from({ length: 64 }).map((_, index) => {
                const file = index % 8;                 // a -> h
                const rank = 7 - Math.floor(index / 8); // 8 -> 1

                const piece = board[file][rank];
                const imgPath = piece ? `/pieces/${piece.color}-${piece.type}.png` : null;
                const isSelected = selected?.file === file && selected?.rank === rank;
                return (
                    <div
                        key={index}
                        onClick={() => piece && handleClick(file, rank)}
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
    );
}

export default BoardComponent;