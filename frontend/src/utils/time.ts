import type { TimeFormat } from "../types/chess";
import { pluralize } from "./text";

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

function formatBaseTime(timeFormat: TimeFormat): string {
  const baseMinutes = Math.floor(timeFormat.initialMs / 60000);
  const baseSeconds = Math.floor((timeFormat.initialMs % 60000) / 1000);

  if (baseMinutes > 0 && baseSeconds === 0) {
    return pluralize(baseMinutes, "minute", "minutes", true);
  }
  if (baseMinutes === 0 && baseSeconds > 0) {
    return pluralize(baseSeconds, "second", "seconds", true);
  }
  if (baseMinutes > 0 && baseSeconds > 0) {
    return `${baseMinutes}m ${baseSeconds}s`;
  }
  return "0 seconds";
}

function formatIncrement(timeFormat: TimeFormat): string {
  const incrementSeconds = Math.floor(timeFormat.incrementMs / 1000);
  if (incrementSeconds > 0) {
    return pluralize(incrementSeconds, "second", "seconds", true);
  }
  return "None";
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

export {
  formatMatchTime,
  formatTimeFormat,
  formatJoinDate,
  formatMatchDate,
  formatBaseTime,
  formatIncrement,
};
