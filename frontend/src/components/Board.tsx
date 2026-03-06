import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Dialog, DialogContent } from "@mui/material";
import { DragDropProvider, useDraggable, useDroppable } from "@dnd-kit/react";
import type { DragStartEvent, DragEndEvent } from "@dnd-kit/react";
import type { Board, Color, PieceType, Square } from "../interfaces/chess";
import { MessageTypeMove } from "../interfaces/message";
import type { LegalMove, WSMessage, MoveData } from "../interfaces/message";
import type { Result } from "../interfaces/result";
import { coordsToString } from "../utils/chess";

const OUTCOME_TEXT = {
    white_win: "White wins",
    black_win: "Black wins",
    draw: "Draw",
} as const;

const REASON_TEXT = {
    checkmate: "by Checkmate",
    stalemate: "by Stalemate",
    threefold_repetition: "by Threefold Repetition",
    fifty_moves_rule: "by Fifty-Move Rule",
    insufficient_material: "by Insufficient Material",
    resignation: "by Resignation",
    agreed_draw: "by Agreement",
} as const;

const PROMOTION_PIECES = ["queen", "rook", "bishop", "knight"] as const;

type DragStartPayload = Parameters<DragStartEvent>[0];
type DragEndPayload = Parameters<DragEndEvent>[0];

function DraggablePiece({
    square,
    pieceColor,
    imgPath,
    disabled,
}: {
    square: Square;
    pieceColor: Color;
    imgPath: string;
    disabled: boolean;
}) {
    const { ref, isDragging } = useDraggable<{ square: Square, color: Color }>({
        id: `piece-${square.file}-${square.rank}`,
        data: { square: square, color: pieceColor },
        disabled,
    });

    return (
        <img
            ref={ref}
            src={imgPath}
            style={{
                width: "80px",
                height: "80px",
                opacity: isDragging ? 0.55 : 1,
                cursor: disabled ? "pointer" : "grab",
                userSelect: "none",
                WebkitUserSelect: "none",
            }}
            draggable={false}
        />
    );
}

function BoardSquare({
    square,
    isSelected,
    isMoveTarget,
    onClick,
    children,
}: {
    square: Square;
    isSelected: boolean;
    isMoveTarget: boolean;
    onClick: () => void;
    children: React.ReactNode;
}) {
    const { ref, isDropTarget } = useDroppable<{ square: Square }>({
        id: `square-${square.file}-${square.rank}`,
        data: { square: square },
    });

    return (
        <div
            ref={ref}
            onClick={onClick}
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
                backgroundColor: (square.file + square.rank) % 2 === 0 ? "#C8C4BE" : "#F5F3F0",
                outline: isDropTarget ? "2px dashed #1976d2" : "none",
                outlineOffset: isDropTarget ? "-2px" : 0,
                userSelect: "none",
                WebkitUserSelect: "none",
            }}
        >
            {children}
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
}

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
    const [selected, setSelected] = useState<Square | null>(null);
    const [promotionDialog, setPromotionDialog] = useState<{ from: string; to: string; color: Color } | null>(null);

    if (!board) {
        return <div>Loading board...</div>;
    }
    const currentBoard = board;

    function getResultString() {
        if (!matchResult) return null;

        const outcomeText = OUTCOME_TEXT[matchResult.outcome];
        const reasonText = REASON_TEXT[matchResult.reason];

        return `${outcomeText} ${reasonText}`;
    }

    const resultString = getResultString();

    const selectedSquareString = selected ? coordsToString(selected) : null;
    const selectedMoves = selectedSquareString && legalMoves ? (legalMoves[selectedSquareString] ?? []) : [];
    const selectedMoveTargets = new Set(selectedMoves.map((move) => move.to));

    function hasLegalMoves(square: Square) {
        const squareString = coordsToString(square);
        return (legalMoves?.[squareString]?.length ?? 0) > 0;
    }

    function sendMove(from: string, to: string, promotion: PieceType | null) {
        const moveData: MoveData = { from, to, promotion };
        const message: WSMessage = {
            type: MessageTypeMove,
            data: moveData,
        };
        sendMessage(message);
    }

    function tryMove(from: Square, to: Square, sourceColor: Color) {
        const fromSquareString = coordsToString(from);
        const toSquareString = coordsToString(to);

        if (!legalMoves || (from.file === to.file && from.rank === to.rank)) {
            setSelected(null);
            return;
        }

        const move = selectedMoves.find(m => m.to === toSquareString);
        if (!move) {
            setSelected(null);
            return;
        }

        if (move.promotion) {
            setPromotionDialog({ from: fromSquareString, to: toSquareString, color: sourceColor });
            return;
        }

        sendMove(fromSquareString, toSquareString, null);
        setSelected(null);
    }

    function handlePromotion(pieceType: PieceType) {
        if (promotionDialog) {
            sendMove(promotionDialog.from, promotionDialog.to, pieceType);
            setPromotionDialog(null);
            setSelected(null);
        }
    }

    function handleClick(square: Square, isMoveTarget: boolean) {
        const clickedPiece = currentBoard[square.file][square.rank];
        const clickedPieceHasMoves = clickedPiece && hasLegalMoves(square);

        // Nothing selected
        if (!selected) {
            if (clickedPieceHasMoves) {
                setSelected(square);
            }
            return;
        }

        // Move piece if the target square is legal
        if (isMoveTarget) {
            const sourcePiece = currentBoard[selected.file][selected.rank];
            if (sourcePiece) {
                tryMove(selected, square, sourcePiece.color);
                return;
            }
        }

        if (clickedPieceHasMoves && !(selected.file === square.file && selected.rank === square.rank)) {
            setSelected(square);
        } else {
            setSelected(null);
        }
    }

    function handleDragStart(event: DragStartPayload) {
        const sourceSquare = event.operation.source?.data.square;
        if (!sourceSquare || !hasLegalMoves(sourceSquare)) {
            setSelected(null);
            return;
        }

        setSelected(sourceSquare);
    }

    function handleDragEnd(event: DragEndPayload) {
        const from = event.operation.source?.data.square;
        const to = event.operation.target?.data.square;
        const sourceColor = event.operation.source?.data.color;

        if (!from || !to || !sourceColor || event.canceled) {
            setSelected(null);
            return;
        }

        tryMove(from, to, sourceColor);
    }

    return (
        <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
            <DragDropProvider onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
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
                        const square: Square = { file, rank };

                        const piece = board[file][rank];
                        const imgPath = piece ? `/pieces/${piece.color}-${piece.type}.png` : null;
                        const isSelected = selected?.file === file && selected?.rank === rank;
                        const isMoveTarget = selectedMoveTargets.has(coordsToString(square));
                        const canDragPiece = piece && hasLegalMoves(square);

                        return (
                            <BoardSquare
                                key={index}
                                square={square}
                                onClick={() => handleClick(square, isMoveTarget)}
                                isSelected={isSelected}
                                isMoveTarget={isMoveTarget}
                            >
                                {imgPath && piece && (
                                    <DraggablePiece
                                        square={square}
                                        pieceColor={piece.color}
                                        imgPath={imgPath}
                                        disabled={!canDragPiece}
                                    />
                                )}
                            </BoardSquare>
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
            </DragDropProvider>
            <Dialog open={!!promotionDialog} onClose={() => setPromotionDialog(null)}>
                <DialogContent>
                    <div style={{ display: "flex", gap: "10px" }}>
                        {PROMOTION_PIECES.map((pieceType) => (
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