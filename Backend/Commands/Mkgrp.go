package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"

	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Mkgrp struct {
	Name       string
	Parameters []string
}

func NewMkgrp(parameters []string) *Mkgrp {
	mkgrp := &Mkgrp{
		Name:       "",
		Parameters: parameters,
	}
	mkgrp.readParameters()
	return mkgrp
}

func (mkgrp *Mkgrp) readParameters() {
	for _, parametro := range mkgrp.Parameters {
		mkgrp.IdentifyParameters(parametro)
	}
	mkgrp.makeGroup()

}

func (mkgrp *Mkgrp) IdentifyParameters(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "name" {
		mkgrp.Name = Stringmake(parameterIdentifier[1])
	}

}

func (mkgrp *Mkgrp) makeGroup() {
	if Logged.User != "root" {
		Concatenar("Lo sentimos, solo el usuario root puede realizar esta acción")
		return
	}

	var path string
	partition := getMount(Logged.Id, &path)
	if partition.Part_start == -1 {
		Concatenar("La particion no ha sido montada")
		return
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		return
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Concatenar("Error al leer el archivo")
		return
	}

	var fb Structs.BloquesArchivos
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(bloque-1), 0)

		data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			Concatenar("Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'G' || linea[2] == 'g' {
			c++
			in := strings.Split(linea, ",")
			if in[2] == mkgrp.Name {
				if linea[0] != '0' {
					Concatenar("El nombre del grupo ya existe")
					return
				}
			}
		}
	}
	txt += strconv.Itoa(c+1) + ",G," + mkgrp.Name + "\n"

	tam := len(txt)
	var cadenasS []string

	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		Concatenar("Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Concatenar("No se ha encontrado el disco")
		return
	}
	for i := 0; i < len(cadenasS); i++ {

		var fbAux Structs.BloquesArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WriteBytes(file, binAux.Bytes())

			diferencia := super.S_blocks_count - super.S_free_blocks_count
			file.Seek(super.S_bm_block_start+diferencia, 0)
			WriteBytes(file, []byte{byte('1')})
			super.S_free_blocks_count--

			inode.I_block[i] = diferencia
		} else {
			fbAux = fb
		}
		copy(fbAux.B_content[:], cadenasS[i])
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		WriteBytes(file, bin6.Bytes())

	}

	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WriteBytes(file, inodos.Bytes())

	file.Seek(partition.Part_start, 0)
	var super2 bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, super2)
	WriteBytes(file, super2.Bytes())
	Concatenar("Grupo" + mkgrp.Name + ", creado correctamente!")

	if super.S_filesystem_type == 3 {
		var journalingbytes bytes.Buffer
		var journaling Structs.Journaling
		copy(journaling.Operacion[:], "mkgrp")
		copy(journaling.Ruta[:], "/users.txt")
		fecha := time.Now().String()
		copy(journaling.Fecha[:], fecha)
		copy(journaling.Contenido[:], "group: "+mkgrp.Name+"grp")
		journaling.Active = '1'
		file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
		binary.Write(&journalingbytes, binary.BigEndian, journaling)
		WriteBytes(file, journalingbytes.Bytes())
	}

	file.Close()
}
