package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type Unmount struct {
	Id         string
	Parameters []string
}

func NewUnmount(parameters []string) *Unmount {
	unmount := &Unmount{
		Id:         "",
		Parameters: parameters,
	}
	unmount.readParameters()
	return unmount
}

func (unmount *Unmount) readParameters() {
	for _, parametro := range unmount.Parameters {
		unmount.identifyParameters(parametro)
	}
	unmount.UnmountPartition()
}

func (unmount *Unmount) identifyParameters(parameter string) {
	parameter_identifier := strings.Split(parameter, "=")

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "id" {
		unmount.Id = Stringmake(parameter_identifier[1])
	}
}

func (unmount *Unmount) UnmountPartition() {
	path := ""
	partition_unmount := getMount(unmount.Id, &path)

	if partition_unmount.Part_start == -1 {
		fmt.Println("No se encontro la particion o no ha sido montada")
		return
	} else {
		for i := range MountedPartitions {
			if MountedPartitions[i].PartitionName == partition_unmount.Part_name {

				UnmountinDisk(MountedPartitions[i].DiskName, partition_unmount.Part_start, partition_unmount.Part_type)
				MountedPartitions[i] = MountedPartition{}
				break

			}
		}
		fmt.Println("La particion ha sido desmontada con éxito")
		ShowMounts()
	}

}

func UnmountinDisk(driveletter string, part_start int64, part_type byte) {

	path := "MIA/P1/" + driveletter + ".dsk"
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo para lectura:", err)
		return
	}
	defer file.Close()

	if part_type == 'P' {
		var mbr Structs.MBR
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &mbr)

		indice := 0
		for i := range mbr.Partitions {
			if part_start == mbr.Partitions[i].Part_start {
				indice = i
				break
			}
		}
		mbr.Partitions[indice].Part_id = [4]byte{}

		RewriteMBR(&mbr, path)
		return
	}

	if part_type == 'L' {

		var ebr Structs.EBR

		_, err := file.Seek(part_start, 0)
		if err != nil {
			fmt.Println("Error al establecer la posición de escritura:", err)
			os.Exit(1)
		}
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)

		ebr.Part_mount = 0
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &ebr)
		file.Seek(part_start, 0)
		WriteBytes(file, binario.Bytes())

		return
	}

}
