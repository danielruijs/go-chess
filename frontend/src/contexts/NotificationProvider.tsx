import type { ReactNode } from "react";
import { Alert, Snackbar } from "@mui/material";
import { useWebSocket } from "./WebSocketContext";

function NotificationProvider({ children }: { children: ReactNode }) {
  const { isDrawDeclinedNoticeOpen, closeDrawDeclinedNotice } = useWebSocket();

  return (
    <>
      {children}
      <Snackbar
        open={isDrawDeclinedNoticeOpen}
        autoHideDuration={5000}
        onClose={closeDrawDeclinedNotice}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      >
        <Alert onClose={closeDrawDeclinedNotice} severity="info" sx={{ width: "100%" }}>
          Your draw offer was declined.
        </Alert>
      </Snackbar>
    </>
  );
}

export default NotificationProvider;
