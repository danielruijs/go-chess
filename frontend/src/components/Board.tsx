import { useMemo, useState, type ReactNode } from "react";
import { Button, Dialog, DialogContent } from "@mui/material";
import { DragDropProvider, useDraggable, useDroppable } from "@dnd-kit/react";
import type { DragStartEvent, DragEndEvent } from "@dnd-kit/react";
import type { Board, Color, PieceType, Square } from "../types/chess";
import type { LegalMove } from "../types/message";
import { OUTCOME_TEXT, REASON_TEXT, type Result } from "../types/result";
import { coordsToString, displayIndexToSquare } from "../utils/chess";

const PROMOTION_PIECES = ["queen", "rook", "bishop", "knight"] as const;

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
  const { ref, isDragging } = useDraggable<{ square: Square; color: Color }>({
    id: `piece-${square.file}-${square.rank}`,
    data: { square: square, color: pieceColor },
    disabled,
  });

  return (
    <img
      ref={ref}
      src={imgPath}
      className="select-none"
      style={{
        opacity: isDragging ? 0.55 : 1,
        cursor: disabled ? "pointer" : "grab",
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
  children: ReactNode;
}) {
  const { ref, isDropTarget } = useDroppable<{ square: Square }>({
    id: `square-${square.file}-${square.rank}`,
    data: { square: square },
  });

  return (
    <div
      ref={ref}
      onClick={onClick}
      className={`w-full h-full flex items-center justify-center relative box-border select-none }`}
      style={{
        border: isSelected ? "3px solid #1976d2" : "1px solid #000000",
        cursor: "pointer",
        backgroundColor: (square.file + square.rank) % 2 === 0 ? "#C8C4BE" : "#F5F3F0",
        outline: isDropTarget ? "2px dashed #1976d2" : "none",
        outlineOffset: isDropTarget ? "-2px" : 0,
      }}
    >
      {children}
      {isMoveTarget && <div className="absolute w-[25%] h-[25%] rounded-full bg-red-600" />}
    </div>
  );
}

type BoardComponentProps = {
  color: Color;
  board: Board;
  legalMoves: Record<string, LegalMove[]> | null;
  matchResult: Result | null;
  sendMoveMessage: (from: string, to: string, promotion: PieceType | null) => void;
  onBackToMenu: () => void;
};

function BoardComponent({
  color,
  board,
  legalMoves,
  matchResult,
  sendMoveMessage,
  onBackToMenu,
}: BoardComponentProps) {
  const [selected, setSelected] = useState<Square | null>(null);
  const [promotionDialog, setPromotionDialog] = useState<{
    from: string;
    to: string;
    color: Color;
  } | null>(null);
  const { selectedMoves, selectedMoveTargets } = useMemo(() => {
    const selectedSquareString = selected ? coordsToString(selected) : null;
    const selectedMoves =
      selectedSquareString && legalMoves ? (legalMoves[selectedSquareString] ?? []) : [];
    const selectedMoveTargets = new Set(selectedMoves.map((move) => move.to));
    return { selectedMoves, selectedMoveTargets };
  }, [selected, legalMoves]);

  function getResultString() {
    if (!matchResult) return null;

    const outcomeText = OUTCOME_TEXT[matchResult.outcome];
    const reasonText = REASON_TEXT[matchResult.reason];

    return `${outcomeText} ${reasonText}`;
  }

  const resultString = getResultString();

  function hasLegalMoves(square: Square) {
    const squareString = coordsToString(square);
    return (legalMoves?.[squareString]?.length ?? 0) > 0;
  }

  function tryMove(from: Square, to: Square, sourceColor: Color) {
    const fromStr = coordsToString(from);
    const toStr = coordsToString(to);
    const move = selectedMoves.find((m) => m.to === toStr);
    if (!move) {
      setSelected(null);
      return;
    }

    if (move.promotion) {
      setPromotionDialog({ from: fromStr, to: toStr, color: sourceColor });
    } else {
      sendMoveMessage(fromStr, toStr, null);
      setSelected(null);
    }
  }

  function handlePromotion(pieceType: PieceType) {
    if (promotionDialog) {
      sendMoveMessage(promotionDialog.from, promotionDialog.to, pieceType);
      setPromotionDialog(null);
      setSelected(null);
    }
  }

  function handleClick(square: Square, isMoveTarget: boolean) {
    const piece = board[square.file][square.rank];
    const isSameSquare = selected && selected.file === square.file && selected.rank === square.rank;

    if (selected && isMoveTarget) {
      const sourcePiece = board[selected.file][selected.rank];
      tryMove(selected, square, sourcePiece!.color);
    } else if (isSameSquare) {
      setSelected(null); // deselect if clicking the same square again
    } else if (piece && hasLegalMoves(square)) {
      setSelected(square); // nothing selected, select piece with legal moves
    } else {
      setSelected(null);
    }
  }

  function handleDragStart(event: DragStartEvent) {
    const sourceSquare = event.operation.source?.data.square;
    if (!sourceSquare || !hasLegalMoves(sourceSquare)) {
      setSelected(null);
      return;
    }

    setSelected(sourceSquare);
  }

  function handleDragEnd(event: DragEndEvent) {
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
    <div className="flex flex-col gap-2.5">
      <DragDropProvider onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
        <div className="relative grid grid-cols-8 grid-rows-8 max-w-[80vh] shadow-xl">
          {Array.from({ length: 64 }).map((_, index) => {
            const square = displayIndexToSquare(index, color);

            const piece = board[square.file][square.rank];
            const imgPath = piece ? `/pieces/${piece.color}-${piece.type}.png` : null;
            const isSelected = selected?.file === square.file && selected?.rank === square.rank;
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
              className={`
                            absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2
                            z-10 p-5
                            bg-gray-200 border-2 border-blue-600 rounded
                            text-blue-600 text-lg font-bold text-center
                            flex flex-col gap-3 items-center
                        `}
            >
              <div>{resultString}</div>
              <Button onClick={onBackToMenu} variant="contained">
                Back to Menu
              </Button>
            </div>
          )}
        </div>
      </DragDropProvider>
      <Dialog open={!!promotionDialog} onClose={() => setPromotionDialog(null)}>
        <DialogContent>
          <div className="flex gap-2">
            {PROMOTION_PIECES.map((pieceType) => (
              <div
                key={pieceType}
                onClick={() => handlePromotion(pieceType)}
                className={`
                                    cursor-pointer border-2 border-blue-600 rounded bg-gray-100
                                    w-[15vw] h-[15vw] max-w-[100px] max-h-[100px]
                                `}
              >
                <img
                  src={`/pieces/${promotionDialog?.color}-${pieceType}.png`}
                  className="w-full h-full"
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
