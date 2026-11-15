import { formatMatchTime } from "../utils/time";

type MatchClockProps = {
  timeMs: number;
  isActive: boolean;
};

function MatchClock({ timeMs, isActive }: MatchClockProps) {
  return (
    <div
      className={`flex justify-center items-center rounded-xl border font-bold tabular-nums transition-all duration-200 shadow-md px-3 py-1 text-sm md:py-3 md:text-2xl ${
        isActive
          ? "border-slate-300 bg-slate-100 text-slate-900 ring-2 ring-slate-300"
          : "border-slate-700 bg-slate-800 text-slate-400"
      }`}
    >
      {formatMatchTime(timeMs)}
    </div>
  );
}

export default MatchClock;
