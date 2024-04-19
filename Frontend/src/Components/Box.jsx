import React, { useState } from 'react';
import Box from '@mui/material/Box';
import '../Styles/Proyecto.css';
import Commands from './Commands';
import image from '../Images/anime2.jpg';
import ListDisks from './ListDisks';
import ListPartitions from './ListPartitions';
export default function ContentDisplayBox() {
  const [showCommands, setCommands] = useState(true);
  const [showDisks, setDisks] = useState(false);
  const [showPartitions, setPartitions] = useState(false);
  const [Partitionname, setPartitionname] = useState("");
  const [fileContent, setFileContent] = useState("");
  const [anterior, setAnterior] = useState("Disk");
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

  const handlePruebaClick = (namedisk) => {
    setDisks(false); 
    setPartitions(true);
    setPartitionname(namedisk);
  };

  const handleCommandsClick = () => {
    setCommands(true);
    setDisks(false);
    setPartitions(false);
    
  };

  const handleDisksClick = () => {
    if (anterior === "Part"){
      setCommands(false);
      setDisks(false);
      setPartitions(true);
      return
    }
    else if (anterior === "Disk"){
      setCommands(false);
      setDisks(true);
      setPartitions(false);
      return
    }
  };

  const changeanterior = (anterior1) => {
    setAnterior(anterior1)
  }


  const returntodisk = () =>
      {
        setCommands(false);
        setDisks(true);
        setPartitions(false);
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
          marginLeft={5}
          backgroundColor="#333333"
        >
          <div style={{ width:'200px', padding: '20px', borderRadius: '5px', backgroundColor:'#cccccc', boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)', color:'black' }}>
            <button style={{width:'100%', height:'36px', fontSize:'14px', marginTop:'10px'}} onClick={handleCommandsClick}>Enviar Comandos</button>
            <button style={{width:'100%', height:'36px', fontSize:'14px', marginTop:'10px'}} onClick={handleDisksClick}>Discos</button>
          </div>

          {showDisks && <ListDisks showdisks={showDisks} prueba1={handlePruebaClick} changeanterior={changeanterior} />}
          {showCommands && <Commands showCommands={showCommands} commands1={fileContent}/>}
          {showPartitions && <ListPartitions namedisk={Partitionname} returntodisk={returntodisk} changeanterior={changeanterior} />}
        </Box>

        <img src={image} alt='imagen' style={{width:'20%', height:'30%', marginLeft:'15px'}} />
      </div>
    </>
  );
}
