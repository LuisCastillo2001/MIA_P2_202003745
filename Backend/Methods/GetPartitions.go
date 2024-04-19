package Methods

import (
	"Proyecto_1/Commands"
	"Proyecto_1/Structs"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func GetPartitions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	disk := vars["disk"]

	partitions := GetPartitionsForDisk(disk)

	jsonData, err := json.Marshal(partitions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}

func GetPartitionsForDisk(disk string) []string {
	var partitions []string
	path := "MIA/P1/" + disk
	file, err := os.Open(path)
	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	if err != nil {
		Commands.Concatenar("Error al leer el MBR:")
		return partitions
	}
	var seek int64
	seek = 0
	correlativo := 1

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_start != -1 {

			partitions = append(partitions, string(mbr.Partitions[i].Part_name[:]))
			if mbr.Partitions[i].Part_type == 'E' {
				seek = mbr.Partitions[i].Part_start
			}

		}
	}

	if seek != 0 {
		var ebr Structs.EBR

		file.Seek(seek, 0)

		for {

			_, err := file.Seek(seek, 0)
			if err != nil {
				Commands.Concatenar("Error al establecer la posición de escritura:")
				os.Exit(1)
			}

			reader := bufio.NewReader(file)
			err = binary.Read(reader, binary.BigEndian, &ebr)
			if ebr.Part_start == 0 {
				return partitions
			}
			if err != nil {
				Commands.Concatenar("Error al leer el EBR:")
				os.Exit(1)
			}

			partitions = append(partitions, string(ebr.Part_name[:]))

			if ebr.Part_next == -1 {

				break
			}
			correlativo++

			seek = ebr.Part_next

		}
	}
	return partitions

}
