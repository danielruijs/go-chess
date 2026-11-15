import { Button } from "@mui/material";
import MoveList from "./MoveList";
import type { MoveRow } from "../utils/chess";

type MatchInfoPanelProps = {
  groupedMoves: MoveRow[];
  isDrawOfferPending: boolean;
  isDrawOfferSent: boolean;
  inMatch: boolean;
  onRespondToDrawOffer: (accept: boolean) => void;
  onResign: () => void;
  onOfferDraw: () => void;
};

function MatchInfoPanel({
  groupedMoves,
  isDrawOfferPending,
  isDrawOfferSent,
  inMatch,
  onRespondToDrawOffer,
  onResign,
  onOfferDraw,
}: MatchInfoPanelProps) {
  return (
    <div className="w-full md:w-48 h-48 md:h-72 p-2 border border-gray-600 rounded flex flex-col box-border">
      <MoveList groupedMoves={groupedMoves} />

      <div className="mt-2 shrink-0">
        {isDrawOfferPending ? (
          <div className="flex flex-col items-center gap-1">
            <div className="text-sm font-semibold">Draw offered</div>
            <div className="flex justify-center gap-2 w-full">
              <Button
                sx={{
                  flex: 1,
                  minWidth: 0,
                  height: "32px",
                  fontSize: "18px",
                }}
                variant="contained"
                disabled={!inMatch}
                color="success"
                onClick={() => onRespondToDrawOffer(true)}
              >
                ✓
              </Button>
              <Button
                sx={{
                  flex: 1,
                  minWidth: 0,
                  height: "32px",
                  fontSize: "18px",
                }}
                variant="contained"
                disabled={!inMatch}
                color="error"
                onClick={() => onRespondToDrawOffer(false)}
              >
                ✕
              </Button>
            </div>
          </div>
        ) : (
          <div className="flex justify-center gap-2 w-full">
            <Button
              variant="contained"
              size="small"
              disabled={!inMatch}
              onClick={onResign}
              sx={{ flex: 1 }}
            >
              Resign
            </Button>
            <Button
              variant="contained"
              size="small"
              disabled={!inMatch || isDrawOfferSent}
              onClick={onOfferDraw}
              sx={{ flex: 1 }}
            >
              {isDrawOfferSent ? "Offered" : "Draw"}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

export default MatchInfoPanel;
