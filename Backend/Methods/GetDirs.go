package Methods

import (
	"Proyecto_1/Commands"
	"encoding/json"
	"net/http"
	"strings"
)

type GetDirs struct {
	Diskname  string `json:"diskname"`
	Partition string `json:"partition"`
	Path      string `json:"path"`
}

func Getdirs(w http.ResponseWriter, r *http.Request) {
	var getdirs GetDirs
	err := json.NewDecoder(r.Body).Decode(&getdirs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var x bool
	if strings.HasSuffix(getdirs.Path, "/") && len(getdirs.Path) > 1 {
		getdirs.Path = getdirs.Path[:len(getdirs.Path)-1]
	}

	dirs := Commands.ReturnDirs(getdirs.Diskname, getdirs.Path, &x)
	if x == false {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(dirs)
}
