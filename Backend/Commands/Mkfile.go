package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type Mkfile struct {
	Path string
	R    int
	Size int
	Cont string
}

func NewMkfile(parameters []string) *Mkfile {
	mkfile := &Mkfile{
		Path: "",
		R:    0,
		Size: 0,
		Cont: "",
	}
	mkfile.readParameters(parameters)
	return mkfile
}

func (mkfile *Mkfile) readParameters(parameters []string) {
	for _, parametro := range parameters {
		mkfile.identifyParameters(parametro)
	}
	if mkfile.R == 1 {

		return
	}
	mkfile.Makefile(mkfile.Path, 0)
}

func (mkfile *Mkfile) identifyParameters(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	parameterName := strings.ToLower(strings.TrimSpace(parameterIdentifier[0]))

	if parameterName == "path" {
		mkfile.Path = Stringmake(parameterIdentifier[1])
	}

	if parameterName == "r" {
		mkfile.R = 1
	}

	if parameterName == "size" {
		size, err := strconv.Atoi(parameterIdentifier[1])
		if err != nil {
			fmt.Println("Error al convertir el tamaño:", err)
			return
		}
		mkfile.Size = size
	}

	if parameterName == "cont" {
		mkfile.Cont = Stringmake(parameterIdentifier[1])
	}
}

func (mkfile *Mkfile) Makefile(pathaux string, numero int) bool {
	var path string
	partition := getMount(Logged.Id, &path)
	if partition.Part_start == -1 {
		fmt.Println("MKUSER", "No se encontró la partición montada con el id: "+Logged.Id)
		return false
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("MKDIR", "No se ha encontrado el disco.")
		return false
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		return false
	}

	if err_ != nil {
		fmt.Println("Error al leer el archivo")
		return false
	}
	path1 := strings.Split(pathaux, "/")
	fnd := false
	fmt.Println(path1)
	fmt.Println(fnd)
	return true

}
