import React from 'react'
import { useEffect } from 'react';
export default function Directorys({Logout, changeanterior}) {
    useEffect(() => {
        changeanterior("Dir");
      },[]);


  return (
    <>
    <button onClick={Logout}/>
    <div>
      <h1>Se activaron los directorios</h1>
    </div>
    </>
  )
}
