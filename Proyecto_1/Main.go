// main.go

package main

import (
	"Proyecto_1/Commands"
	"bufio"
	"fmt"
	"os"
)

func main() {
	Commands.ExploreDisks()
	Commands.ShowMounts()
	fmt.Println("Bienvenido al proyecto de " +
		"Manejo e implementación de archivos")

	for {

		fmt.Println("Introduzca el comando que desea ejecutar")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := scanner.Text()
		Commands.LeerComando(name)
	}
}
