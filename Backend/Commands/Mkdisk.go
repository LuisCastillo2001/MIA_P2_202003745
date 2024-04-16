package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var alfabeto = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

type MKDISK struct {
	size       int64
	fit        string
	unit       byte
	parameters []string
	mbr        Structs.MBR
	filename   string
	Error      string
}

func NewMKDisk(parameters []string) *MKDISK {
	mkdisk := &MKDISK{
		size:       0,
		fit:        "FF",
		unit:       'M',
		parameters: parameters,
		Error:      "",
	}
	mkdisk.readParameters()
	return mkdisk
}

func (mkdisk *MKDISK) readParameters() {

	for _, parametro := range mkdisk.parameters {
		mkdisk.identifyParameters(parametro)
	}
	if mkdisk.size == 0 {
		mkdisk.Error = "No se puede crear un disco sin su tamaño o hay un parametro mal escrito"
		return
	}
	mkdisk.CreateDisk()
}

func (mkdisk *MKDISK) identifyParameters(parameter string) {
	parameter_identifier := strings.Split(parameter, "=")

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "size" {
		sizeInt, err := strconv.ParseInt(Stringmake(parameter_identifier[1]), 10, 64)
		if err != nil {
			fmt.Println("Error al convertir a int32:", err)
			mkdisk.Error = "Solo se aceptan enteros"
			return
		}

		mkdisk.size = int64(sizeInt)
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "unit" {
		mkdisk.unit = []byte(Stringmake(strings.ToUpper(parameter_identifier[1])))[0]
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "fit" {
		mkdisk.fit = strings.ToUpper(Stringmake(parameter_identifier[1]))
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) != "size" && strings.ToLower(strings.TrimSpace(parameter_identifier[0])) != "unit" && strings.ToLower(strings.TrimSpace(parameter_identifier[0])) != "fit" {
		mkdisk.Error = "El parametro " + parameter_identifier[0] + " no es valido"
		mkdisk.size = 0
		return
	}

}

func (mkdisk *MKDISK) CreateDisk() {

	if mkdisk.unit == 'K' || mkdisk.unit == 'k' {
		mkdisk.size = 1024 * mkdisk.size
	}

	if mkdisk.unit == 'M' {
		mkdisk.size = 1024 * 1024 * mkdisk.size
	}

	content := make([]byte, mkdisk.size)

	filename := FileName()
	err := os.WriteFile(filename, content, 0644)

	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		mkdisk.Error = "El archivo no se pudo crear"
		return
	}

	mkdisk.filename = filename
	mkdisk.addMBR()

}

func (mkdisk *MKDISK) addMBR() {
	fit := []byte(mkdisk.fit)
	firtsByte := fit[0]

	mkdisk.mbr = Structs.NewMBR(mkdisk.size, firtsByte)
	file, err := os.OpenFile(mkdisk.filename, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("El archivo no se pudo abrir correctamente.")
		mkdisk.Error = "El archivo no se pudo abrir correctamente"
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error al mover el puntero al principio del archivo:", err)
		mkdisk.Error = "Hubo un error al mover el puntero"
		os.Exit(1)
	}

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &mkdisk.mbr)
	WriteBytes(file, binario.Bytes())
}

func FileName() string {
	directorio := "MIA/P1"
	ext := ".dsk"

	for _, letra := range alfabeto {
		filename := letra + ext
		rutaCompleta := filepath.Join(directorio, filename)

		if _, err := os.Stat(rutaCompleta); os.IsNotExist(err) {

			if err := os.MkdirAll(directorio, 0755); err != nil {

				panic(err)
			}

			file, err := os.Create(rutaCompleta)
			if err != nil {

				panic(err)
			}
			defer file.Close()

			return rutaCompleta
		}
	}

	return ""
}

func leerMbr(path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("El archivo no se pudo abrir correctamente.")

		os.Exit(1)
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error al mover el puntero al principio del archivo:", err)

		os.Exit(1)
	}

}
