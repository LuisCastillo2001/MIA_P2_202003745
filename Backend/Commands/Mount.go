package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
)

type MountedPartition struct {
	DiskName      string
	Id            string
	PartitionName [16]byte
}

var MountedPartitions [50]MountedPartition

type Mount struct {
	Driveletter string
	Name        [16]byte
	id          string
	Parameters  []string
	Correlativo int
}

func NewMount(parameters []string) *Mount {
	mount := &Mount{
		Driveletter: "",
		Name:        [16]byte{},
		Parameters:  parameters,
		Correlativo: 1,
	}
	mount.readParameters()
	return mount
}

func (mount *Mount) readParameters() {
	for _, parametro := range mount.Parameters {
		mount.identifyParameters(parametro)

	}

	mount.MountPartition()

}

func (mount *Mount) identifyParameters(parameter string) {
	parameter_identifier := strings.Split(parameter, "=")

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "driveletter" {
		mount.Driveletter = strings.ToUpper(Stringmake(parameter_identifier[1]))
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "name" {
		mount.Name = Array16bytes(Stringmake(parameter_identifier[1]))
	}

}

func (mount *Mount) MountPartition() {
	path := "MIA/P1/" + mount.Driveletter + ".dsk"
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		Concatenar("Error al abrir el archivo para lectura:")
		return
	}
	defer file.Close()

	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	indice := -1
	var seek int64
	seek = 0

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_start != -1 {
			if mbr.Partitions[i].Part_type == 'E' {
				seek = mbr.Partitions[i].Part_start
			}
			if mbr.Partitions[i].Part_name == mount.Name {
				indice = i
				break
			}
			mount.Correlativo++
		}

	}

	if indice != -1 {
		if mbr.Partitions[indice].Part_type == 'E' {
			Concatenar("La particion no puede ser montada")
			return
		} else if mbr.Partitions[indice].Part_type == 'P' {
			comprobar := mbr.Partitions[indice].Part_id

			x := 0
			for i := range comprobar {
				if comprobar[i] == 0 {
					x++
				}
			}

			if x != 4 {
				Concatenar("La particion ya ha sido montada")
				ShowMounts()
				return
			}
			mount.Addmount()
			copy(mbr.Partitions[indice].Part_id[:], mount.id)
			RewriteMBR(&mbr, path)
			Concatenar("Particion Montada con éxito")

		}
	}

	if seek != 0 && indice == -1 {

		var ebr Structs.EBR

		file.Seek(seek, 0)

		for {
			_, err := file.Seek(seek, 0)
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

			if ebr.Part_name == mount.Name {
				break
			}

			if ebr.Part_start == -1 {
				break
			}
			seek = ebr.Part_next
			mount.Correlativo++
		}

		if ebr.Part_mount == 1 {
			Concatenar("La partición lógica ya fue montada")
			return
		}

		ebr.Part_mount = 1
		file.Seek(ebr.Part_start, 0)
		err = binary.Write(file, binary.BigEndian, ebr)
		mount.Addmount()
		Concatenar("Particion lógica Montada con éxito")

	}

}

func (mount *Mount) Addmount() {

	indice := 0
	for i := range MountedPartitions {
		if MountedPartitions[i].Id == "" {
			indice = i
			break
		}
	}

	mount.id = mount.Driveletter + strconv.Itoa(mount.Correlativo) + "45"

	newmount := MountedPartition{
		DiskName:      mount.Driveletter,
		Id:            mount.id,
		PartitionName: mount.Name,
	}
	MountedPartitions[indice] = newmount
	ShowMounts()

}

func ShowMounts() {
	Concatenar("-------Particiones montadas---------")
	for i := range MountedPartitions {
		if MountedPartitions[i].Id != "" {
			Concatenar2("Disk: %s, Id: %s, PartitionName: %s\n",
				MountedPartitions[i].DiskName,
				MountedPartitions[i].Id,
				string(MountedPartitions[i].PartitionName[:]),
			)
		}

	}
}

func getMount(id string, path2 *string) *Structs.Partition {
	indice := -1
	var tmppartition Structs.Partition
	tmppartition.Part_start = -1
	for i := range MountedPartitions {
		if MountedPartitions[i].Id == id {
			indice = i
			break
		}
	}

	if indice == -1 {
		tmppartition.Part_start = -1
		return &tmppartition
	}
	path := "MIA/P1/" + MountedPartitions[indice].DiskName + ".dsk"
	*path2 = path

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		Concatenar("Error al abrir el archivo para lectura:")
		return &tmppartition
	}
	defer file.Close()

	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	indiceaux := -1
	var seek int64
	seek = 0

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'E' {
			seek = mbr.Partitions[i].Part_start
		}
		if mbr.Partitions[i].Part_name == MountedPartitions[indice].PartitionName {
			indiceaux = i
			tmppartition = mbr.Partitions[indiceaux]

			return &tmppartition
		}

	}

	if seek != 0 && indiceaux == -1 {

		var ebr = findLogic(file, MountedPartitions[indice].PartitionName, seek)
		if ebr == nil {
			Concatenar("No se pudo encontrar la particion, verifique que no la haya borrado")
			return &tmppartition
		}
		tmppartition.Part_type = 'L'
		tmppartition.Part_name = ebr.Part_name
		tmppartition.Part_status = '1'
		tmppartition.Part_start = ebr.Part_start
		tmppartition.Part_fit = ebr.Part_fit
		tmppartition.Part_s = ebr.Part_s
		tmppartition.Part_status = '1'

		return &tmppartition

	}
	return &tmppartition

}
