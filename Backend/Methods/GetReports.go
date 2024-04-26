package Methods

import (
	"io/ioutil"
	"path/filepath"
)

type Reporte struct {
	name string
	dot  string
}

func ObtenerReportes() []Reporte {

	var reportes []Reporte
	files, err := ioutil.ReadDir("./MIA/Reportes")
	if err != nil {
		return reportes
	}

	for _, file := range files {

		name := file.Name()

		content, err := ioutil.ReadFile(filepath.Join("./MIA/Reportes", name))
		if err != nil {
			return reportes
		}

		reporte := Reporte{
			name: name,
			dot:  string(content),
		}
		reportes = append(reportes, reporte)
	}

	return reportes
}
