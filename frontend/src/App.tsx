import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Home from "./pages/Home.tsx"
import Game from './pages/Game.tsx'
import { WebSocketProvider } from './contexts/WebSocketContext.tsx'
import NotificationProvider from './contexts/NotificationProvider.tsx'

function App() {

  return (
    <BrowserRouter>
      <WebSocketProvider>
        <NotificationProvider>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/game" element={<Game />} />
          </Routes>
        </NotificationProvider>
      </WebSocketProvider>
    </BrowserRouter>
  )
}

export default App
