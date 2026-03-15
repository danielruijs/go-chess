import type { TimeFormat } from "../interfaces/chess";

function formatMatchTime(ms: number): string {
  const minutes = Math.floor(ms / 60000);
  const seconds = Math.floor((ms % 60000) / 1000);

  if (ms > 0 && minutes === 0 && seconds < 10) {
    const tenths = Math.floor((ms % 1000) / 100);
    return `${minutes}:${seconds.toString().padStart(2, "0")}.${tenths}`;
  }

  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
}

function formatTimeFormat(timeFormat: TimeFormat): string {
  const initialMinutes = timeFormat.initialMs / 60000;
  const incrementSeconds = timeFormat.incrementMs / 1000;
  return `${initialMinutes}+${incrementSeconds}`;
}

export { formatMatchTime, formatTimeFormat };
