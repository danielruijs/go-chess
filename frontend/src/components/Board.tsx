import { useState } from "react";
import type { Position } from "../interfaces/chess";

function Board({ position, socket }: { position: Position | null, socket: WebSocket | null }) {
    const [selected, setSelected] = useState<{ row: number; col: number } | null>(null);

    if (!position) {
        return <div>Loading board...</div>;
    }

    function handleClick(i: number, j: number, piece: any) {
        // Nothing selected
        if (!selected) {
            setSelected({ row: i, col: j });
            return;
        }

        // Deselect if square clicked again
        if (selected.row === i && selected.col === j) {
            setSelected(null);
            return;
        }

        // Move piece
        const from = { row: selected.row, col: selected.col };
        const to = { row: i, col: j };
        console.log(`${from.row},${from.col} -> ${to.row},${to.col}`);

        setSelected(null);
    }

    return <div>
        {position.board.map((row, i) => (
            <div key={i} style={{ display: "flex" }}>
                {row.map((piece, j) => {
                    const imgPath = piece ? `/pieces/${piece.color}-${piece.type}.png` : null;
                    const isSelected = selected?.row === i && selected?.col === j;
                    return (
                        <div
                            onClick={() => imgPath && handleClick(i, j, piece)}
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