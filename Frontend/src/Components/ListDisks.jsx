import React, { useEffect, useState } from 'react';
import Disk from './Disk';

export default function ListDisks({ showdisks, prueba1, changeanterior }) {
  const [disks, setDisks] = useState([]);
  useEffect(() => {
    changeanterior("Disk");
  },[]); 
  useEffect(() => {
    fetch('http://localhost:3000/ListaDiscos')
      .then(response => {
        if (!response.ok) {
          throw new Error('Error al obtener los datos');
        }
        return response.json();
      })
      .then(data => {
        if (data == null){
          setDisks([]);
          return;
        }
        setDisks(data);
      })
      .catch(error => {
        console.error('Error:', error);
      });
  }, []); 

  return (
    <div style={{ width: '99%', height: '100%', textAlign: 'left' }}>
      {disks.length > 0 && disks.map((disk, index) => (
        <Disk  name={disk} onPruebaClick={prueba1} />
      ))}
      
    </div>
  );
}
