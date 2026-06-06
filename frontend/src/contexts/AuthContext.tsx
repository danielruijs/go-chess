import { useState, useEffect, type ReactNode } from "react";
import { AuthContext } from "./AuthContext";
import {
  checkSession as apiCheckSession,
  logout as apiLogout,
  login as apiLogin,
  register as apiRegister,
} from "../api/auth";
import type { Credentials } from "../types/auth";
import type { PlayerInfoData } from "../types/message";

export function AuthProvider({ children }: { children: ReactNode }) {
  const [playerInfo, setPlayerInfo] = useState<PlayerInfoData | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  async function checkSession() {
    try {
      const playerInfoData = await apiCheckSession();
      setPlayerInfo(playerInfoData);
    } catch (e) {
      console.error("Session check failed", e);
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    checkSession();
  }, []);

  async function login(credentials: Credentials) {
    const err = await apiLogin(credentials);
    if (!err) {
      await checkSession();
    }
    return err;
  }

  async function register(credentials: Credentials) {
    const err = await apiRegister(credentials);
    if (!err) {
      await checkSession();
    }
    return err;
  }

  async function logout() {
    try {
      await apiLogout();
    } catch (e) {
      console.error("Logout failed", e);
    }
    setPlayerInfo(null);
  }

  return (
    <AuthContext.Provider
      value={{
        playerInfo,
        isLoading,
        login,
        register,
        logout,
        setPlayerInfo,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
