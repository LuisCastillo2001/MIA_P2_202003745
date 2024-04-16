package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Rmgrp struct {
	Name       string
	Parameters []string
}

func newRmgrp(parameters []string) *Rmgrp {
	rmgrp := &Rmgrp{
		Name:       "",
		Parameters: parameters,
	}
	rmgrp.readParameters()
	return rmgrp
}

func (rmgrp *Rmgrp) readParameters() {
	for _, parametro := range rmgrp.Parameters {
		rmgrp.identifyParameters(parametro)

	}
	rmgrp.RemoveGroup()
}

func (rmgrp *Rmgrp) identifyParameters(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "name" {
		rmgrp.Name = Stringmake(parameterIdentifier[1])
	}
}

func (rmgrp *Rmgrp) RemoveGroup() {

	if !Compare(Logged.User, "root") {
		fmt.Println("RMGRP", "Solo el usuario \"root\" puede acceder a estos comandos.")
		return
	}

	var path string
	partition := getMount(Logged.Id, &path)
	if partition.Part_start == -1 {
		fmt.Println("No se encontro la partición montada, hubo un fmt.println")
		return
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("Hubo un fmt.println al abrir el disco")
		return
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		fmt.Println("Error al leer el archivo")
		return
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		fmt.Println("Error al leer el archivo")
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
			fmt.Println("Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""

	vctr := strings.Split(txt, "\n")
	existe := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'G' || linea[2] == 'g') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[2] == rmgrp.Name {
				existe = true
				aux += strconv.Itoa(0) + ",G," + in[2] + "\n"
				continue
			}
		}
		aux += linea + "\n"
	}
	if !existe {
		fmt.Println("No se encontró el grupo \"" + rmgrp.Name + "\".")
		return
	}
	txt = aux

	tam := len(txt)
	var contenido []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			contenido = append(contenido, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			contenido = append(contenido, txt)
		}
	} else {
		contenido = append(contenido, txt)
	}
	if len(contenido) > 16 {
		fmt.Println("Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("No se ha encontrado el disco.")
		return
	}
	for i := 0; i < len(contenido); i++ {

		var fbAux Structs.BloquesArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WriteBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], contenido[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		WriteBytes(file, bin6.Bytes())

	}

	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WriteBytes(file, inodos.Bytes())

	fmt.Println("Grupo " + rmgrp.Name + ", eliminado correctamente!")
	if super.S_filesystem_type == 3 {
		var journalingbytes bytes.Buffer
		var journaling Structs.Journaling
		copy(journaling.Operacion[:], "rmgrp")
		copy(journaling.Ruta[:], "/user.txt")
		fecha := time.Now().String()
		copy(journaling.Fecha[:], fecha)
		copy(journaling.Contenido[:], "grp: "+rmgrp.Name)
		journaling.Active = '1'
		file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
		binary.Write(&journalingbytes, binary.BigEndian, journaling)
		WriteBytes(file, journalingbytes.Bytes())
	}

	file.Close()

}
