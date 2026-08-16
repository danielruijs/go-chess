import { useState } from "react";
import { Link } from "react-router-dom";
import { useWebSocket } from "../contexts/WebSocketContext";
import { useAuth } from "../contexts/AuthContext";
import { Button } from "@mui/material";
import Auth from "./Auth";

export default function NavBar() {
  const { isConnected } = useWebSocket();
  const { playerInfo, logout } = useAuth();
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);

  return (
    <nav className="w-full bg-white/80 px-6 py-3 flex justify-between items-center shadow-sm relative z-50">
      {/* Left: Logo */}
      <Link
        to="/"
        className="flex gap-0.5 font-bold text-xl text-slate-800 tracking-tight select-none"
      >
        <span className="text-[#00ADD8] font-extrabold">Go</span>Chess
      </Link>

      {/* Right: Info & Actions */}
      <div className="flex gap-6">
        {/* Identity & Button */}
        {isConnected && playerInfo?.displayName && (
          <div className="flex items-center gap-4">
            {playerInfo.isAuthenticated ? (
              <>
                <span className="text-sm font-medium text-slate-700 select-none flex flex-col items-end leading-tight">
                  <strong className="text-slate-900 font-bold">{playerInfo.displayName}</strong>
                  <span className="text-xs text-slate-400 font-normal">@{playerInfo.username}</span>
                </span>
                <Button variant="outlined" color="primary" size="small" onClick={logout}>
                  Log Out
                </Button>
              </>
            ) : (
              <>
                <span className="text-sm font-medium text-slate-700 select-none">
                  <strong className="text-slate-950 font-bold">{playerInfo.displayName}</strong>
                </span>
                <Button variant="contained" size="small" onClick={() => setIsAuthModalOpen(true)}>
                  Log In
                </Button>
              </>
            )}
          </div>
        )}
      </div>

      {/* Login / Register Modal */}
      <Auth open={isAuthModalOpen} onClose={() => setIsAuthModalOpen(false)} />
    </nav>
  );
}
