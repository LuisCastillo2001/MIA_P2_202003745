
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Inicio from './Inicio';
import Proyecto from "./Pages/Proyecto";
//Esto lo usuare para aprender react
//La pantalla principal la pondré en otro lado
function App(){
  return (
  
    
    <Router>
      <Routes>
        <Route path="/" element={<Inicio />} />
        <Route path="/Proyecto" element={<Proyecto />} />
      </Routes>
    </Router>
  
    
  )
}

export default App
