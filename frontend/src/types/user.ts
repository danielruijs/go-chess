import type { Color, TimeFormat } from "./chess";
import type { Result } from "./result";

type UserInfo = {
  username: string;
  displayName: string;
  createdAt: string;
};

type GameRecord = {
  wins: number;
  losses: number;
  draws: number;
};

type UserStats = {
  white: GameRecord;
  black: GameRecord;
};

type UserMatchItem = {
  publicId: string;
  playedColor: Color;
  opponentDisplayName: string;
  opponentUsername?: string;
  result: Result;
  timeFormat: TimeFormat;
  moveCount: number;
  createdAt: string;
};

type UserProfile = {
  user: UserInfo;
  stats: UserStats;
  matches: UserMatchItem[];
};

export type { UserInfo, GameRecord, UserStats, UserMatchItem, UserProfile };
