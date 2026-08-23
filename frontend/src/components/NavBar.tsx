import { useState, type MouseEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useWebSocket } from "../contexts/WebSocketContext";
import { useAuth } from "../contexts/AuthContext";
import { Button, Divider, ListItemIcon, ListItemText, Menu, MenuItem } from "@mui/material";
import PersonOutlineIcon from "@mui/icons-material/Person";
import LogoutIcon from "@mui/icons-material/Logout";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import Auth from "./Auth";
import UserAvatar from "./UserAvatar";
import ConnectionStatus from "./ConnectionStatus";

export default function NavBar() {
  const navigate = useNavigate();
  const { isConnected } = useWebSocket();
  const { playerInfo, logout } = useAuth();
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);
  const [userMenuAnchorElement, setUserMenuAnchorElement] = useState<null | HTMLElement>(null);
  const isUserMenuOpen = Boolean(userMenuAnchorElement);

  const handleOpenUserMenu = (event: MouseEvent<HTMLElement>) => {
    setUserMenuAnchorElement(event.currentTarget);
  };

  const handleCloseUserMenu = () => {
    setUserMenuAnchorElement(null);
  };

  const handleProfileClick = () => {
    handleCloseUserMenu();
    if (playerInfo?.username) {
      navigate(`/user/${playerInfo.username}`);
    }
  };

  const handleLogoutClick = () => {
    handleCloseUserMenu();
    logout();
  };

  return (
    <nav className="w-full h-16 bg-white/80 px-4 flex justify-between items-center shadow-sm relative z-50">
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
                <button
                  type="button"
                  onClick={handleOpenUserMenu}
                  className="group flex items-center gap-2 px-2 py-1.5 rounded-xl border border-slate-200 hover:border-slate-300 hover:bg-slate-50 active:bg-slate-100 transition cursor-pointer select-none focus:outline-none"
                >
                  <UserAvatar name={playerInfo.displayName} size="sm" />
                  <span className="text-sm font-medium text-slate-700 select-none flex flex-col items-start leading-tight">
                    <strong className="text-slate-900 font-bold group-hover:text-blue-600 transition-colors">
                      {playerInfo.displayName}
                    </strong>
                    <span className="text-xs text-slate-400 font-normal">
                      @{playerInfo.username}
                    </span>
                  </span>
                  <KeyboardArrowDownIcon
                    fontSize="small"
                    className={`text-slate-400 transition-transform duration-200 ${
                      isUserMenuOpen ? "rotate-180" : ""
                    }`}
                  />
                </button>

                <Menu
                  id="user-dropdown-menu"
                  anchorEl={userMenuAnchorElement}
                  open={isUserMenuOpen}
                  onClose={handleCloseUserMenu}
                >
                  <MenuItem onClick={handleProfileClick}>
                    <ListItemIcon>
                      <PersonOutlineIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText primary="Profile" />
                  </MenuItem>
                  <Divider sx={{ my: 0.5 }} />
                  <MenuItem onClick={handleLogoutClick} sx={{ color: "error.main" }}>
                    <ListItemIcon sx={{ color: "error.main" }}>
                      <LogoutIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText primary="Log Out" />
                  </MenuItem>
                </Menu>
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
