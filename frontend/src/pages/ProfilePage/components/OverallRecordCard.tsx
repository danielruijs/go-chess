import { Tooltip } from "@mui/material";
import LeaderboardIcon from "@mui/icons-material/Leaderboard";
import { getPct } from "../../../utils/math";
import { pluralize } from "../../../utils/text";
import type { UserStats } from "../../../types/user";

interface OverallRecordCardProps {
  stats: UserStats;
}

function OverallRecordCard({ stats }: OverallRecordCardProps) {
  const totalWins = stats.white.wins + stats.black.wins;
  const totalDraws = stats.white.draws + stats.black.draws;
  const totalLosses = stats.white.losses + stats.black.losses;
  const totalGames = totalWins + totalDraws + totalLosses;
  const winPct = getPct(totalWins, totalGames);
  const drawPct = getPct(totalDraws, totalGames);
  const lossPct = getPct(totalLosses, totalGames);

  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 flex flex-col justify-between">
      <div>
        <div className="flex items-center gap-2 text-slate-500 font-semibold text-sm">
          <LeaderboardIcon fontSize="small" className="text-blue-600" />
          <span>Overall Record</span>
        </div>

        <div className="grid grid-cols-3 gap-2 text-center my-4">
          <div className="bg-emerald-50 border border-emerald-100 rounded-xl p-3">
            <div className="text-xl font-bold text-emerald-700">{totalWins}</div>
            <div className="text-xs font-semibold text-emerald-600 uppercase mt-0.5">
              {pluralize(totalWins, "Win", "Wins", false)}
            </div>
          </div>
          <div className="bg-slate-50 border border-slate-200 rounded-xl p-3">
            <div className="text-xl font-bold text-slate-700">{totalDraws}</div>
            <div className="text-xs font-semibold text-slate-600 uppercase mt-0.5">
              {pluralize(totalDraws, "Draw", "Draws", false)}
            </div>
          </div>
          <div className="bg-rose-50 border border-rose-100 rounded-xl p-3">
            <div className="text-xl font-bold text-rose-700">{totalLosses}</div>
            <div className="text-xs font-semibold text-rose-600 uppercase mt-0.5">
              {pluralize(totalLosses, "Loss", "Losses", false)}
            </div>
          </div>
        </div>
      </div>

      {totalGames > 0 ? (
        <div className="mt-2">
          <div className="h-3 w-full bg-slate-100 rounded-full flex overflow-hidden">
            {winPct > 0 && (
              <Tooltip title={`Wins: ${winPct.toFixed(1)}%`}>
                <div style={{ width: `${winPct}%` }} className="bg-emerald-500" />
              </Tooltip>
            )}
            {drawPct > 0 && (
              <Tooltip title={`Draws: ${drawPct.toFixed(1)}%`}>
                <div style={{ width: `${drawPct}%` }} className="bg-slate-400" />
              </Tooltip>
            )}
            {lossPct > 0 && (
              <Tooltip title={`Losses: ${lossPct.toFixed(1)}%`}>
                <div style={{ width: `${lossPct}%` }} className="bg-rose-500" />
              </Tooltip>
            )}
          </div>
          <div className="flex justify-between items-center text-xs font-semibold mt-2 px-1">
            <div className="flex items-center gap-1.5 text-emerald-600">
              <span className="w-2 h-2 rounded-full bg-emerald-500 shrink-0" />
              <span>{winPct.toFixed(0)}% W</span>
            </div>
            <div className="flex items-center gap-1.5 text-slate-500">
              <span className="w-2 h-2 rounded-full bg-slate-400 shrink-0" />
              <span>{drawPct.toFixed(0)}% D</span>
            </div>
            <div className="flex items-center gap-1.5 text-rose-600">
              <span className="w-2 h-2 rounded-full bg-rose-500 shrink-0" />
              <span>{lossPct.toFixed(0)}% L</span>
            </div>
          </div>
        </div>
      ) : (
        <p className="text-xs text-slate-400 italic text-center py-2">No matches played yet</p>
      )}
    </div>
  );
}

export default OverallRecordCard;
