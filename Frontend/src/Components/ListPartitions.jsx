import React from 'react'
import Partition from './Partition';
import { useEffect, useState } from 'react';
import { PiKeyReturnBold } from "react-icons/pi";
export default function ListPartitions({namedisk, returntodisk, changeanterior, showLogins}) {
  const [partitions, setPartitions] = useState([]);
  useEffect(() => {
    changeanterior("Part");
  },[]);

  
  useEffect(() => {
    fetch(`http://localhost:3000/ListaParticiones/${namedisk}`
    , {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }
  )
      .then(response => {
        if (response.ok) {
          return response.json();
        }
        throw new Error('Error al obtener las particiones del disco');
      })
      .then(data => {
        if (data == null){
          setPartitions([]);
          return;
        }
        console.log('Particiones del disco:', data);
        setPartitions(data); 
      })
      .catch(error => {
        console.error('Error:', error);
      });
  }, []);
  return (
    <>
    <div>
      <button onClick={returntodisk}>
      <PiKeyReturnBold />
       
      </button>
    </div>
    <div style={{ width: '99%', height: '100%', textAlign: 'left' }}>
        {partitions.length > 0 && partitions.map((partition) => (
        <Partition  name={partition}  diskname={namedisk} showLogin={showLogins} />
      ))}
    </div>
    </>
  )
}
