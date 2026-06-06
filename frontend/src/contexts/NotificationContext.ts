import { createContext, useContext } from "react";

export type NotificationSeverity = "success" | "info" | "warning" | "error";

export interface NotificationContextType {
  showNotification: (message: string, severity: NotificationSeverity) => void;
}

export const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export function useNotification() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error("useNotification must be used within a NotificationProvider");
  }
  return context;
}
