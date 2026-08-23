import { pluralize } from "../../../utils/text";
import { getPct } from "../../../utils/math";
import type { GameRecord } from "../../../types/user";
import type { Color } from "../../../types/chess";

interface ColorPerformanceCardProps {
  color: Color;
  record: GameRecord;
}

function ColorPerformanceCard({ color, record }: ColorPerformanceCardProps) {
  const isWhite = color === "white";
  const games = record.wins + record.draws + record.losses;
  const winPct = getPct(record.wins, games);

  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 flex flex-col justify-between">
      <div>
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2 text-slate-700 font-semibold text-sm">
            <span
              className={`w-3.5 h-3.5 rounded-full shrink-0 ${
                isWhite ? "bg-white border-2 border-slate-700" : "bg-slate-900"
              }`}
            />
            <span>{isWhite ? "As White" : "As Black"}</span>
          </div>
          <span className="text-xs font-bold text-slate-500 bg-slate-100 px-2 py-0.5 rounded-full">
            {pluralize(games, "game", "games", true)}
          </span>
        </div>

        <div className="my-2">
          <div className="flex justify-between items-baseline mb-1">
            <span className="text-sm text-slate-500 font-medium">Win Rate</span>
            <span className="text-lg font-bold text-slate-900">
              {Math.round(winPct * 10) / 10}%
            </span>
          </div>
          <div className="bg-slate-100 rounded-full h-2 overflow-hidden">
            <div
              className="h-2 rounded-full bg-slate-800"
              style={{ width: `${Math.min(100, winPct)}%` }}
            />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-2 text-center text-xs mt-4 pt-4 border-t border-slate-100">
        <div>
          <span className="text-slate-400 block">Won</span>
          <span className="font-bold text-emerald-600">{record.wins}</span>
        </div>
        <div>
          <span className="text-slate-400 block">Drawn</span>
          <span className="font-bold text-slate-600">{record.draws}</span>
        </div>
        <div>
          <span className="text-slate-400 block">Lost</span>
          <span className="font-bold text-rose-600">{record.losses}</span>
        </div>
      </div>
    </div>
  );
}

export default ColorPerformanceCard;
