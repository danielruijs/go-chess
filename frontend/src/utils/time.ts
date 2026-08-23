import type { TimeFormat } from "../types/chess";

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

function formatJoinDate(dateStr: string): string {
  try {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return dateStr;
    return date.toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  } catch {
    return dateStr;
  }
}

function formatMatchDate(dateStr: string): string {
  try {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return dateStr;

    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60 / 1000);
    const diffHours = Math.floor(diffMs / 60 / 60 / 1000);
    const diffDays = Math.floor(diffMs / 24 / 60 / 60 / 1000);

    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;

    return date.toLocaleDateString(undefined, {
      year: date.getFullYear() !== now.getFullYear() ? "numeric" : undefined,
      month: "short",
      day: "numeric",
    });
  } catch {
    return dateStr;
  }
}

export { formatMatchTime, formatTimeFormat, formatJoinDate, formatMatchDate };
