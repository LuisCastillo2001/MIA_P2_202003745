package Commands

import (
	"fmt"
	"os"
	"strings"
)

func Eliminardisco(name string) {
	namedisk := strings.Split(name, "=")
	directorydisk := "MIA/P1/" + strings.TrimSpace(namedisk[1]) + ".dsk"

	fmt.Printf("¿Desea eliminar el disco %s? (Si/No): ", directorydisk)
	var respuesta string
	fmt.Scanln(&respuesta)

	if strings.ToLower(respuesta) == "si" {
		err := os.Remove(directorydisk)
		if err != nil {
			fmt.Println("Error al eliminar el archivo:", err)
			return
		}
		fmt.Println("Disco eliminado correctamente")
	} else {
		fmt.Println("El disco no se ha eliminado.")
	}
}

/*
if len(command) == 1 {
			fmt.Println("Se leyó un comentario")
		} else {
			if command_execute == "mkdisk" {
				parameters := command[1:] // Tomar todos los elementos después del primero
				Commands.NewMKDisk(parameters)

			}

			if command_execute == "rmdisk" {
				Commands.Eliminardisco(command[1])
			}
		}
		delay(3)

func delay(secs int) {
	for i := (int32(time.Now().Unix()) + int32(secs)); int32(time.Now().Unix()) != i; time.Sleep(time.Second) {
	}
}
*/
