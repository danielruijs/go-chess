import type { QueueData } from "../types/message";
import type { MouseEventHandler } from "react";
import { formatBaseTime, formatIncrement, formatTimeFormat } from "../utils/time";
import type { TimeFormat } from "../types/chess";
import { IconButton, Tooltip } from "@mui/material";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";

type MatchmakingComponentProps = {
  timeFormat: TimeFormat;
  queueData: QueueData | undefined;
  displayName: string;
  isConnected: boolean;
  inMatch: boolean;
  onToggleQueue: MouseEventHandler<HTMLDivElement>;
};

function MatchmakingComponent({
  timeFormat,
  queueData,
  displayName,
  isConnected,
  inMatch,
  onToggleQueue,
}: MatchmakingComponentProps) {
  const queueLength = queueData?.queueLength ?? 0;
  const inQueue = queueData?.inQueue ?? false;
  const disabled = !displayName || !isConnected || inMatch;
  const actionText = inQueue ? "Click to leave queue" : "Click to join queue";

  const baseTimeLabel = formatBaseTime(timeFormat);
  const incrementLabel = formatIncrement(timeFormat);

  const tooltipTitle = (
    <div className="text-xs p-1 space-y-1">
      <div className="font-semibold text-white mb-1.5">
        Time Control ({formatTimeFormat(timeFormat)})
      </div>

      <div className="flex justify-between gap-3">
        <span className="text-slate-300">Starting clock</span>
        <span className="font-medium text-white">{baseTimeLabel} / player</span>
      </div>

      <div className="flex justify-between gap-3">
        <span className="text-slate-300">Increment</span>
        <span className="font-medium text-white">
          {incrementLabel === "None" ? "None" : `${incrementLabel} / move`}
        </span>
      </div>
    </div>
  );

  return (
    <div className="flex justify-center items-center relative">
      {inQueue && (
        <div className="absolute top-2 left-2 bg-green-500 text-white px-2 py-1 rounded-full text-xs font-semibold flex items-center gap-1 z-10">
          <span className="w-2 h-2 bg-white rounded-full animate-pulse"></span>
          In Queue
        </div>
      )}

      <div className="absolute top-2 right-2 z-10">
        <Tooltip
          title={tooltipTitle}
          arrow
          placement="top"
          enterTouchDelay={0}
          leaveTouchDelay={60000}
        >
          <IconButton
            size="small"
            aria-label="Time control info"
            onClick={(e) => e.stopPropagation()}
          >
            <InfoOutlinedIcon sx={{ fontSize: 18 }} />
          </IconButton>
        </Tooltip>
      </div>

      <div
        onClick={onToggleQueue}
        className={[
          "w-60 p-5 rounded-2xl backdrop-blur-xl border text-slate-800",
          "flex flex-col items-center gap-3 transition shadow-xl shadow-slate-400/40 relative",
          disabled
            ? "bg-white/50 border-white/60 cursor-not-allowed opacity-70"
            : inQueue
              ? "bg-green-50/80 border-green-300 cursor-pointer hover:bg-green-100/80"
              : "bg-white/70 border-white/80 cursor-pointer hover:shadow-2xl hover:-translate-y-1 hover:bg-white/80",
        ].join(" ")}
      >
        <p className="text-3xl font-bold">{formatTimeFormat(timeFormat)}</p>
        <p className="text-sm text-slate-500">{actionText}</p>
        <p className="text-sm text-slate-600">Queue length: {queueLength}</p>
      </div>
    </div>
  );
}

export default MatchmakingComponent;
