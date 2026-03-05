import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Dialog, DialogContent } from "@mui/material";
import type { Board, Color, PieceType } from "../interfaces/chess";
import { MessageTypeMove } from "../interfaces/message";
import type { LegalMove, WSMessage, MoveData } from "../interfaces/message";
import type { Result } from "../interfaces/result";
import { coordsToString } from "../utils/chess";

function BoardComponent({
    color,
    board,
    legalMoves,
    sendMessage,
    matchResult,
}: {
    color: Color,
    board: Board | null,
    legalMoves: Record<string, LegalMove[]> | null,
    sendMessage: (message: WSMessage) => void,
    matchResult: Result | null,
}) {
    const navigate = useNavigate();
    const [selected, setSelected] = useState<{ file: number; rank: number } | null>(null);
    const [promotionDialog, setPromotionDialog] = useState<{ from: string; to: string; color: "white" | "black" } | null>(null);

    function getResultString() {
        if (!matchResult) return null;

        const outcomeText = {
            white_win: "White wins",
            black_win: "Black wins",
            draw: "Draw",
        }[matchResult.outcome];

        const reasonText = {
            checkmate: "by Checkmate",
            stalemate: "by Stalemate",
            threefold_repetition: "by Threefold Repetition",
            fifty_moves_rule: "by Fifty-Move Rule",
            insufficient_material: "by Insufficient Material",
            resignation: "by Resignation",
            agreed_draw: "by Agreement",
        }[matchResult.reason];

        return `${outcomeText} ${reasonText}`;
    };

    const resultString = getResultString();

    if (!board) {
        return <div>Loading board...</div>;
    }

    const selectedSquare = selected ? coordsToString(selected.file, selected.rank) : null;
    const selectedMoves = selectedSquare && legalMoves ? (legalMoves[selectedSquare] ?? []) : [];
    const selectedMoveTargets = new Set(selectedMoves.map((move) => move.to));

    function sendMove(from: string, to: string, promotion: PieceType | null) {
        const moveData: MoveData = { from, to, promotion };
        const message: WSMessage = {
            type: MessageTypeMove,
            data: moveData,
        };
        sendMessage(message);
    }

    function handlePromotion(pieceType: PieceType) {
        if (promotionDialog) {
            sendMove(promotionDialog.from, promotionDialog.to, pieceType);
            setPromotionDialog(null);
            setSelected(null);
        }
    }

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

            const move = selectedMoves.find(m => m.to === to);
            if (move?.promotion) {
                // Promotion
                const piece = board?.[selected.file][selected.rank];
                if (piece) {
                    setPromotionDialog({ from, to, color: piece.color });
                }
            } else {
                // Regular move
                sendMove(from, to, null);
                setSelected(null);
            }
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
            <div
                style={{
                    position: "relative",
                    display: "grid",
                    gridTemplateColumns: "repeat(8, 80px)",
                    gridTemplateRows: "repeat(8, 80px)",
                }}
            >
                {Array.from({ length: 64 }).map((_, index) => {
                    const displayFile = index % 8;
                    const displayRow = Math.floor(index / 8);
                    const file = color === "white" ? displayFile : 7 - displayFile;
                    const rank = color === "white" ? 7 - displayRow : displayRow;

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
                                backgroundColor: (file + rank) % 2 === 0 ? "#C8C4BE" : "#F5F3F0",
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
                {resultString && (
                    <div
                        style={{
                            position: "absolute",
                            top: "50%",
                            left: "50%",
                            transform: "translate(-50%, -50%)",
                            fontSize: "18px",
                            fontWeight: "bold",
                            padding: "20px",
                            backgroundColor: "#e4e4e4",
                            border: "2px solid #1976d2",
                            borderRadius: "4px",
                            textAlign: "center",
                            color: "#1976d2",
                            minWidth: "200px",
                            zIndex: 10,
                            display: "flex",
                            flexDirection: "column",
                            gap: "15px",
                            alignItems: "center",
                        }}
                    >
                        <div>{resultString}</div>
                        <Button
                            onClick={() => navigate("/")}
                            variant="contained"
                        >
                            Back to Menu
                        </Button>
                    </div>
                )}
            </div>
            <Dialog open={!!promotionDialog} onClose={() => setPromotionDialog(null)}>
                <DialogContent>
                    <div style={{ display: "flex", gap: "10px" }}>
                        {(["queen", "rook", "bishop", "knight"] as const).map((pieceType) => (
                            <div
                                key={pieceType}
                                onClick={() => handlePromotion(pieceType)}
                                style={{
                                    cursor: "pointer",
                                    border: "2px solid #1976d2",
                                    borderRadius: "4px",
                                    backgroundColor: "#f0f0f0",
                                }}
                            >
                                <img
                                    src={`/pieces/${promotionDialog?.color}-${pieceType}.png`}
                                    style={{ width: "80px", height: "80px" }}
                                />
                            </div>
                        ))}
                    </div>
                </DialogContent>
            </Dialog>
        </div>
    );
}

export default BoardComponent;