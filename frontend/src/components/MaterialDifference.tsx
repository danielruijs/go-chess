import type { Color, MaterialDiff, PieceType } from "../types/chess";

type MaterialDifferenceProps = {
  materialDiff: MaterialDiff;
  color: Color;
};

function MaterialDifference({ materialDiff, color }: MaterialDifferenceProps) {
  const extraPieces = materialDiff.extraPieces[color];
  const pieceTypes = Object.entries(extraPieces).flatMap(([type, count]) => {
    return Array.from({ length: count }, () => type as PieceType);
  });

  const oppositeColor = color === "white" ? "black" : "white";

  const score = Math.max(color === "white" ? materialDiff.score : -materialDiff.score, 0);

  return (
    <div className="flex items-center">
      <div className="flex">
        {pieceTypes.map((type, index) => {
          return (
            <img
              key={index}
              src={`/pieces/${oppositeColor}-${type}.png`}
              className={`
                h-[1.5em] w-auto object-contain pb-0.5
                ${type === "pawn" ? "ml-[-8px]" : "ml-[-2px]"}
            `}
            />
          );
        })}
      </div>
      {score > 0 && <span className="flex ml-1 text-sm font-bold text-gray-500">+{score}</span>}
    </div>
  );
}

export default MaterialDifference;
