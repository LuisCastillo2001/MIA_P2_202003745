package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
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
	fmt.Printf("Size: %d\n", fdisk.Size)
	fmt.Printf("Driveletter: %s\n", fdisk.Driveletter)
	fmt.Printf("Name: %s\n", fdisk.Name)
	fmt.Printf("Unit: %c\n", fdisk.Unit)
	fmt.Printf("Fit: %s\n", fdisk.Fit)
	fmt.Printf("Type: %c\n", fdisk.Type)
	fmt.Printf("Add: %d\n", fdisk.Add)
	fmt.Printf("Delete: %d\n", fdisk.Delete)
	fmt.Println("--------------------------------------")
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
			fmt.Println("Error al convertir a int32:", err)
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
			fmt.Println("Error al convertir a int64:", err)
			fdisk.Error = "Solo se aceptan enteros"
			return
		}
		fdisk.Add = addInt
	}

	if len(parameter_identifier) == 1 {
		addInt, err := strconv.ParseInt(Stringmake(parameter_identifier[0]), 10, 64)
		if err != nil {
			fmt.Println("Error al convertir a int64:", err)
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
		fmt.Println("Error al abrir el archivo para lectura:", err)
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
		fmt.Println("El archivo no se pudo abrir correctamente:", err)
		os.Exit(1)
	}
	defer file.Close()

	err = binary.Write(file, binary.BigEndian, mbr)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
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
		fmt.Println("Ya no se pueden crear mas particiones")
		return
	}

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'E' && fdisk.Type == 'E' {
			fmt.Println("Ya hay una partición extendida, no se pudo crear la particion")
			return
		}
	}

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_name == Array16bytes(fdisk.Name) {
			fmt.Println("No se pudo crear la partición, el nombre esta repetido")
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
		fmt.Println("Error, no hay espacio suficiente ")
		return
	}

	RewriteMBR(&mbr, path)
	fmt.Println("Partición realizada con éxito")
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
		fmt.Println("El archivo no se pudo abrir correctamente:", err)
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
		fmt.Println("No se ha montado ninguna partición extendida, no podrá crearse la partición lógica")
		return
	}
	part_start := partition.Part_start
	contador := 0

	for {

		_, err := file.Seek(part_start, 0)
		if err != nil {
			fmt.Println("Error al establecer la posición de escritura:", err)
			os.Exit(1)
		}

		var ebr Structs.EBR
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			fmt.Println("aqui es crear")
			fmt.Println("Error al leer el EBR:", err)
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
		fmt.Println("Error al establecer la posición de escritura para el nuevo EBR:", err)
		os.Exit(1)
	}

	err = binary.Write(file, binary.BigEndian, newEbr)
	if err != nil {
		fmt.Println("Error al escribir el nuevo EBR:", err)
		os.Exit(1)
	}

	if contador > 0 {
		partNext(mbr, path, part_start)

	}

	fmt.Println("Partición lógica realizada con éxito ")
	fdisk.printParameters()

}

func partNext(mbr Structs.MBR, path string, next int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("El archivo no se pudo abrir correctamente:", err)
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
			fmt.Println("Error al establecer la posición de escritura:", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			fmt.Println("Aqui es crear la siguiente")
			fmt.Println("Error al leer el EBR:", err)
			os.Exit(1)
		}

		if ebr.Part_next == -1 {
			break
		}
		part_start = ebr.Part_next

	}

	if ebr.Part_start == 0 {
		fmt.Println("No se encontro nada")
		return
	}

	ebr.Part_next = int64(next)
	_, err = file.Seek(part_start, 0)
	if err != nil {
		fmt.Println("Error al establecer la posición de escritura para el nuevo EBR:", err)
		os.Exit(1)
	}

	err = binary.Write(file, binary.BigEndian, ebr)
	if err != nil {
		fmt.Println("Error al escribir el nuevo EBR:", err)
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
		fmt.Println("Error al abrir el archivo para lectura:", err)
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
		fmt.Println("Voy a ir a buscar a las lógicas")
		fdisk.addSpaceLogic(file)
		return
	}
	size := mbr.Partitions[indice].Part_s + fdisk.Add
	if size < 0 {
		fmt.Println("No puede quedar esapcio negativo")
		return
	}
	if indice != 3 {

		if fdisk.Add < 0 {
			mbr.Partitions[indice].Part_s += fdisk.Add

			if mbr.Partitions[indice].Part_s < 0 {
				fmt.Println("No puede quedar espacio negativo")
				return
			}
			RewriteMBR(&mbr, path)
			fmt.Println("Espacio quitado con éxito")
		} else {
			x := fdisk.checkFreespace(indice, mbr, int(size))
			if x {
				mbr.Partitions[indice].Part_s += fdisk.Add
				fmt.Println("Si se pudo añadir espacio a la particion :)")
			} else {
				fmt.Println("No se pudo añadir espacio :(")
				return
			}
		}

	} else {
		if fdisk.Add < 0 {
			mbr.Partitions[indice].Part_s += fdisk.Add
			RewriteMBR(&mbr, path)
			fmt.Println("Espacio quitado con éxito")
			return
		} else {
			mbr.Partitions[indice].Part_s += fdisk.Add
			RewriteMBR(&mbr, path)
			fmt.Println("Espacio añadido con éxito")
			return
		}

	}

	RewriteMBR(&mbr, path)
	//LeerParticiones(path, &mbr)
	fmt.Println("Se añadio/quito correctamente espacio a la partición")
}

func (fdisk *Fdisk) checkFreespace(indice int, mbr Structs.MBR, sizetoadd int) bool {
	fmt.Println("Chequeando")

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
		fmt.Println("Error al abrir el archivo para lectura:", err)
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
			fmt.Println("Error al establecer la posición de escritura:", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			fmt.Println("Error al leer el EBR:", err)
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
	fmt.Println(ebr.Part_s)
	if ebr.Part_next == -1 {
		fmt.Println(part_start)
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &ebr)
		file.Seek(part_start, 0)
		WriteBytes(file, binario.Bytes())
		fmt.Println("Extensión de partición lógica realizada con éxito")

	} else {
		fmt.Println("No se le pudo dar espacio a la partición")
	}

}

func (fdisk *Fdisk) deletePartition() {
	path := "MIA/P1/" + fdisk.Driveletter + ".dsk"
	file, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)

	if err != nil {
		fmt.Println("Error al abrir el archivo para lectura y escritura:", err)
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
		fmt.Println("No se pudo encontrar la partición, verifique el nombre")
		return
	}

	// Llenar de ceros solo la partición a eliminar
	nullBytes := make([]byte, mbr.Partitions[indice].Part_s)
	_, err = file.WriteAt(nullBytes, int64(mbr.Partitions[indice].Part_start))
	if err != nil {
		fmt.Println("Error al escribir bytes nulos:", err)
		return
	}

	// No es necesario cerrar y abrir nuevamente el archivo

	// Actualizar la partición en la MBR
	mbr.Partitions[indice] = Structs.NewPartition()

	// Actualizar la MBR en el archivo
	RewriteMBR(&mbr, path)

	fmt.Println("Partición eliminada con éxito")
}

func LeerLogicas(path string, part_start int64) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error al abrir el archivo para lectura:", err)
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
			fmt.Println("Error al establecer la posición de escritura:", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			fmt.Println("aqui es leer")
			fmt.Println("Error al leer el EBR:", err)
			os.Exit(1)
		}
		fmt.Println("---------Particion logica--------")
		fmt.Println("Part_mount:", ebr.Part_mount)
		fmt.Println("Part_fit:", ebr.Part_fit)
		fmt.Println("Part_start:", ebr.Part_start)
		fmt.Println("Part_s:", ebr.Part_s)
		fmt.Println("Part_next:", ebr.Part_next)
		fmt.Println("Part_name:", string(ebr.Part_name[:]))
		if ebr.Part_start == 0 {
			break
		}
		part_start = ebr.Part_next
	}

}

func LeerParticiones(path string, mbr *Structs.MBR) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error al abrir el archivo para lectura:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, mbr)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	fmt.Printf("--------------MBR leído------------\n")
	fmt.Printf("MBR_tamano: %d\n", mbr.MBR_tamano)
	fmt.Printf("MBR_fecha_creacion: %s\n", string(mbr.MBR_fecha_creacion[:]))
	fmt.Printf("MBR_dsk_signature: %d\n", mbr.MBR_dsk_signature)

	// Iterar sobre las particiones e imprimir todos los atributos
	for i, partition := range mbr.Partitions {
		fmt.Printf("--------------Partición %d------------\n", i+1)
		fmt.Printf("Part_status: %c\n", partition.Part_status)
		fmt.Printf("Part_type: %c\n", partition.Part_type)
		fmt.Printf("Part_fit: %c\n", partition.Part_fit)
		fmt.Printf("Part_start: %d\n", partition.Part_start)
		fmt.Printf("Part_s: %d\n", partition.Part_s)
		fmt.Printf("Part_name: %s\n", string(partition.Part_name[:]))
		fmt.Printf("Part_correlative: %d\n", partition.Part_correlative)
		fmt.Printf("Part_id: %s\n", string(partition.Part_id[:]))
	}
}
