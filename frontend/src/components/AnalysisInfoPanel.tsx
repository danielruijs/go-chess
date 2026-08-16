import { IconButton, Tooltip } from "@mui/material";
import SkipPreviousIcon from "@mui/icons-material/SkipPrevious";
import NavigateBeforeIcon from "@mui/icons-material/NavigateBefore";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";
import SkipNextIcon from "@mui/icons-material/SkipNext";
import FlipIcon from "@mui/icons-material/Flip";
import { getResultString, type Result } from "../types/result";
import MoveList, { type MoveRow } from "./MoveList";

type AnalysisInfoPanelProps = {
  groupedMoves: MoveRow[];
  currentPosition: number;
  totalPositions: number;
  result: Result | null;
  onSelectPosition: (positionIndex: number) => void;
  onFirst: () => void;
  onPrev: () => void;
  onNext: () => void;
  onLast: () => void;
  onFlip: () => void;
};

function AnalysisInfoPanel({
  groupedMoves,
  currentPosition,
  totalPositions,
  result,
  onSelectPosition,
  onFirst,
  onPrev,
  onNext,
  onLast,
  onFlip,
}: AnalysisInfoPanelProps) {
  const isFirst = currentPosition === 0;
  const isLast = currentPosition === totalPositions - 1;

  const resultString = getResultString(result);

  return (
    <div className="p-2 border border-gray-600 h-72 w-48 flex flex-col box-border">
      <MoveList
        groupedMoves={groupedMoves}
        currentPosition={currentPosition}
        onSelectPosition={onSelectPosition}
      />

      {resultString && (
        <div className="text-center text-xs font-semibold py-1 text-blue-700 border-y border-gray-300">
          {resultString}
        </div>
      )}

      <div className="flex justify-center items-center gap-1 border-gray-300 pt-1">
        <Tooltip title="First position">
          <span>
            <IconButton size="small" onClick={onFirst} disabled={isFirst}>
              <SkipPreviousIcon fontSize="small" />
            </IconButton>
          </span>
        </Tooltip>
        <Tooltip title="Previous position">
          <span>
            <IconButton size="small" onClick={onPrev} disabled={isFirst}>
              <NavigateBeforeIcon fontSize="small" />
            </IconButton>
          </span>
        </Tooltip>
        <Tooltip title="Next position">
          <span>
            <IconButton size="small" onClick={onNext} disabled={isLast}>
              <NavigateNextIcon fontSize="small" />
            </IconButton>
          </span>
        </Tooltip>
        <Tooltip title="Last position">
          <span>
            <IconButton size="small" onClick={onLast} disabled={isLast}>
              <SkipNextIcon fontSize="small" />
            </IconButton>
          </span>
        </Tooltip>
        <Tooltip title="Flip board">
          <IconButton size="small" onClick={onFlip}>
            <FlipIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </div>
    </div>
  );
}

export default AnalysisInfoPanel;
