import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "@mui/material/styles";
import { theme } from "./theme.ts";
import Home from "./pages/Home.tsx";
import Game from "./pages/Game.tsx";
import Analysis from "./pages/Analysis.tsx";
import { AuthProvider } from "./contexts/AuthContext.tsx";
import WebSocketProvider from "./contexts/WebSocketContext.tsx";
import NotificationProvider from "./contexts/NotificationProvider.tsx";
import NavBar from "./components/NavBar.tsx";
import ConnectionStatus from "./components/ConnectionStatus.tsx";

function App() {
  return (
    <ThemeProvider theme={theme}>
      <BrowserRouter>
        <NotificationProvider>
          <AuthProvider>
            <WebSocketProvider>
              <div className="flex flex-col min-h-screen">
                <NavBar />
                <div className="flex-1 relative flex flex-col">
                  <div className="absolute top-4 right-4 z-50">
                    <ConnectionStatus />
                  </div>
                  <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/match/:publicId" element={<Game />} />
                    <Route path="/analysis/:publicId" element={<Analysis />} />
                  </Routes>
                </div>
              </div>
            </WebSocketProvider>
          </AuthProvider>
        </NotificationProvider>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
