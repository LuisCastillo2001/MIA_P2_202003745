
import * as React from 'react';
import '../Styles/Proyecto.css'
import { useRef, useState } from 'react';

export default function Commands({ showCommands, commands1}) {
   //Que hice, mande un parametro, y en base a ese parametro decido si muestro o no el componente
   //no fue necesario hacer un useState para mostrar o no el componente
  const textareafile = useRef();
  const inputtext = useRef();
  const [fileContent, setFileContent] = useState("");
  
  const handleenviar = () => {
    if (!textareafile.current.value) {
      if (!inputtext.current.value) {
        window.alert('No hay contenido en el archivo ni comando a ejecutar');
        return;
      }else{
        textareafile.current.value = inputtext.current.value;
      }
     if (!inputtext.current.value) {
        window.alert('No hay contenido en el archivo ni comando a ejecutar');
        return;
      }
    }
    fetch('http://localhost:3000/MandarArchivo', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ fileContent: textareafile.current.value })
      })
      .then(response => {
        if (response.ok) {
          window.alert('Archivo leído correctamente')
        } else {
          console.error('Error al enviar el contenido del archivo al servidor');
        }
      })
      .catch(error => {
        console.error('Error al enviar el contenido del archivo al servidor:', error);
      });
    
    
    };

  
  
  return (
    
    <div style={{width:'100%', height:'100%'}}>
      <textarea id='prueba' ref={textareafile} style={{resize:'none', width:'99%', height:'87%', borderRadius:'10px'}} value={commands1}>
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
              width: '600px',
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
