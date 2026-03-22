import type { Color, MaterialDiff } from "../types/chess";
import MaterialDifference from "./MaterialDifference";

type PlayerInfoProps = {
  color: Color;
  name: string;
  materialDiff: MaterialDiff;
};

function PlayerInfo({ color, name, materialDiff }: PlayerInfoProps) {
  return (
    <div className="flex items-center gap-5">
      <span className="text-lg font-bold">{name}</span>
      <MaterialDifference materialDiff={materialDiff} color={color} />
    </div>
  );
}

export default PlayerInfo;
