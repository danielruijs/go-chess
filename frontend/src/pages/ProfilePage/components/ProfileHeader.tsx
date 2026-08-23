import UserAvatar from "../../../components/UserAvatar";
import { formatJoinDate } from "../../../utils/time";
import { getPct } from "../../../utils/math";
import type { UserInfo, UserStats } from "../../../types/user";

interface ProfileHeaderProps {
  user: UserInfo;
  stats: UserStats;
}

function ProfileHeader({ user, stats }: ProfileHeaderProps) {
  const totalWins = stats.white.wins + stats.black.wins;
  const totalDraws = stats.white.draws + stats.black.draws;
  const totalLosses = stats.white.losses + stats.black.losses;
  const totalGames = totalWins + totalDraws + totalLosses;
  const winPct = getPct(totalWins, totalGames);
  return (
    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 flex flex-col sm:flex-row items-center justify-between gap-6">
      <div className="flex flex-col sm:flex-row items-center gap-5 text-center sm:text-left">
        <UserAvatar name={user.displayName} size="lg" />

        <div className="flex flex-col gap-1">
          <h1 className="text-2xl font-extrabold text-slate-900">{user.displayName}</h1>
          <p className="text-slate-500 font-medium text-sm">@{user.username}</p>
          <p className="text-xs text-slate-400">Member since {formatJoinDate(user.createdAt)}</p>
        </div>
      </div>

      <div className="flex items-center gap-3 bg-slate-50 border border-slate-200 rounded-xl px-5 py-3 text-center">
        <div>
          <div className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Win Rate
          </div>
          <div className="text-2xl font-black text-slate-900">{Math.round(winPct * 10) / 10}%</div>
        </div>
        <div className="h-8 w-px bg-slate-200 mx-1" />
        <div>
          <div className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Played
          </div>
          <div className="text-2xl font-black text-slate-900">{totalGames}</div>
        </div>
      </div>
    </div>
  );
}

export default ProfileHeader;
