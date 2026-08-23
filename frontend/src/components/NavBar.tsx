import { useState } from "react";
import { Link } from "react-router-dom";
import { useWebSocket } from "../contexts/WebSocketContext";
import { useAuth } from "../contexts/AuthContext";
import { Button } from "@mui/material";
import Auth from "./Auth";
import UserAvatar from "./UserAvatar";
import ConnectionStatus from "./ConnectionStatus";

export default function NavBar() {
  const { isConnected } = useWebSocket();
  const { playerInfo, logout } = useAuth();
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);

  return (
    <nav className="w-full h-16 bg-white/80 px-4 sm:px-6 flex justify-between items-center shadow-sm relative z-50">
      {/* Left: Logo */}
      <Link
        to="/"
        className="flex gap-0.5 font-bold text-xl text-slate-800 tracking-tight select-none"
      >
        <span className="text-go font-extrabold">Go</span>Chess
      </Link>

      {/* Right: Info & Actions */}
      <div className="flex items-center gap-2">
        <ConnectionStatus />
        {/* Identity & Button */}
        {isConnected && playerInfo?.displayName && (
          <div className="flex items-center gap-2">
            {playerInfo.isAuthenticated ? (
              <>
                <Link
                  to={`/user/${playerInfo.username}`}
                  className="group flex items-center gap-2 px-1.5 py-1.5 rounded-xl border border-slate-200 hover:border-slate-300 hover:bg-slate-50 active:bg-slate-100 transition cursor-pointer select-none"
                  title="View your profile"
                >
                  <UserAvatar name={playerInfo.displayName} size="sm" />
                  <span className="text-sm font-medium text-slate-700 select-none flex flex-col leading-tight">
                    <strong className="text-slate-900 font-bold group-hover:text-blue-600 transition-colors">
                      {playerInfo.displayName}
                    </strong>
                    <span className="text-xs text-slate-400 font-normal">
                      @{playerInfo.username}
                    </span>
                  </span>
                </Link>
                <Button variant="outlined" color="primary" size="small" onClick={logout}>
                  Log Out
                </Button>
              </>
            ) : (
              <>
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
