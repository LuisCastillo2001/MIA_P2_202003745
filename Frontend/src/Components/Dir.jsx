import React from 'react'
import { useState } from 'react';
import diskimage from '../Images/carpeta.png';
export default function Dirs({name, changeruta}) {
  const [isHovered, setIsHovered] = useState(false);

  const handleClick = () => {
    changeruta(name);
  }

  return (
    <div>
       <div
      onClick={handleClick}
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
      <img src={diskimage} style={{ height: '95%', width: '100%' }} alt="Dirs" />
      <p style={{ textAlign: 'center', marginTop: '10px' , marginLeft:'20px'}}>{name}</p>
    </div>
    </div>
  )
}
