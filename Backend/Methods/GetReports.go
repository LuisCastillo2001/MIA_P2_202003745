package Methods

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type Reporte struct {
	Name string `json:"name"`
	Dot  string `json:"dot"`
}

func ObtenerReportes() []Reporte {
	var reportes []Reporte
	files, err := ioutil.ReadDir("./MIA/Reportes")
	if err != nil {
		return reportes
	}

	for _, file := range files {
		name := file.Name()

		// Leer el contenido del archivo
		filePath := filepath.Join("./MIA/Reportes", name)
		contentBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		content := string(contentBytes)

		reporte := Reporte{
			Name: name,
			Dot:  content,
		}
		reportes = append(reportes, reporte)
	}

	return reportes
}

func SendReportes(w http.ResponseWriter, r *http.Request) {
	reportes := ObtenerReportes()

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(reportes)
}
