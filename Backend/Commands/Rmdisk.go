package Commands

import (
	"os"
	"strings"
)

func Eliminardisco(name string) {
	namedisk := strings.Split(name, "=")
	directorydisk := "MIA/P1/" + strings.TrimSpace(namedisk[1]) + ".dsk"

	respuesta := "si"

	if strings.ToLower(respuesta) == "si" {
		err := os.Remove(directorydisk)
		if err != nil {
			Concatenar("Error al eliminar el archivo:")
			return
		}
		Concatenar("Disco eliminado correctamente")
	} else {
		Concatenar("El disco no se ha eliminado.")
	}
}

/*
if len(command) == 1 {
			Concatenar("Se leyó un comentario")
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
