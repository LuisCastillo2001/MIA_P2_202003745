import React from 'react';
import { useState, useEffect } from 'react';
import Springrep from './springRep';
export default function Reportes() {
  const [reportes, setReportes] = useState([]);

  useEffect(() => {
    fetch('http://localhost:3000/Reportes',{
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Error al obtener los datos');
        }
        return response.json();
      })
      .then(data => {
        if (data == null){
          return;
        }
        setReportes(data);
        console.log(data);
      })
      .catch(error => {
        console.error('Error:', error);
      });
    }, []);

  return (
    <div style={{ width: '90%', height: '100%', textAlign: 'left' }}>
      <div style={{display:'flex', flexWrap:'wrap'}}>
      {reportes.length > 0 && reportes.map((reporte, index) => (
        <Springrep
          key={index} 
          name={reporte.name} 
          dot={reporte.dot} 
        />
      ))}
    </div>
    </div>
  );
}
