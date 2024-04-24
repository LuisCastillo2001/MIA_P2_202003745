package Methods

import (
	"Proyecto_1/Commands"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetDisks() []string {
	Commands.Logout()
	dir := "./MIA/P1"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var fileNames []string

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames
}

func SendDisks(w http.ResponseWriter, r *http.Request) {

	arreglo := GetDisks()

	jsonData, err := json.Marshal(arreglo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
