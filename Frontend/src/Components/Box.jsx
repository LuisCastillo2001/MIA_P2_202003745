import * as React from 'react';
import Box from '@mui/material/Box';
import '../Styles/Proyecto.css'
import Commands from './Commands';
import { useState, useRef } from 'react';
import image from '../Components/anime2.jpg'
import Disks from './Disks';
export default function ContentDisplayBox() {
  const [fileContent, setFileContent] = useState("");
  const [commands, setCommands] = useState(true); 
  const [disks, setDisks] = useState(false);
 
  const divStyle = {
    width:'200px',
    padding: '20px',
    borderRadius: '5px',
    backgroundColor:'#cccccc',
    boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)',
    color:'black'
  };
 
  const handleFileChange = (event) => {
    const file = event.target.files[0];
    
    const reader = new FileReader();

    reader.onload = (e) => {
      const content = e.target.result;
      setFileContent(content);
      
    }
      
    reader.readAsText(file);
    alert('Archivo subido correctamente');
  };

  const showDisks= () => {
    setCommands(false);
    setDisks(true);
  }

  const showCommands = () => {
    setCommands(true);
    setDisks(false);
  }
  
  return (

    <>
    
    <div style={{ display: 'flex', alignItems:'center', marginTop:'10px', marginLeft:'112px'}}>
      <label htmlFor="file-upload" className="custom-file-upload">
      Subir Archivo
      </label>
      <input
        id="file-upload"
        type="file"
        style={{ display: 'none' }}
        onChange={handleFileChange} 
      />
      </div>

    <div style={{ display: 'flex', alignItems:'center', marginTop:'10px', marginLeft:'112px'}}>
      
    <Box
        height={650}
        width={1000}
        my={4}
        display="flex"   
        sx={{
        border: '2px solid grey',
        borderRadius: '10px', 
        boxShadow: '0 0 10px rgba(0, 0, 0, 0.5)'
    }}
        marginLeft = {5}
        backgroundColor = '#333333'
   
    >
      
    <div style={divStyle}>
        
      <button style={{width:'100%', height:'36px', fontSize:'14px', marginTop:'10px'}} onClick={showCommands}>Enviar Comandos</button>
      <button style={{width:'100%', height:'36px', fontSize:'14px', marginTop:'10px'}} onClick={showDisks}>Discos</button>
    </div>
    
  
    {disks === true && <Disks showdisks={disks} ></Disks>}
    {commands === true && <Commands showCommands={commands} commands1={fileContent} ></Commands>}
    </Box>

    <img src={image}  alt='imagen' style={{width:'20%', height:'30%', marginLeft:'15px'}}></img>
    
    
    </div>
    
    
    </>
    
  );
}

