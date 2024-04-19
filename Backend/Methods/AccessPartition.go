package Methods

import (
	"Proyecto_1/Commands"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func Access(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	disk := vars["disk"]
	partition := vars["partition"]
	namedisk := QuitarPunto(disk)
	for i := 0; i < len(Commands.MountedPartitions); i++ {
		partName := strings.TrimSpace(string(Commands.MountedPartitions[i].PartitionName[:]))
		if partName == partition && namedisk == Commands.MountedPartitions[i].DiskName {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.NotFound(w, r)
}

func QuitarPunto(cadena string) string {
	cadena_sin_punto := strings.TrimSuffix(cadena, ".dsk")
	return cadena_sin_punto
}
