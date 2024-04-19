package Commands

import "fmt"

var Cadena []string

func Concatenar(nuevo string) {
	Cadena = append(Cadena, nuevo)
}

func Concatenar2(formato string, args ...interface{}) {

	cadena := fmt.Sprintf(formato, args...)

	Cadena = append(Cadena, cadena)
}
