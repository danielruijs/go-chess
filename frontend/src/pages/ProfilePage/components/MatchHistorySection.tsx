import HistoryIcon from "@mui/icons-material/History";
import SportsEsportsIcon from "@mui/icons-material/SportsEsports";
import MatchRow from "./MatchRow";
import type { UserMatchItem } from "../../../types/user";

interface MatchHistorySectionProps {
  matches: UserMatchItem[];
}

function MatchHistorySection({ matches }: MatchHistorySectionProps) {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <HistoryIcon className="text-go" />
        <h2 className="text-xl font-bold text-slate-900 tracking-tight">Match History</h2>
        <span className="bg-slate-100 text-slate-600 text-xs font-bold px-2 py-0.5 rounded-full">
          {matches.length}
        </span>
      </div>

      {matches.length === 0 ? (
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-12 text-center flex flex-col items-center">
          <div className="w-12 h-12 rounded-full bg-slate-100 text-slate-400 flex items-center justify-center mb-3">
            <SportsEsportsIcon />
          </div>
          <h3 className="font-bold text-slate-800">No Matches Played Yet</h3>
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {matches.map((match) => (
            <MatchRow key={match.publicId} match={match} />
          ))}
        </div>
      )}
    </div>
  );
}

export default MatchHistorySection;
