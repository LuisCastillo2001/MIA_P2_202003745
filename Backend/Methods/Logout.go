package Methods

import (
	"Proyecto_1/Commands"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	x := Commands.Logout()

	if x == false {
		http.NotFound(w, r)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
