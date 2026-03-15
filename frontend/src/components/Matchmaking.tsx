import type { QueueData } from "../interfaces/message";
import type { MouseEventHandler } from "react";
import { formatTimeFormat } from "../utils/time";
import type { TimeFormat } from "../interfaces/chess";

type MatchmakingComponentProps = {
  timeFormat: TimeFormat;
  queueData: QueueData | undefined;
  playerName: string;
  isConnected: boolean;
  inMatch: boolean;
  onJoinQueue: MouseEventHandler<HTMLDivElement>;
};

function MatchmakingComponent({
  timeFormat,
  queueData,
  playerName,
  isConnected,
  inMatch,
  onJoinQueue,
}: MatchmakingComponentProps) {
  const queueLength = queueData?.queueLength ?? 0;
  const inQueue = queueData?.inQueue ?? false;
  const disabled = !playerName || !isConnected || inQueue || inMatch;

  return (
    <div className="flex justify-center items-center relative">
      {inQueue && (
        <div className="absolute top-2 right-2 bg-green-500 text-white px-2 py-1 rounded-full text-xs font-semibold flex items-center gap-1">
          <span className="w-2 h-2 bg-white rounded-full animate-pulse"></span>
          In Queue
        </div>
      )}

      <div
        onClick={onJoinQueue}
        className={[
          "w-60 p-5 rounded-2xl backdrop-blur-xl border text-slate-800",
          "flex flex-col items-center gap-4 transition shadow-xl shadow-slate-400/40 relative",
          disabled
            ? "bg-white/50 border-white/60 cursor-not-allowed opacity-70"
            : "bg-white/70 border-white/80 cursor-pointer hover:shadow-2xl hover:-translate-y-1 hover:bg-white/80",
        ].join(" ")}
      >
        <p className="text-3xl font-bold">{formatTimeFormat(timeFormat)}</p>
        <p className="text-sm text-slate-600">Queue length: {queueLength}</p>
      </div>
    </div>
  );
}

export default MatchmakingComponent;
