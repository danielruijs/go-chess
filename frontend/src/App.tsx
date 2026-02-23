import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Home from "./pages/Home.tsx"
import Game from './pages/Game.tsx'
import { WebSocketProvider } from './contexts/WebSocketContext.tsx'

function App() {

  return (
    <WebSocketProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/game" element={<Game />} />
        </Routes>
      </BrowserRouter>
    </WebSocketProvider>
  )
}

export default App
