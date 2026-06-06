import { createContext, useContext } from "react";
import type { Credentials } from "../types/auth";
import type { PlayerInfoData } from "../types/message";

export interface AuthContextType {
  playerInfo: PlayerInfoData | null;
  isLoading: boolean;
  login: (credentials: Credentials) => Promise<Error | null>;
  register: (credentials: Credentials) => Promise<Error | null>;
  logout: () => Promise<void>;
  setPlayerInfo: (playerInfo: PlayerInfoData | null) => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
