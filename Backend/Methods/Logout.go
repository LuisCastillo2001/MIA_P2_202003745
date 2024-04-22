package Methods

import (
	"Proyecto_1/Commands"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	Commands.Logout()
}
