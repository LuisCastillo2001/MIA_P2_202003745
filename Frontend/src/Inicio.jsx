import { useRef, useState } from 'react'
import image from './inicio.jpg'
import './App.css'
import { useNavigate } from 'react-router-dom'

//Esto lo usuare para aprender react
//La pantalla principal la pondré en otro lado
function Inicio(){
  let titulo = "Manejo e implementación de archivos"
  let titulo2 = <h2>Aqui pondre las cosas para aprender react</h2>
  
    const navigate = useNavigate()
  //aqui uso el hook usestate el cual primero va la variable y luego la funcion que la modifica, le puedo
  //dar un valor inicial
  const [contador, setContador] = useState(1)
  const [numero1, setNumero1] = useState(0)
  const [numero2, setNumero2] = useState(0)
  const [resultado, setResultado] = useState(0)

  //Para hacer referencias
  const handleonclick = () => {
    navigate('/Proyecto')
  }


 
  const incrementar = (e) =>{
      //Formas de acceder a elementos del DOM
      //e.target.style.backgroundColor = 'red'
      //document.getElementById('contador').innerHTML = Number(document.getElementById('contador').innerHTML) + 1

      //el ref siempre lleva current para saber que es una referencia
      //console.log(contadorRef.current.innerHTML)
      //valor que voy a modificar
      setContador(contador + 1)
      if (contador === 10){
        setContador(0)
      }
  }

 
  const suma = () =>{

  }
  return (
    <>
    <h1> {titulo} </h1> 
    {titulo2} 
    <div className='prueba'>
      <div className='contendor1'>
        <div className='contednor1'>
        <button className='boton_img' onClick={incrementar}>Contador de clicks</button>

        <p>{contador}</p>
        {/*
        <p>suma</p>
        <p>numero1</p>
        <input type='number' value={numero1} onChange={modificarnumero1}/>
        <p>numero2</p>
        <input type='number'value={numero2} onChange={modificarnumero2}/>
        <p>resultado</p>
        <input type='number' readOnly value={resultado}/>
        <button onClick={mostrarresultado}>Sumar</button>
        */}
        </div>
       

        
       

        <img src = {image} alt = "imagen" />
       
      </div>
    
    <button className='boton_prueba' onClick={handleonclick}>Boton de prueba</button>
     
   
    </div>
    </>
  )
}

export default Inicio
