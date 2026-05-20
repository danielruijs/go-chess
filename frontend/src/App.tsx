import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/Home.tsx";
import Game from "./pages/Game.tsx";
import WebSocketProvider from "./contexts/WebSocketContext.tsx";
import NotificationProvider from "./contexts/NotificationProvider.tsx";
import ConnectionStatus from "./components/ConnectionStatus.tsx";

function App() {
  return (
    <BrowserRouter>
      <WebSocketProvider>
        <NotificationProvider>
          <div className="relative min-h-screen">
            <div className="absolute top-4 right-4 z-50">
              <ConnectionStatus />
            </div>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/game" element={<Game />} />
            </Routes>
          </div>
        </NotificationProvider>
      </WebSocketProvider>
    </BrowserRouter>
  );
}

export default App;
