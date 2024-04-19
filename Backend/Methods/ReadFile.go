package Methods

import (
	"Proyecto_1/Commands"
	"encoding/json"
	"net/http"
)

func MandarArchivo(w http.ResponseWriter, r *http.Request) {
	var archivo Filecont

	err := json.NewDecoder(r.Body).Decode(&archivo)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	Commands.LeerComando(archivo.FileContent)

	cadenaJSON, err := json.Marshal(Commands.Cadena)
	if err != nil {
		http.Error(w, "Error al serializar el slice Cadena a JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(cadenaJSON)
	if err != nil {
		http.Error(w, "Error al escribir la respuesta", http.StatusInternalServerError)
		return
	}
	Commands.Cadena = []string{}

}
