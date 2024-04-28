import React from 'react';
import { useEffect, useState } from 'react';
import { CiLogout } from "react-icons/ci";
import {InputGroup, FormControl, Button} from 'react-bootstrap';
import { FaSearch } from "react-icons/fa";
import Dirs from './Dir';
import '../Styles/Dirs.css';

export default function Directorys({Logout, changeanterior}) {

    const [ruta, setRuta] = useState("/");
    const [dirs, setDirs] = useState([]);
    
    const disco = JSON.parse(localStorage.getItem('disk'));
    const diskname = disco.diskname + ".dsk"
    const partition = disco.partitionname
    
    


    

    useEffect(() => {
        const getdir = {
            diskname : diskname,
            partition : partition,
            path : ruta
        
        }
        changeanterior("Dir");
        fetch('http://54.163.43.245:3000/ObtenerDirectorios', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }, body: JSON.stringify(getdir)       
    }).then(response => {
        if (!response.ok) {
            throw new Error(' Error en la petición');
        }
        return response.json(); 
        
    }).then(data => {
       
        setDirs(data);


    }).catch(error => {
        console.error('Error:', error);
    });
    },[]);


    const changeruta = (dir) => {
       
        const newRuta = ruta + dir + "/";
        setRuta(newRuta);
        search(newRuta);
    }

    const handlerutaChange = (e) => {
        setRuta(e.target.value);
    }

    const search = (searchRuta = ruta) => {

        if (typeof searchRuta !== 'string') {
            
            searchRuta = ruta

            

            if (searchRuta[searchRuta.length-1] != "/") {
                searchRuta = searchRuta + "/";
                setRuta(searchRuta);
            }
        } 

        const getdir = {
            diskname : diskname,
            partition : partition,
            path : searchRuta
        
        }
        alert("Buscando en: " + searchRuta);
        fetch('http://54.163.43.245:3000/ObtenerDirectorios', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }, body: JSON.stringify(getdir)       
        }).then(response => {
            if (!response.ok) {
                window.alert("La carpeta no existe")
                throw new Error(' Error en la petición');
            }
            return response.json(); 
            
        }).then(data => {
            if (data == null){
                setDirs([]);
                return;
              }
           
            setDirs(data);
    
    
        }).catch(error => {
            console.error('Error:', error);
        });

    }

    return (
        <>
            
          <div className="search-container" style={{marginTop:'10px', marginLeft:'10px'}}>
              <InputGroup className="mb-3">
                  <InputGroup.Text>
                      <Button variant="primary" className="search-btn" style={{height:'40px', marginTop:'5px'}} 
                      onClick={search}>
                      <FaSearch />
                      </Button>
                  </InputGroup.Text>
                  <FormControl style={{height:'30px', fontSize:'24px', width:'650px'}}
                      placeholder="Buscar Archivo"
                      aria-label="Buscar Archivo"
                      aria-describedby="basic-addon1"
                      value={ruta}
                      onChange={handlerutaChange}
                      id="ruta"
                  />
              </InputGroup>
              <div style={{display:'flex', flexWrap:'wrap'}} id = "directorios">
              {dirs.length > 0 && dirs.map((dir, index) => (
                 <Dirs  name={dir} changeruta={changeruta}  />
                ))}
                </div>
          </div>
         
            
        </>
    );
}
