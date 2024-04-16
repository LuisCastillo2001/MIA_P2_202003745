import React from 'react'
import { useRef, useState } from 'react'
import ContentDisplayBox from '../Components/Box'
import SpringModal from '../Components/springmodal'
import '../Styles/Proyecto.css'
export default function Proyecto() {
  const prueba = useRef()
 
  /*
  function inicio(){
    
    personas.map(persona =>{
      //Este es el map, va recorriendo los valores del arreglo que le mannde
      prueba.current.innerHTML += `<h1> ${persona.nombre} ${persona.edad} </h1>`
    })
  }
  */
 
  return (
    <>
    
     <SpringModal></SpringModal>
      
      
        <ContentDisplayBox>
      
        </ContentDisplayBox>
  
     
    </>
  )
}

