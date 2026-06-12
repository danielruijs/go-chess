import { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  Tabs,
  Tab,
  TextField,
  Button,
  Box,
  Alert,
  IconButton,
} from "@mui/material";
import { useAuth } from "../contexts/AuthContext";
import CloseIcon from "@mui/icons-material/Close";
import type { Credentials } from "../types/auth";

type AuthProps = {
  open: boolean;
  onClose: () => void;
};

function Auth({ open, onClose }: AuthProps) {
  const { login, register } = useAuth();
  const [tabIndex, setTabIndex] = useState(0);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabIndex(newValue);
    setError(null);
    setPassword("");
    setDisplayName("");
  };

  const handleSubmit = async (e: React.SubmitEvent) => {
    // Prevent the default browser form submission behavior
    e.preventDefault();

    setError(null);
    setLoading(true);

    const credentials: Credentials = { username, password, displayName };
    const err = tabIndex === 0 ? await login(credentials) : await register(credentials);

    setLoading(false);
    if (err) {
      setError(err.message);
    } else {
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 1 }}>
        <DialogTitle sx={{ p: 0, fontWeight: 700, fontSize: "1.5rem", color: "#1e293b" }}>
          {tabIndex === 0 ? "Welcome Back" : "Create Account"}
        </DialogTitle>
        <IconButton onClick={onClose} size="small" sx={{ color: "#64748b" }}>
          <CloseIcon />
        </IconButton>
      </Box>

      <Tabs value={tabIndex} onChange={handleTabChange} variant="fullWidth" sx={{ mb: 2 }}>
        <Tab label="Log In" />
        <Tab label="Sign Up" />
      </Tabs>

      <DialogContent sx={{ p: 0, overflow: "visible" }}>
        <form onSubmit={handleSubmit}>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            {error && (
              <Alert severity="error" sx={{ borderRadius: 4 }}>
                {error.charAt(0).toUpperCase() + error.slice(1)}
              </Alert>
            )}

            <TextField
              label="Username"
              variant="outlined"
              fullWidth
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              disabled={loading}
              required
            />

            <TextField
              label="Password"
              type="password"
              variant="outlined"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={loading}
              required
            />

            {tabIndex === 1 && (
              <TextField
                label="Display Name"
                variant="outlined"
                fullWidth
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                disabled={loading}
                required
              />
            )}

            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={loading}
              sx={{
                mt: 1,
                py: 1.5,
                fontSize: "1rem",
              }}
            >
              {loading ? "Please wait..." : tabIndex === 0 ? "Log In" : "Sign Up"}
            </Button>
          </Box>
        </form>
      </DialogContent>
    </Dialog>
  );
}

export default Auth;
