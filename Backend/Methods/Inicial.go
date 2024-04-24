package Methods

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Filecont struct {
	FileContent string `json:"fileContent"`
}

var router *mux.Router = mux.NewRouter()

func Iniciar() {

	//Estos son los endpoints que voy a declarar y el tipo de metodo
	router.HandleFunc("/", inicial).Methods("GET")
	router.HandleFunc("/MandarArchivo", MandarArchivo).Methods("POST")
	router.HandleFunc("/ListaDiscos", SendDisks).Methods("GET")
	router.HandleFunc("/ListaParticiones/{disk}", GetPartitions).Methods("GET")
	router.HandleFunc("/AccederParticion/{disk}/{partition}", Access).Methods("GET")
	router.HandleFunc("/Login", makeLogin).Methods("POST")
	router.HandleFunc("/Logout", Logout).Methods("POST")
	router.HandleFunc("/ObtenerDirectorios", Getdirs).Methods("POST")
	handler := allowCORS(router)
	fmt.Println("Servidor en http://localhost:3000/")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
func inicial(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>¡Hola Mundo2!</h1>")
}
