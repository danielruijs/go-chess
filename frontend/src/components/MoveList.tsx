import { useEffect, useRef } from "react";

type MoveItem = {
  san: string;
  positionIndex: number;
};

export type MoveRow = {
  moveNumber: number;
  white: MoveItem;
  black: MoveItem | null;
};

type MoveListProps = {
  groupedMoves: MoveRow[];
  currentPosition?: number;
  onSelectPosition?: (positionIndex: number) => void;
};

function MoveList({ groupedMoves, currentPosition, onSelectPosition }: MoveListProps) {
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);
  const selectedMoveCellRef = useRef<HTMLTableCellElement | null>(null);

  useEffect(() => {
    // Scroll back to the top when navigating to the starting position
    if (currentPosition === 0 && scrollContainerRef.current) {
      scrollContainerRef.current.scrollTo({ top: 0, behavior: "smooth" });
      // Scroll the selected move into view if it is partially or fully outside the visible area
    } else if (selectedMoveCellRef.current) {
      selectedMoveCellRef.current.scrollIntoView({
        block: "nearest",
        inline: "nearest",
        behavior: "smooth",
      });
      // Live match mode: auto-scroll to bottom when a new move is made
    } else if (currentPosition === undefined && scrollContainerRef.current) {
      scrollContainerRef.current.scrollTo({
        top: scrollContainerRef.current.scrollHeight,
        behavior: "smooth",
      });
    }
  }, [currentPosition, groupedMoves]);

  const renderCell = (move: MoveItem | null) => {
    if (!move) {
      return <td className="px-1" />;
    }

    if (!onSelectPosition) {
      return <td className="px-1">{move.san}</td>;
    }

    const isSelected = currentPosition !== undefined && currentPosition === move.positionIndex;

    return (
      <td
        ref={isSelected ? selectedMoveCellRef : null} // Attach ref only to the currently active move
        onClick={isSelected ? undefined : () => onSelectPosition(move.positionIndex)}
        className={
          isSelected
            ? "bg-blue-300 font-bold px-1 rounded cursor-default"
            : "cursor-pointer px-1 rounded hover:bg-blue-100"
        }
      >
        {move.san}
      </td>
    );
  };

  return (
    <div ref={scrollContainerRef} className="flex-1 overflow-y-auto">
      <table className="w-full border-collapse">
        <thead>
          <tr className="border-b border-gray-400">
            <th className="text-left w-1/6 pl-1">#</th>
            <th className="text-left w-2/5 pl-1">White</th>
            <th className="text-left w-2/5 pl-1">Black</th>
          </tr>
        </thead>
        <tbody>
          {groupedMoves.map(({ moveNumber, white, black }, index) => (
            <tr key={moveNumber} className={index % 2 === 0 ? "bg-gray-200" : ""}>
              <td className="pl-1 text-gray-500">{moveNumber}.</td>
              {renderCell(white)}
              {renderCell(black)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default MoveList;
