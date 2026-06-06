import { useState, useEffect, useCallback, type ReactNode } from "react";
import { AuthContext } from "./AuthContext";
import {
  checkSession as apiCheckSession,
  logout as apiLogout,
  login as apiLogin,
  register as apiRegister,
} from "../api/auth";
import type { Credentials } from "../types/auth";
import type { PlayerInfoData } from "../types/message";
import { useNotification } from "./NotificationContext";

export function AuthProvider({ children }: { children: ReactNode }) {
  const [playerInfo, setPlayerInfo] = useState<PlayerInfoData | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const { showNotification } = useNotification();

  const checkSession = useCallback(async () => {
    try {
      const playerInfoData = await apiCheckSession();
      setPlayerInfo(playerInfoData);
    } catch (e) {
      console.error("Session check failed", e);
      showNotification("Failed to connect to authentication server.", "error");
    } finally {
      setIsLoading(false);
    }
  }, [showNotification]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    checkSession();
  }, [checkSession]);

  async function login(credentials: Credentials) {
    const { data, error } = await apiLogin(credentials);
    if (data) {
      setPlayerInfo(data);
    }
    return error;
  }

  async function register(credentials: Credentials) {
    const { data, error } = await apiRegister(credentials);
    if (data) {
      setPlayerInfo(data);
    }
    return error;
  }

  async function logout() {
    try {
      await apiLogout();
    } catch (e) {
      console.error("Logout failed", e);
      showNotification("Logout failed. Please try again.", "error");
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
