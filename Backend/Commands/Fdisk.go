package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type Fdisk struct {
	Size        int
	Driveletter string
	Name        string
	Unit        byte
	Type        byte
	Fit         string
	Add         int64
	Parameters  []string
	Error       string
	Delete      int
}

func NewFdisk(parameters []string) *Fdisk {
	fdisk := &Fdisk{
		Size:        0,
		Driveletter: "",
		Name:        "",
		Unit:        'B',
		Fit:         "WF",
		Type:        'P',
		Parameters:  parameters,
		Add:         -1,
		Delete:      -1,
	}
	fdisk.readParameters()
	return fdisk
}

func (fdisk *Fdisk) printParameters() {
	Concatenar2("Size: %d\n", fdisk.Size)
	Concatenar2("Driveletter: %s\n", fdisk.Driveletter)
	Concatenar2("Name: %s\n", fdisk.Name)
	Concatenar2("Unit: %c\n", fdisk.Unit)
	Concatenar2("Fit: %s\n", fdisk.Fit)
	Concatenar2("Type: %c\n", fdisk.Type)
	Concatenar2("Add: %d\n", fdisk.Add)
	Concatenar2("Delete: %d\n", fdisk.Delete)
	Concatenar("--------------------------------------")
}

func (fdisk *Fdisk) readParameters() {
	for _, parametro := range fdisk.Parameters {
		fdisk.identifyParameters(parametro)

	}

	if fdisk.Add != -1 {
		fdisk.addP_E()
		return
	}
	if fdisk.Delete != -1 {
		fdisk.deletePartition()
		return
	}

	fdisk.createPartition()
}

func (fdisk *Fdisk) identifyParameters(parameter string) {

	parameter_identifier := strings.Split(parameter, "=")

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "size" {
		sizeInt, err := strconv.Atoi(Stringmake(parameter_identifier[1]))
		if err != nil {
			Concatenar("Error al convertir a int32:")
			fdisk.Error = "Solo se aceptan enteros"
			return
		}
		fdisk.Size = sizeInt
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "unit" {
		fdisk.Unit = []byte(strings.ToUpper(Stringmake(parameter_identifier[1])))[0]
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "fit" {
		fdisk.Fit = strings.ToUpper(Stringmake(parameter_identifier[1]))
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "name" {
		fdisk.Name = Stringmake(parameter_identifier[1])
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "driveletter" {

		fdisk.Driveletter = strings.ToUpper(Stringmake(parameter_identifier[1]))
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "type" {
		fdisk.Type = []byte(strings.ToUpper(Stringmake(parameter_identifier[1])))[0]
	}
	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "add" {
		if Stringmake(parameter_identifier[1]) == "" {

			return
		}
		addInt, err := strconv.ParseInt(Stringmake(parameter_identifier[1]), 10, 64)
		if err != nil {
			Concatenar("Error al convertir a int64:")
			fdisk.Error = "Solo se aceptan enteros"
			return
		}
		fdisk.Add = addInt
	}

	if len(parameter_identifier) == 1 {
		addInt, err := strconv.ParseInt(Stringmake(parameter_identifier[0]), 10, 64)
		if err != nil {
			Concatenar("Error al convertir a int64:")
			fdisk.Error = "Solo se aceptan enteros"
			return
		}
		fdisk.Add = addInt * -1
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "delete" {
		fdisk.Delete = 1
	}

}

func (fdisk *Fdisk) createPartition() {
	path := "MIA/P1/" + fdisk.Driveletter + ".dsk"
	file, err := os.Open(path)
	if err != nil {
		Concatenar("Error al abrir el archivo para lectura:")
		return
	}
	defer file.Close()
	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	if fdisk.Type == 'E' || fdisk.Type == 'P' {
		fdisk.CreateP_E(mbr, path)
	}
	if fdisk.Type == 'L' {
		fdisk.createLogic(mbr, path)
	}
	//LeerParticiones(path, &mbr)

}

func RewriteMBR(mbr *Structs.MBR, path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		Concatenar("El archivo no se pudo abrir correctamente:")
		os.Exit(1)
	}
	defer file.Close()

	err = binary.Write(file, binary.BigEndian, mbr)
	if err != nil {
		Concatenar("Error al escribir en el archivo:")
		os.Exit(1)
	}

}

func (fdisk *Fdisk) CreateP_E(mbr Structs.MBR, path string) {

	if fdisk.Unit == 'K' {

		fdisk.Size = 1024 * fdisk.Size
	}
	if fdisk.Unit == 'M' {
		fdisk.Size = 1024 * 1024 * fdisk.Size
	}

	indice := 0
	part_start := 0

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_status == '0' {
			indice = i
			if indice == 0 {
				part_start = int(unsafe.Sizeof(mbr))
			} else {
				size := mbr.Partitions[indice-1].Part_start + mbr.Partitions[indice-1].Part_s
				part_start = int(size)

			}
			break
		}
	}
	contador := 1
	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_status == '1' {
			contador++
		}
	}

	if contador == 5 {
		Concatenar("Ya no se pueden crear mas particiones")
		return
	}

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'E' && fdisk.Type == 'E' {
			Concatenar("Ya hay una partición extendida, no se pudo crear la particion")
			return
		}
	}

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_name == Array16bytes(fdisk.Name) {
			Concatenar("No se pudo crear la partición, el nombre esta repetido")
			return
		}
	}
	fit := []byte(fdisk.Fit)
	firtsByte := fit[0]

	mbr.Partitions[indice].Part_status = '1'
	mbr.Partitions[indice].Part_name = Array16bytes(fdisk.Name)
	mbr.Partitions[indice].Part_fit = firtsByte
	mbr.Partitions[indice].Part_type = fdisk.Type
	mbr.Partitions[indice].Part_start = int64(part_start)
	mbr.Partitions[indice].Part_s = int64(fdisk.Size)
	mbr.Partitions[indice].Part_correlative = int64(indice)

	if (mbr.Partitions[indice].Part_start + mbr.Partitions[indice].Part_s) > mbr.MBR_tamano {
		Concatenar("Error, no hay espacio suficiente ")
		return
	}

	RewriteMBR(&mbr, path)
	Concatenar("Partición realizada con éxito")
	fdisk.printParameters()

}

func (fdisk *Fdisk) createLogic(mbr Structs.MBR, path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if fdisk.Unit == 'K' {
		fdisk.Size = 1024 * fdisk.Size
	}
	if fdisk.Unit == 'M' {
		fdisk.Size = 1024 * 1024 * fdisk.Size
	}
	if err != nil {
		Concatenar("El archivo no se pudo abrir correctamente:")
		os.Exit(1)
	}

	defer file.Close()

	partition := Structs.NewPartition()

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'E' {
			partition = mbr.Partitions[i]

			break
		}
	}
	if partition.Part_start == -1 {
		Concatenar("No se ha montado ninguna partición extendida, no podrá crearse la partición lógica")
		return
	}
	part_start := partition.Part_start
	contador := 0

	for {

		_, err := file.Seek(part_start, 0)
		if err != nil {
			Concatenar("Error al establecer la posición de escritura:")
			os.Exit(1)
		}

		var ebr Structs.EBR
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			Concatenar("aqui es crear")
			Concatenar("Error al leer el EBR:")
			os.Exit(1)
		}

		if ebr.Part_start == 0 {

			break
		}

		part_start = ebr.Part_start + ebr.Part_s
		contador++

	}
	fit := []byte(fdisk.Fit)
	firtsByte := fit[0]

	newEbr := Structs.NewEBR()
	newEbr.Part_start = part_start
	newEbr.Part_fit = firtsByte
	newEbr.Part_mount = '0'
	newEbr.Part_s = int64(fdisk.Size)
	newEbr.Part_name = Array16bytes(fdisk.Name)

	_, err = file.Seek(part_start, 0)
	if err != nil {
		Concatenar("Error al establecer la posición de escritura para el nuevo EBR:")
		os.Exit(1)
	}

	err = binary.Write(file, binary.BigEndian, newEbr)
	if err != nil {
		Concatenar("Error al escribir el nuevo EBR:")
		os.Exit(1)
	}

	if contador > 0 {
		partNext(mbr, path, part_start)

	}

	Concatenar("Partición lógica realizada con éxito ")
	fdisk.printParameters()

}

func partNext(mbr Structs.MBR, path string, next int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		Concatenar("El archivo no se pudo abrir correctamente:")
		os.Exit(1)
	}
	defer file.Close()

	partition := Structs.NewPartition()
	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'E' {
			partition = mbr.Partitions[i]
			break
		}
	}

	part_start := partition.Part_start
	var ebr Structs.EBR
	for {

		_, err := file.Seek(part_start, 0)
		if err != nil {
			Concatenar("Error al establecer la posición de escritura:")
			os.Exit(1)
		}

		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			Concatenar("Aqui es crear la siguiente")
			Concatenar("Error al leer el EBR:")
			os.Exit(1)
		}

		if ebr.Part_next == -1 {
			break
		}
		part_start = ebr.Part_next

	}

	if ebr.Part_start == 0 {
		Concatenar("No se encontro nada")
		return
	}

	ebr.Part_next = int64(next)
	_, err = file.Seek(part_start, 0)
	if err != nil {
		Concatenar("Error al establecer la posición de escritura para el nuevo EBR:")
		os.Exit(1)
	}

	err = binary.Write(file, binary.BigEndian, ebr)
	if err != nil {
		Concatenar("Error al escribir el nuevo EBR:")
		os.Exit(1)
	}

}

func (fdisk *Fdisk) addP_E() {
	path := "MIA/P1/" + fdisk.Driveletter + ".dsk"
	file, err := os.Open(path)
	if fdisk.Unit == 'K' {
		fdisk.Add = 1024 * fdisk.Add
	}
	if fdisk.Unit == 'M' {
		fdisk.Add = 1024 * 1024 * fdisk.Add
	}
	if err != nil {
		Concatenar("Error al abrir el archivo para lectura:")
		return
	}
	defer file.Close()
	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	indice := -1

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_name == Array16bytes(fdisk.Name) {
			indice = i

			break
		}
	}
	if indice == -1 {
		Concatenar("Voy a ir a buscar a las lógicas")
		fdisk.addSpaceLogic(file)
		return
	}
	size := mbr.Partitions[indice].Part_s + fdisk.Add
	if size < 0 {
		Concatenar("No puede quedar esapcio negativo")
		return
	}
	if indice != 3 {

		if fdisk.Add < 0 {
			mbr.Partitions[indice].Part_s += fdisk.Add

			if mbr.Partitions[indice].Part_s < 0 {
				Concatenar("No puede quedar espacio negativo")
				return
			}
			RewriteMBR(&mbr, path)
			Concatenar("Espacio quitado con éxito")
		} else {
			x := fdisk.checkFreespace(indice, mbr, int(size))
			if x {
				mbr.Partitions[indice].Part_s += fdisk.Add
				Concatenar("Si se pudo añadir espacio a la particion :)")
			} else {
				Concatenar("No se pudo añadir espacio :(")
				return
			}
		}

	} else {
		if fdisk.Add < 0 {
			mbr.Partitions[indice].Part_s += fdisk.Add
			RewriteMBR(&mbr, path)
			Concatenar("Espacio quitado con éxito")
			return
		} else {
			mbr.Partitions[indice].Part_s += fdisk.Add
			RewriteMBR(&mbr, path)
			Concatenar("Espacio añadido con éxito")
			return
		}

	}

	RewriteMBR(&mbr, path)
	//LeerParticiones(path, &mbr)
	Concatenar("Se añadio/quito correctamente espacio a la partición")
}

func (fdisk *Fdisk) checkFreespace(indice int, mbr Structs.MBR, sizetoadd int) bool {
	Concatenar("Chequeando")

	var diferencia int64
	diferencia = 0

	if mbr.Partitions[indice].Part_start+mbr.Partitions[indice].Part_s+fdisk.Add < mbr.Partitions[indice+1].Part_start && mbr.Partitions[indice+1].Part_status == '1' {
		return true
	}

	for i := indice + 1; i < 4; i++ {
		if mbr.Partitions[i].Part_status == '1' {
			diferencia = mbr.Partitions[i].Part_start
			break
		}
	}
	if diferencia != 0 {
		if fdisk.Add+mbr.Partitions[indice].Part_s+mbr.Partitions[indice].Part_start > diferencia {
			return false
		} else {

			return true
		}
	} else {
		if mbr.Partitions[indice].Part_start+mbr.Partitions[indice].Part_s+fdisk.Add > mbr.MBR_tamano {
			return false
		} else {
			return true
		}

	}
	return false
}

func (fdisk *Fdisk) addSpaceLogic(file *os.File) {
	path := "MIA/P1/" + fdisk.Driveletter + ".dsk"
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if fdisk.Unit == 'K' {
		fdisk.Add = 1024 * fdisk.Add
	}
	if fdisk.Unit == 'M' {
		fdisk.Add = 1024 * 1024 * fdisk.Add
	}
	if err != nil {
		Concatenar("Error al abrir el archivo para lectura:")
		return
	}
	defer file.Close()
	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)

	indice := -1
	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'E' {
			indice = i
			break
		}
	}

	if indice == -1 {
		return
	}
	part_start := mbr.Partitions[indice].Part_start

	var ebr Structs.EBR

	file.Seek(part_start, 0)

	for {
		_, err := file.Seek(part_start, 0)
		if err != nil {
			Concatenar("Error al establecer la posición de escritura:")
			os.Exit(1)
		}

		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			Concatenar("Error al leer el EBR:")
			os.Exit(1)
		}

		if ebr.Part_name == Array16bytes(fdisk.Name) {
			break
		}

		if ebr.Part_start == 0 {
			break
		}
		part_start = ebr.Part_next

	}

	ebr.Part_s += fdisk.Add

	if ebr.Part_next == -1 {

		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &ebr)
		file.Seek(part_start, 0)
		WriteBytes(file, binario.Bytes())
		Concatenar("Extensión de partición lógica realizada con éxito")

	} else {
		Concatenar("No se le pudo dar espacio a la partición")
	}

}

func (fdisk *Fdisk) deletePartition() {
	path := "MIA/P1/" + fdisk.Driveletter + ".dsk"
	file, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)

	if err != nil {
		Concatenar("Error al abrir el archivo para lectura y escritura:")
		return
	}
	defer file.Close()

	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)

	indice := -1
	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_name == Array16bytes(fdisk.Name) {
			indice = i
			break
		}
	}
	if indice == -1 {
		Concatenar("No se pudo encontrar la partición, verifique el nombre")
		return
	}

	// Llenar de ceros solo la partición a eliminar
	nullBytes := make([]byte, mbr.Partitions[indice].Part_s)
	_, err = file.WriteAt(nullBytes, int64(mbr.Partitions[indice].Part_start))
	if err != nil {
		Concatenar("Error al escribir bytes nulos:")
		return
	}

	// No es necesario cerrar y abrir nuevamente el archivo

	// Actualizar la partición en la MBR
	mbr.Partitions[indice] = Structs.NewPartition()

	// Actualizar la MBR en el archivo
	RewriteMBR(&mbr, path)

	Concatenar("Partición eliminada con éxito")
}

/*
	func LeerLogicas(path string, part_start int64) {
		file, err := os.Open(path)
		if err != nil {
			Concatenar("Error al abrir el archivo para lectura:")
			return
		}
		defer file.Close()
		file.Seek(part_start, 0)

		var ebr Structs.EBR
		for {
			if part_start == -1 {
				break
			}
			_, err := file.Seek(part_start, 0)
			if err != nil {
				Concatenar("Error al establecer la posición de escritura:")
				os.Exit(1)
			}

			reader := bufio.NewReader(file)
			err = binary.Read(reader, binary.BigEndian, &ebr)
			if err != nil {
				Concatenar("aqui es leer")
				Concatenar("Error al leer el EBR:")
				os.Exit(1)
			}
			Concatenar("---------Particion logica--------")
			Concatenar("Part_mount:", ebr.Part_mount)
			Concatenar("Part_fit:", ebr.Part_fit)
			Concatenar("Part_start:", ebr.Part_start)
			Concatenar("Part_s:", ebr.Part_s)
			Concatenar("Part_next:", ebr.Part_next)
			Concatenar("Part_name:", string(ebr.Part_name[:]))
			if ebr.Part_start == 0 {
				break
			}
			part_start = ebr.Part_next
		}

}
*/
func LeerParticiones(path string, mbr *Structs.MBR) {
	file, err := os.Open(path)
	if err != nil {
		Concatenar("Error al abrir el archivo para lectura:")
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, mbr)
	if err != nil {
		Concatenar("Error al leer el archivo:")
		return
	}

	Concatenar2("--------------MBR leído------------\n")
	Concatenar2("MBR_tamano: %d\n", mbr.MBR_tamano)
	Concatenar2("MBR_fecha_creacion: %s\n", string(mbr.MBR_fecha_creacion[:]))
	Concatenar2("MBR_dsk_signature: %d\n", mbr.MBR_dsk_signature)

	// Iterar sobre las particiones e imprimir todos los atributos
	for i, partition := range mbr.Partitions {
		Concatenar2("--------------Partición %d------------\n", i+1)
		Concatenar2("Part_status: %c\n", partition.Part_status)
		Concatenar2("Part_type: %c\n", partition.Part_type)
		Concatenar2("Part_fit: %c\n", partition.Part_fit)
		Concatenar2("Part_start: %d\n", partition.Part_start)
		Concatenar2("Part_s: %d\n", partition.Part_s)
		Concatenar2("Part_name: %s\n", string(partition.Part_name[:]))
		Concatenar2("Part_correlative: %d\n", partition.Part_correlative)
		Concatenar2("Part_id: %s\n", string(partition.Part_id[:]))
	}
}
