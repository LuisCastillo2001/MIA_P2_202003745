package Methods

import (
	"Proyecto_1/Commands"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func MandarArchivo(w http.ResponseWriter, r *http.Request) {
	var archivo Filecont


	err := json.NewDecoder(r.Body).Decode(&archivo)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	lineas := strings.Split(archivo.FileContent, "\n")

	// Recorrer cada línea e imprimir
	for _, linea := range lineas {
		Commands.LeerComando(linea)
	}
	fmt.Println("fin")

}
