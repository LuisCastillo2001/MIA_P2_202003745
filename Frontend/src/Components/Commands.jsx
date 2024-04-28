
import * as React from 'react';
import '../Styles/Proyecto.css'
import { useRef, useState, useEffect } from 'react';

export default function Commands({ showCommands, commands1}) {
   //Que hice, mande un parametro, y en base a ese parametro decido si muestro o no el componente
   //no fue necesario hacer un useState para mostrar o no el componente
  const textareafile = useRef();
  const inputtext = useRef();
  
  useEffect(() => {
    textareafile.current.value = commands1;
  }, [commands1]);
  
  const handleenviar = () => {
    let lines = ""
    if (textareafile.current.value === '') {
      if (inputtext.current.value === '') {
        return;
      }
      lines = inputtext.current.value.split('\n');
      inputtext.current.value = '';
    }else{
      lines = textareafile.current.value.split('\n');
      textareafile.current.value = '';
    }

    
    
  
   
    const fetchPromises = [];
  
    lines.forEach((line, index) => {
      const promise = new Promise((resolve, reject) => {
        setTimeout(() => {
          fetch('http://54.163.43.245:3000/MandarArchivo', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({ fileContent: line })
          })
            .then(response => {
              if (response.ok) {
                return response.json();
              } else {
                throw new Error('Error al enviar el contenido del archivo al servidor');
              }
            })
            .then(data => {
            
              data.forEach(item => {
                textareafile.current.value += item + '\n';
                commands1 += item + '\n';
              });
              resolve(); // Marcar la promesa como resuelta
            })
            .catch(error => {
              console.error('Error al enviar el contenido del archivo al servidor:', error);
              reject(error); 
            });
        }, index * 1000); 
      });
  
      fetchPromises.push(promise); 
    });
  
   
    Promise.all(fetchPromises)
      .then(() => {
        
        descargarArchivo('comandos.txt', textareafile.current.value);
      })
      .catch(error => {
        console.error('Error al procesar las solicitudes fetch:', error);
      });
  };

  function descargarArchivo(nombreArchivo, contenidoArchivo) {
    
    const blob = new Blob([contenidoArchivo], { type: 'text/plain' });
  
   
    const url = URL.createObjectURL(blob);
  
   
    const enlaceDescarga = document.createElement('a');
    enlaceDescarga.href = url;
    enlaceDescarga.download = nombreArchivo;
  
   
    document.body.appendChild(enlaceDescarga);
    enlaceDescarga.click();
  
    
    URL.revokeObjectURL(url);
  
   
    document.body.removeChild(enlaceDescarga);
  }
  
  return (
    
    <div style={{width:'100%', height:'100%'}}>
      <textarea id='prueba' ref={textareafile} style={{resize:'none', width:'99%', height:'87%', borderRadius:'10px', fontSize:'15px'}}>
      </textarea>
      <div style={{display:'flex'}}>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems:'flex-start' }}>
          <label htmlFor="texto-input" style={{ color: 'white', marginBottom: '5px', marginLeft:'15px' }}>
            Comando a ejecutar
          </label>
          <input
            ref={inputtext}
            id="texto-input"
            type="text"
            style={{
              backgroundColor: '#333333',
              width: '700px',
              height: '26px',
              borderRadius: '10px',
              border: '2px solid white',
              marginLeft:'5px'
            }}
          />
        </div>
        <button
          style={{
            width:'100px',
            height:'36px',
            marginLeft:'10px',
            marginTop:'24px',
            fontSize:'14px',
          }} 
          onClick={handleenviar}
        >
          Enviar
        </button>
      </div>
    </div>
  ) 
}
