package Methods

import (
	"Proyecto_1/Commands"
	"Proyecto_1/Structs"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Partition_web struct {
	ID            string `json:"id"`
	DiskName      string `json:"diskname"`
	PartitionName string `json:"partitionname"`
}

func Access(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	disk := vars["disk"]
	partition := vars["partition"]
	namedisk := quitarPunto(disk)
	for i := 0; i < len(Commands.MountedPartitions); i++ {
		partName := strings.TrimSpace(string(Commands.MountedPartitions[i].PartitionName[:]))
		if partName == partition && namedisk == Commands.MountedPartitions[i].DiskName {
			x := verifymkfs(partition, namedisk)
			if x == false {
				http.NotFound(w, r)
				return
			}
			part := Partition_web{
				ID:            Commands.MountedPartitions[i].Id,
				DiskName:      Commands.MountedPartitions[i].DiskName,
				PartitionName: partName,
			}
			json.NewEncoder(w).Encode(part)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.NotFound(w, r)
}

func quitarPunto(cadena string) string {
	cadena_sin_punto := strings.TrimSuffix(cadena, ".dsk")
	return cadena_sin_punto
}

func verifymkfs(partition string, disk string) bool {
	directory := "MIA/P1/" + disk + ".dsk"
	file, err := os.OpenFile(directory, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()
	//Leer el MBR
	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)

	if err != nil {
		fmt.Println(err)
		return false
	}
	//Obtner la particion

	var seek int64
	seek = 0
	for i := 0; i < 4; i++ {
		partitionname := strings.TrimSpace(string(mbr.Partitions[i].Part_name[:]))

		if partitionname == partition {
			seek = int64(mbr.Partitions[i].Part_start)
			break
		}
	}

	//Leer si existe un superbloque

	file.Seek(int64(seek), 0)
	superbloque := Structs.NewSuperBloque()
	reader = bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	if err != nil {

		return false
	}

	if superbloque.S_magic == 0xEF53 {
		return true
	}
	return false

}
