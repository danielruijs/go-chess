import { useState } from "react";
import type { Board } from "../interfaces/chess";
import { MessageTypeMove } from "../interfaces/message";
import type { LegalMove, WSMessage, MoveData } from "../interfaces/message";
import { coordsToString } from "../utils/chess";

function BoardComponent({
    board,
    legalMoves,
    socket,
    whiteName,
    blackName,
}: {
    board: Board | null,
    legalMoves: Record<string, LegalMove[]> | null,
    socket: WebSocket | null,
    whiteName: string,
    blackName: string,
}) {
    const [selected, setSelected] = useState<{ file: number; rank: number } | null>(null);

    if (!board) {
        return <div>Loading board...</div>;
    }

    const selectedSquare = selected ? coordsToString(selected.file, selected.rank) : null;
    const selectedMoves = selectedSquare && legalMoves ? (legalMoves[selectedSquare] ?? []) : [];
    const selectedMoveTargets = new Set(selectedMoves.map((move) => move.to));

    function handleClick(i: number, j: number) {
        const clickedSquare = coordsToString(i, j);
        const clickedPiece = board?.[i][j] ?? null;
        const clickedPieceWithMoves = clickedPiece ? legalMoves?.[clickedSquare] : null;

        // Nothing selected
        if (!selected) {
            if (clickedPieceWithMoves) {
                setSelected({ file: i, rank: j });
            }
            return;
        }

        // Move piece if the target square is legal
        if (selectedMoveTargets.has(clickedSquare)) {
            const from = coordsToString(selected.file, selected.rank);
            const to = clickedSquare;
            console.log(`${from} -> ${to}`);
            const moveData: MoveData = { from: from, to: to, promotion: null };
            const message: WSMessage = {
                type: MessageTypeMove,
                data: moveData,
            };
            socket?.send(JSON.stringify(message));

            setSelected(null);
            return;
        }

        if (clickedPieceWithMoves && !(selected.file === i && selected.rank === j)) {
            setSelected({ file: i, rank: j });
        } else {
            setSelected(null);
        }
    }

    return (
        <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
            <div style={{ fontSize: "18px", fontWeight: "bold" }}>
                {blackName}
            </div>
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
                    const isMoveTarget = selectedMoveTargets.has(coordsToString(file, rank));
                    return (
                        <div
                            key={index}
                            onClick={() => handleClick(file, rank)}
                            style={{
                                width: 80,
                                height: 80,
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                                position: "relative",
                                border: isSelected ? "3px solid #1976d2" : "1px solid #000000",
                                cursor: "pointer",
                                boxSizing: "border-box",
                            }}
                        >
                            {imgPath && <img src={imgPath} style={{ width: "80px", height: "80px" }} />}
                            {isMoveTarget && (
                                <div
                                    style={{
                                        position: "absolute",
                                        width: 20,
                                        height: 20,
                                        borderRadius: "50%",
                                        backgroundColor: "rgb(255, 0, 0)",
                                    }}
                                />
                            )}
                        </div>
                    );
                })}
            </div>
            <div style={{ fontSize: "18px", fontWeight: "bold" }}>
                {whiteName}
            </div>
        </div>
    );
}

export default BoardComponent;