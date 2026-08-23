import { Link } from "react-router-dom";
import { Button, Tooltip } from "@mui/material";
import AccessTimeIcon from "@mui/icons-material/AccessTime";
import AnalyticsOutlinedIcon from "@mui/icons-material/AnalyticsOutlined";
import { formatMatchDate, formatTimeFormat } from "../../../utils/time";
import { getPlayerResult } from "../../../utils/result";
import { pluralize } from "../../../utils/text";
import { REASON_TEXT } from "../../../types/result";
import type { UserMatchItem } from "../../../types/user";

interface MatchRowProps {
  match: UserMatchItem;
}

function MatchRow({ match }: MatchRowProps) {
  const playerResult = getPlayerResult(match.result, match.playedColor);
  const isWin = playerResult === "win";
  const isLoss = playerResult === "loss";
  const asWhite = match.playedColor === "white";

  return (
    <div className="bg-white rounded-xl border border-slate-200 hover:border-blue-300 hover:shadow-md transition p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      {/* Left: Outcome badge + Opponent info */}
      <div className="flex items-center gap-4">
        <div
          className={`w-16 shrink-0 text-center py-1.5 rounded-lg font-black text-sm capitalize select-none ${
            isWin
              ? "bg-emerald-100 text-emerald-800 border border-emerald-200"
              : isLoss
                ? "bg-rose-100 text-rose-800 border border-rose-200"
                : "bg-slate-100 text-slate-700 border border-slate-200"
          }`}
        >
          {playerResult}
        </div>

        <div className="flex flex-col">
          <div className="flex items-center gap-2">
            <Tooltip title={asWhite ? "Played as White" : "Played as Black"}>
              <span
                className={`w-3.5 h-3.5 rounded-full shrink-0 ${
                  asWhite ? "bg-white border-2 border-slate-700" : "bg-slate-900"
                }`}
              />
            </Tooltip>
            <span className="text-slate-500 text-xs font-medium">vs</span>
            {match.opponentUsername ? (
              <Link
                to={`/user/${match.opponentUsername}`}
                className="font-bold text-slate-900 text-sm hover:text-blue-600 hover:underline transition-colors"
                title={`View ${match.opponentDisplayName}'s profile`}
              >
                {match.opponentDisplayName}
              </Link>
            ) : (
              <span className="font-bold text-slate-900 text-sm">{match.opponentDisplayName}</span>
            )}
          </div>
          <span className="text-xs text-slate-400 mt-0.5">
            {REASON_TEXT[match.result.reason]} • {pluralize(match.moveCount, "move", "moves")}
          </span>
        </div>
      </div>

      {/* Right: Time control, Date, Analysis button */}
      <div className="flex items-center justify-between gap-5">
        <div className="flex items-center gap-4 text-xs text-slate-500 font-medium">
          <span className="flex items-center gap-1 bg-slate-50 px-2.5 py-1 rounded-md border border-slate-200">
            <AccessTimeIcon fontSize="inherit" className="text-slate-400" />
            {formatTimeFormat(match.timeFormat)}
          </span>
          <span>{formatMatchDate(match.createdAt)}</span>
        </div>

        <Link to={`/analysis/${match.publicId}`}>
          <Button
            variant="outlined"
            size="small"
            color="primary"
            startIcon={<AnalyticsOutlinedIcon fontSize="small" />}
          >
            Analyze
          </Button>
        </Link>
      </div>
    </div>
  );
}

export default MatchRow;
