import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Home from "./pages/Home.tsx"
import Game from './pages/Game.tsx'
import { WebSocketProvider } from './contexts/WebSocketContext.tsx'

function App() {

  return (
    <BrowserRouter>
      <WebSocketProvider>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/game" element={<Game />} />
        </Routes>
      </WebSocketProvider>
    </BrowserRouter>
  )
}

export default App
