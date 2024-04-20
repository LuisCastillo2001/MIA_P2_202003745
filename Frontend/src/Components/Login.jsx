import React from 'react'
import { useEffect, useState } from 'react';
import { PiKeyReturnBold } from "react-icons/pi";
import { Box, TextField } from '@mui/material';
import '../Styles/Login.css';
export default function Login({returntoPartitions, changeanterior}) {
  const [partitions, setPartitions] = useState([]);
  useEffect(() => {
    changeanterior("Login");
  },[]);

  const prueba = () => {
    console.log("hola");
    alert('Usuario o contraseña incorrectos');
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
                    <input required="" name="user" type="text" />
                    <label>Nombre de usuario</label>
                </div>
                <div className="user-box">
                    <input required="" name="password" type="password" />
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
