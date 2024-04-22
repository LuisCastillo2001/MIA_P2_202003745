import React from 'react'
import partitionimage from '../Images/partition.png';
import  { useState, useRef } from 'react';
export default function Partition({name, diskname, showLogin}) {
  const [isHovered, setIsHovered] = useState(false);
  const currentpartition = useRef();
  const obtenerParticion = () => {
    fetch(`http://localhost:3000/AccederParticion/${diskname}/${name}`, {
      method: 'GET', 
      headers: {
        'Content-Type': 'application/json'
      }
    })
    .then(response => {
      if (!response.ok) {
        window.alert("Lo sentimos esta partición no tiene un mkfs en el sistema")
        throw new Error(' Error en la petición');
      }
      //Traigo los datos decodificados del json
      return response.json(); 
    })
    .then(data => {
      
      const discosave = JSON.stringify(data);
      localStorage.setItem('disk', discosave);
      showLogin();
      alert("Peticion enviada con exito")
    })
    .catch(error => {
      console.error('Hubo un problema al hacer el fetch', error);
    });
  }
  return (
    <div
      ref={currentpartition}
      onClick={obtenerParticion}
      style={{
        display: 'inline-block',
        marginLeft: '34px',
        marginTop: '25px',
        height: '100px',
        width: '104px',
        position: 'relative',
        backgroundColor: isHovered ? 'rgba(0, 0, 255, 0.2)' : 'transparent',
        transition: 'background-color 0.1s linear',
        
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
    <img src={partitionimage} style={{ height: '90%', width: '85%' }} alt="Partition" />
    <p style={{ textAlign: 'center', marginTop: '5px' }}>{name}</p>
  </div>
  )
}
