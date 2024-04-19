import React, { useState, useRef } from 'react';
import diskimage from '../Images/disco.png';

export default function Disk({ name, onPruebaClick }) {
  const [isHovered, setIsHovered] = useState(false);
  const currentdisk = useRef();

  const prueba = (e) => {
    e.target.parentNode.parentNode.innerHTML = ""
    onPruebaClick(name); 
    
    
  };

  return (
    <div
      ref={currentdisk}
      onClick={prueba}
      style={{
        display: 'inline-block',
        marginLeft: '34px',
        marginTop: '25px',
        height: '100px',
        width: '104px',
        position: 'relative',
        backgroundColor: isHovered ? 'rgba(0, 0, 255, 0.3)' : 'transparent',
        transition: 'background-color 0.31s ease',
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <img src={diskimage} style={{ height: '90%', width: '100%' }} alt="Disk" />
      <p style={{ textAlign: 'center', marginTop: '5px' }}>{name}</p>
    </div>
  );
}
