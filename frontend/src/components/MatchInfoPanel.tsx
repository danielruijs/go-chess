import { Button } from "@mui/material";
type MoveRow = {
  moveNumber: number;
  whiteMove: string;
  blackMove: string;
};

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
      <div className="flex-1 overflow-y-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr>
              <th className="text-left w-1/6">#</th>
              <th className="text-left w-2/5">White</th>
              <th className="text-left w-2/5">Black</th>
            </tr>
          </thead>
          <tbody>
            {groupedMoves.map(({ moveNumber, whiteMove, blackMove }, index) => (
              <tr key={moveNumber} className={index % 2 === 0 ? "bg-gray-200" : ""}>
                <td>{moveNumber}.</td>
                <td>{whiteMove}</td>
                <td>{blackMove}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
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
