import React, { useState } from 'react';
import Box from '@mui/material/Box';
import '../Styles/Proyecto.css';
import Commands from './Commands';
import image from '../Images/anime2.jpg';
import ListDisks from './ListDisks';
import ListPartitions from './ListPartitions';
import Directorys from './Directorys';
import Login from './Login';
import Reportes from './Reportes';
export default function ContentDisplayBox() {
  const [showCommands, setCommands] = useState(true);
  const [showDisks, setDisks] = useState(false);
  const [showPartitions, setPartitions] = useState(false);
  const [showLogin, setLogin] = useState(false);
  const [Partitionname, setPartitionname] = useState("");
  const [fileContent, setFileContent] = useState("");
  const [showDirectorys, setDirectorys] = useState("");
  const [anterior, setAnterior] = useState("Disk");
  const [showReportes, setReportes] = useState(false);
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
    setLogin(false);
    setDirectorys(false);
  };

  const showLogins = () => {
    setCommands(false);
    setDisks(false);
    setPartitions(false);
    setLogin(true);
    setDirectorys(false);
  }

  const handleCommandsClick = () => {
    setCommands(true);
    setDisks(false);
    setPartitions(false);
    setLogin(false);
    setDirectorys(false);
    
  };

  const handleDirectoryclick = () => {
    setCommands(false);
    setDisks(false);
    setPartitions(false);
    setLogin(false);
    setDirectorys(true);
  }

  const handleDisksClick = () => {
    if (anterior === "Part"){
      setCommands(false);
      setDisks(false);
      setPartitions(true);
      setLogin(false)
      setDirectorys(false)
      return
    }
    else if (anterior === "Disk"){
      setCommands(false);
      setDisks(true);
      setPartitions(false);
      setLogin(false)
      setDirectorys(false)
      return
    }else if (anterior === "Login"){
      setCommands(false);
      setDisks(false);
      setPartitions(false);
      setLogin(true)
      setDirectorys(false)
      return
    }else if (anterior === "Dir"){
      setCommands(false);
      setDisks(false);
      setPartitions(false);
      setLogin(false)
      setDirectorys(true)
      return
    }
  };

  const changeanterior = (anterior1) => {
    setAnterior(anterior1)
  }

  const Logout = () => {

   

    setCommands(false);
    setDisks(false);
    setPartitions(false);
    setLogin(true);
    setDirectorys(false);
    fetch('http://localhost:3000/Logout', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
     
    }).then(response => {
      if (!response.ok) {
        alert("Puede que ya haya un usuario logueado")
        throw new Error(' Error en la petición');
       
      }else{
        window.alert("Se ha cerrado la sesión")
      }
      return response.json(); 
    }).catch(error => {
      console.error('Hubo un problema al hacer el fetch', error);
    });
  }

  const returntodisk = () =>
      {
        setCommands(false);
        setDisks(true);
        setPartitions(false);
        setLogin(false);
        setDirectorys(false);
      }

  const returntoPartitions = () =>
      {
        setCommands(false);
        setDisks(false);
        setPartitions(true);
        setLogin(false);
        setDirectorys(false);
      }

  const handleReportes = () => {
    setCommands(false);
    setDisks(false);
    setPartitions(false);
    setLogin(false);
    setDirectorys(false);
    setReportes(true);
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
            <button style={{width:'100%', height:'36px', fontSize:'14px', marginTop:'10px'}} onClick={handleReportes}>Reportes</button>
          </div>

          {showDisks && <ListDisks showdisks={showDisks} prueba1={handlePruebaClick} changeanterior={changeanterior} />}
          {showCommands && <Commands showCommands={showCommands} commands1={fileContent}/>}
          {showPartitions && <ListPartitions namedisk={Partitionname} returntodisk={returntodisk} changeanterior={changeanterior} showLogins={showLogins} />}
          {showLogin && <Login returntoPartitions={returntoPartitions} changeanterior={changeanterior} showDirectorys={handleDirectoryclick}  />}
          {showDirectorys && <Directorys  changeanterior={changeanterior} Logout={Logout} /> } 
        </Box>

        <div style={{ position: 'relative', width: '20%', height: '30%' }}>
      <img src={image} alt='imagen' style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
      <div style={{ position: 'absolute', bottom: 0, left: '20px', width: '100%', textAlign: 'center' }}>
        <button style={{ border: '2px solid grey' }} onClick={Logout}>Logout</button>
      </div>
    </div>
      </div>
      
    </>
  );
}
