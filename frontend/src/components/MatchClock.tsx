import { formatTime } from "../utils/chess";

function MatchClock({ timeMs, isActive }: { timeMs: number; isActive: boolean }) {
    return (
        <div
            className={[
                "flex justify-center rounded-xl border py-3 transition-all duration-200 shadow-md",
                isActive
                    ? "border-slate-300 bg-slate-100 text-slate-900 ring-2 ring-slate-300"
                    : "border-slate-700 bg-slate-800 text-slate-400",
            ].join(" ")}
        >
            <div className="text-2xl font-bold tabular-nums">{formatTime(timeMs)}</div>
        </div>
    );
}

export default MatchClock;