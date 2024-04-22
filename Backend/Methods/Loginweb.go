package Methods

import (
	"Proyecto_1/Commands"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Id       string `json:"id"`
}

func makeLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login := Commands.Login{
		User: user.User,
		Pass: user.Password,
		Id:   user.Id,
	}
	fmt.Println(user.Id)
	fmt.Println(user.Password)
	fmt.Println(user.User)
	x := login.Makelogin()
	fmt.Println(x)
	if x == true {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.NotFound(w, r)
		return
	}
}
