package Commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func execute(filePath string) {
	path := strings.Split(filePath, "=")
	filePath = strings.TrimSpace(path[1])
	fmt.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		LeerComando(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}
