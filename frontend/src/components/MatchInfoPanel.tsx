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
    <div className="p-2 border border-gray-600 h-72 w-48 flex flex-col box-border">
      <MoveList groupedMoves={groupedMoves} />

      <div className="mt-2">
        {isDrawOfferPending ? (
          <div className="flex flex-col items-center gap-2">
            <div className="text-base font-semibold">Draw offered</div>
            <div className="flex justify-center gap-2">
              <Button
                sx={{ width: "40px", height: "30px", fontSize: "22px" }}
                variant="contained"
                disabled={!inMatch}
                color="success"
                onClick={() => onRespondToDrawOffer(true)}
              >
                ✓
              </Button>
              <Button
                sx={{ width: "40px", height: "30px", fontSize: "22px" }}
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
          <div className="flex justify-center gap-2">
            <Button
              sx={{ width: "100px" }}
              variant="contained"
              disabled={!inMatch}
              onClick={onResign}
            >
              Resign
            </Button>
            <Button
              sx={{ width: "100px" }}
              variant="contained"
              disabled={!inMatch || isDrawOfferSent}
              onClick={onOfferDraw}
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
