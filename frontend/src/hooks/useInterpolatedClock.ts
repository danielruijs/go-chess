import { useEffect, useState } from "react";
import type { Color } from "../types/chess";
import type { ClockData } from "../types/message";

type UseInterpolatedClockProps = {
  clock: ClockData | null;
  activeColor: Color | null;
  inMatch: boolean;
  intervalMs: number;
};

function useInterpolatedClock({
  clock,
  activeColor,
  inMatch,
  intervalMs,
}: UseInterpolatedClockProps): ClockData | null {
  const [displayedClock, setDisplayedClock] = useState<ClockData | null>(clock);
  const [prevClock, setPrevClock] = useState<ClockData | null>(clock);

  if (clock !== prevClock) {
    setPrevClock(clock);
    setDisplayedClock(clock);
  }

  useEffect(() => {
    if (!clock || !inMatch || !activeColor) return;

    const startTime = Date.now();

    const updateClock = () => {
      const elapsed = Date.now() - startTime;

      setDisplayedClock({
        whiteTimeMs:
          activeColor === "white" ? Math.max(0, clock.whiteTimeMs - elapsed) : clock.whiteTimeMs,
        blackTimeMs:
          activeColor === "black" ? Math.max(0, clock.blackTimeMs - elapsed) : clock.blackTimeMs,
        incrementMs: clock.incrementMs,
      });
    };

    const intervalId = setInterval(updateClock, intervalMs);

    return () => clearInterval(intervalId);
  }, [clock, activeColor, inMatch, intervalMs]);

  return displayedClock;
}

export default useInterpolatedClock;
