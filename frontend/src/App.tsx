import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Home from "./pages/Home.tsx"
import Game from './pages/Game.tsx'

function App() {

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/game" element={<Game />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
