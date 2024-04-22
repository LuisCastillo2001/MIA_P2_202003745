import React from 'react'
import { useEffect, useState } from 'react';
import { PiKeyReturnBold } from "react-icons/pi";
import { Box, TextField } from '@mui/material';
import '../Styles/Login.css';
export default function Login({returntoPartitions, changeanterior, showDirectorys}) {
  const [partitions, setPartitions] = useState([]);
  const [usuario, setUsuario] = useState("");
  const [password, setPassword] = useState("");
  useEffect(() => {
    changeanterior("Login");
  },[]);

  const prueba = () => {
    const disco = JSON.parse(localStorage.getItem('disk'));
    const iddisk = disco.id;

    const loginweb = {
      user : usuario,
      password : password,
      id : iddisk
    }

    fetch('http://localhost:3000/Login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(loginweb)
  }).then(response => {
    if (!response.ok) {
      window.alert("Usuario o contraseña incorrectos")
      throw new Error(' Error en la petición');
    }else{
      showDirectorys();
      
    }
    return response.json(); 
  }).catch(error => {
    console.error('Hubo un problema al hacer el fetch', error);
  });

  }

  const onchangeUsuario = (e) => {
    setUsuario(e.target.value);
  }

  const onchangePassword = (e) => {
    setPassword(e.target.value);
  }

  return (
    <>
    <div>
      <button onClick={returntoPartitions}>
      <PiKeyReturnBold />
       
      </button>
    </div>
    <div style={{ width: '99%', height: '100%', textAlign: 'left' }}>
       <div style={{width: '80%', height: '80%', textAlign: 'center', marginTop:'40px', marginLeft:'20px'}}>
       <div className="login-box">
            <p>Login</p>
            <form>
                <div className="user-box">
                    <input required="" name="user" type="text" value={usuario} onChange={onchangeUsuario}/>
                    <label>Nombre de usuario</label>
                </div>
                <div className="user-box">
                    <input required="" name="password" type="password" value={password} onChange={onchangePassword}/>
                    <label>Contraseña</label>
                </div>
                <a href="#" onClick={prueba}>
                    <span></span>
                    <span></span>
                    <span></span>
                    <span></span>
                    Verificar
                </a>
            </form>
            
        </div>
      
        
        </div>
    </div>
    </>
  )
}
