import type { Color, MaterialDiff } from "../types/chess";
import MaterialDifference from "./MaterialDifference";
import MatchClock from "./MatchClock";

type PlayerInfoProps = {
  color: Color;
  name: string;
  materialDiff: MaterialDiff;
  timeMs?: number;
  isActive?: boolean;
};

function PlayerInfo({ color, name, materialDiff, timeMs, isActive = false }: PlayerInfoProps) {
  return (
    <div className="flex items-center justify-between w-full px-1 py-0.5">
      <div className="flex items-center gap-2 sm:gap-4 min-w-0">
        <span className="text-base sm:text-lg font-bold truncate text-slate-900">{name}</span>
        <MaterialDifference materialDiff={materialDiff} color={color} />
      </div>
      {timeMs !== undefined && (
        <div className="md:hidden flex-shrink-0">
          <MatchClock timeMs={timeMs} isActive={isActive} />
        </div>
      )}
    </div>
  );
}

export default PlayerInfo;
