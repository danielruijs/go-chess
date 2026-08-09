import { useEffect, useMemo, useCallback, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import BoardComponent from "../components/Board";
import MatchClock from "../components/MatchClock";
import AnalysisInfoPanel from "../components/AnalysisInfoPanel";
import PlayerInfo from "../components/PlayerInfo";
import { fetchMatch } from "../api/match";
import { getMaterialDiff } from "../utils/chess";
import type { Match } from "../types/match";
import type { Color } from "../types/chess";
import type { MoveRow } from "../components/MoveList";

function Analysis() {
  const navigate = useNavigate();
  const { publicId } = useParams<{ publicId: string }>();

  const [matchData, setMatchData] = useState<Match | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPosition, setCurrentPosition] = useState(0);
  const [boardOrientation, setBoardOrientation] = useState<Color>("white");

  // Load analysis data on mount
  useEffect(() => {
    if (!publicId) return;

    setLoading(true);
    setError(null);

    fetchMatch(publicId)
      .then((res) => {
        switch (res.status) {
          case "in-progress":
            navigate(`/match/${publicId}`, { replace: true });
            break;
          case "not-found":
            setError("Match not found.");
            break;
          case "ok":
            setMatchData(res.data);
            setCurrentPosition(0);
            break;
        }
      })
      .catch(() => setError("Failed to load match data."))
      .finally(() => setLoading(false));
  }, [publicId, navigate]);

  // Navigation handlers
  const goFirst = useCallback(() => setCurrentPosition(0), []);
  const goPrev = useCallback(() => setCurrentPosition((p) => Math.max(0, p - 1)), []);
  const goNext = useCallback(
    () =>
      setCurrentPosition((p) => (matchData ? Math.min(matchData.positions.length - 1, p + 1) : p)),
    [matchData]
  );
  const goLast = useCallback(
    () => setCurrentPosition(matchData ? matchData.positions.length - 1 : 0),
    [matchData]
  );

  // Keyboard navigation
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === "ArrowLeft") {
        e.preventDefault();
        goPrev();
      } else if (e.key === "ArrowRight") {
        e.preventDefault();
        goNext();
      } else if (e.key === "Home") {
        e.preventDefault();
        goFirst();
      } else if (e.key === "End") {
        e.preventDefault();
        goLast();
      }
    }
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [goPrev, goNext, goFirst, goLast]);

  const currentPositionData = matchData?.positions[currentPosition];

  const groupedMoves = useMemo<MoveRow[]>(() => {
    if (!matchData) return [];
    const movePositions = matchData.positions.slice(1); // skip starting position
    const rows: MoveRow[] = [];

    for (let i = 0; i < movePositions.length; i += 2) {
      const whitePos = movePositions[i];
      const blackPos = movePositions[i + 1];
      rows.push({
        moveNumber: Math.floor(i / 2) + 1,
        white: { san: whitePos.san ?? "", positionIndex: whitePos.index },
        black: blackPos ? { san: blackPos.san ?? "", positionIndex: blackPos.index } : null,
      });
    }
    return rows;
  }, [matchData]);

  const materialDiff = useMemo(
    () => getMaterialDiff(currentPositionData?.board ?? null),
    [currentPositionData]
  );

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center text-lg">Loading analysis...</div>
    );
  }

  if (error || !matchData || !currentPositionData) {
    return (
      <div className="flex-1 flex items-center justify-center text-lg text-red-600">
        {error ?? "Something went wrong."}
      </div>
    );
  }

  const { whitePlayerName, blackPlayerName, result } = matchData;

  const bottomPlayer = boardOrientation; // "white" or "black"
  const topPlayer = bottomPlayer === "white" ? "black" : "white";
  const bottomPlayerName = bottomPlayer === "white" ? whitePlayerName : blackPlayerName;
  const topPlayerName = topPlayer === "white" ? whitePlayerName : blackPlayerName;

  const topTimeMs =
    topPlayer === "white" ? currentPositionData.whiteTimeMs : currentPositionData.blackTimeMs;
  const bottomTimeMs =
    bottomPlayer === "white" ? currentPositionData.whiteTimeMs : currentPositionData.blackTimeMs;

  return (
    <div className="flex-1 flex items-center">
      <div className="flex gap-4 items-center justify-center w-full">
        <div className="flex flex-col gap-1">
          <PlayerInfo color={topPlayer} name={topPlayerName} materialDiff={materialDiff} />

          <div className="flex gap-4">
            <BoardComponent
              color={boardOrientation}
              board={currentPositionData.board}
              legalMoves={null}
              matchResult={null}
              sendMoveMessage={() => {}}
              onBackToMenu={() => navigate("/")}
            />

            <div className="flex flex-col justify-between">
              <MatchClock timeMs={topTimeMs} isActive={false} />
              <AnalysisInfoPanel
                groupedMoves={groupedMoves}
                currentPosition={currentPosition}
                totalPositions={matchData.positions.length}
                result={result}
                onSelectPosition={(positionIndex) => setCurrentPosition(positionIndex)}
                onFirst={goFirst}
                onPrev={goPrev}
                onNext={goNext}
                onLast={goLast}
                onFlip={() => setBoardOrientation((o) => (o === "white" ? "black" : "white"))}
              />
              <MatchClock timeMs={bottomTimeMs} isActive={false} />
            </div>
          </div>

          <PlayerInfo color={bottomPlayer} name={bottomPlayerName} materialDiff={materialDiff} />
        </div>
      </div>
    </div>
  );
}

export default Analysis;
