package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func LeerComando(name string) {

	for {

		command := strings.Split(name, "-")
		command_execute := strings.ToLower(strings.ReplaceAll(command[0], " ", ""))

		if len(command) == 1 && (command_execute == "pause" || strings.HasPrefix(name, "#")) {
			if command_execute == "pause" {
				fmt.Println("Ejecutando pausa...")
				fmt.Scanln()
			} else {
				fmt.Println("Se leyó un comentario.")
			}
		} else if command_execute == "mounted" {
			ShowMounts()
		} else if command_execute == "logout" {
			Logout()
		} else {
			if command_execute == "mkdisk" {
				parameters := command[1:] // Tomar todos los elementos después del primero
				fmt.Println("-------COMANDO MKDISK--------")
				x := NewMKDisk(parameters)
				if x.Error != "" {
					fmt.Println(x.Error)
				} else {
					fmt.Println("DISCO CREADO CORRECTAMENTE")
				}

			}

			if command_execute == "rmdisk" {
				fmt.Println("--------COMANDO RMDISK-------")
				Eliminardisco(command[1])
			}

			if command_execute == "fdisk" {
				fmt.Println("---------COMANDO FDISK-------")
				parameters := command[1:]
				NewFdisk(parameters)
			}

			if command_execute == "mount" {
				fmt.Println("----------COMANDO MOUNT---------")
				parameters := command[1:]
				NewMount(parameters)
			}

			if command_execute == "unmount" {
				fmt.Println("---------COMANDO UNMOUNT---------")
				parameters := command[1:]
				NewUnmount(parameters)
			}

			if command_execute == "mkfs" {
				fmt.Println("------------COMANDO MKFS-----------")
				parameters := command[1:]
				NewMkfs(parameters)
			}

			if command_execute == "login" {
				fmt.Println("----------COMANDO LOGIN------------")
				parameters := command[1:]
				NewLogin(parameters)
			}

			if command_execute == "mkgrp" {
				fmt.Println("----------COMANDO MKGRP------------")
				parameters := command[1:]
				NewMkgrp(parameters)
			}

			if command_execute == "rmgrp" {
				fmt.Println("-----------COMANDO RMGRP-----------")
				parameters := command[1:]
				newRmgrp(parameters)
			}

			if command_execute == "mkusr" {
				fmt.Println("-----------COMANDO MKUSR-----------")
				parameters := command[1:]
				NewMkusr(parameters)
			}

			if command_execute == "rmusr" {
				fmt.Println("-----------COMANDO RMUSR-----------")
				parameters := command[1:]
				NewRmuser(parameters)
			}

			if command_execute == "rep" {
				fmt.Println("--------GENERACIÓN DE REPORTES-------")
				parameters := command[1:]
				NewRep(parameters)
			}

			if command_execute == "mkdir" {
				fmt.Println("---------COMANDO MKDIR---------------")
				parameters := command[1:]
				NewMkdir(parameters)
			}

			if command_execute == "execute" {
				execute(command[1])
			}

		}
		delay(2)
		return
	}
}

func delay(secs int) {
	for i := (int32(time.Now().Unix()) + int32(secs)); int32(time.Now().Unix()) != i; time.Sleep(time.Second) {
	}
}

func Stringmake(str string) string {

	trimmedStr := strings.TrimSpace(str)

	result := strings.ReplaceAll(trimmedStr, "\"", "")

	return result
}
func WriteBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func getPath(path string) []string {
	var result []string
	if path == "" {
		return result
	}
	aux := strings.Split(path, "/")
	for i := 1; i < len(aux); i++ {
		result = append(result, aux[i])
	}
	return result
}

func Array16bytes(valor string) [16]byte {

	valorSinEspacios := strings.TrimSpace(valor)

	valorSinComillas := strings.ReplaceAll(valorSinEspacios, "\"", "")

	if len(valorSinComillas) > 16 {
		valorSinComillas = valorSinComillas[:16]
	} else {
		valorSinComillas = fmt.Sprintf("%-16s", valorSinComillas)
	}

	var arregloBytes [16]byte
	copy(arregloBytes[:], []byte(valorSinComillas))

	return arregloBytes
}

func BytesToString(bytes []byte) string {

	str := string(bytes[:])

	str = strings.TrimSpace(str)

	str = strings.ReplaceAll(str, "\"", "")

	return str
}

func findLogic(file *os.File, name [16]byte, seek int64) *Structs.EBR {
	var ebr Structs.EBR
	file.Seek(seek, 0)
	flag := false
	for {

		_, err := file.Seek(seek, 0)
		if err != nil {
			fmt.Println("Error al establecer la posición de escritura:", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &ebr)
		if err != nil {
			fmt.Println("Error al leer el EBR:", err)
			os.Exit(1)
		}

		if ebr.Part_name == name {
			flag = true
			break
		}

		if ebr.Part_next == -1 {

			break
		}

		seek = ebr.Part_next

	}
	if flag == true {
		return &ebr
	}
	return nil

}

func ExploreDisks() {

	alfabeto := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for _, letra := range alfabeto {
		path := "MIA/P1/" + letra + ".dsk"

		file, err := os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {

			continue
		}
		defer file.Close()
		CheckMountedPartitions(file, letra)

	}

}

func CheckMountedPartitions(file *os.File, diskLetter string) {
	var mbr Structs.MBR
	reader := bufio.NewReader(file)
	err := binary.Read(reader, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}
	var seek int64
	seek = 0
	correlativo := 1

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_start != -1 {

			if mbr.Partitions[i].Part_type == 'E' {
				seek = mbr.Partitions[i].Part_start
			}

			if mbr.Partitions[i].Part_type == 'P' {
				comprobar := mbr.Partitions[i].Part_id

				x := 0
				for i := range comprobar {
					if comprobar[i] == 0 {
						x++
					}
				}

				if x != 4 {
					for j := range MountedPartitions {
						if MountedPartitions[j].Id == "" {
							MountedPartitions[j].Id = string(mbr.Partitions[i].Part_id[:])
							MountedPartitions[j].PartitionName = mbr.Partitions[i].Part_name
							MountedPartitions[j].DiskName = diskLetter
							break
						}
					}
				}
			}

			correlativo++
		}
	}

	if seek != 0 {
		var ebr Structs.EBR

		file.Seek(seek, 0)

		for {

			_, err := file.Seek(seek, 0)
			if err != nil {
				fmt.Println("Error al establecer la posición de escritura:", err)
				os.Exit(1)
			}

			reader := bufio.NewReader(file)
			err = binary.Read(reader, binary.BigEndian, &ebr)
			if ebr.Part_start == 0 {
				return
			}
			if err != nil {
				fmt.Println("Error al leer el EBR:", err)
				os.Exit(1)
			}

			if ebr.Part_mount == 1 {
				for j := range MountedPartitions {
					if MountedPartitions[j].Id == "" {
						MountedPartitions[j].PartitionName = ebr.Part_name
						MountedPartitions[j].DiskName = diskLetter
						MountedPartitions[j].Id = diskLetter + strconv.Itoa(correlativo) + "45"
						break
					}
				}
			}

			if ebr.Part_next == -1 {

				break
			}
			correlativo++

			seek = ebr.Part_next

		}
	}

}

func ReadBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func Compare(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

func findfreejournaling(journal_start int64, path string, superbloque Structs.SuperBloque) int64 {
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("MKUSER", "No se ha encontrado el disco.")
		return -1
	}
	file.Seek(journal_start, 0)

	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		var journaling Structs.Journaling
		data := ReadBytes(file, int(unsafe.Sizeof(Structs.Journaling{})))
		buffer := bytes.NewBuffer(data)
		err_ := binary.Read(buffer, binary.BigEndian, &journaling)
		if err_ != nil {
			fmt.Println("MKUSER", "Error al leer el archivo")
			return -1
		}
		if journaling.Active == 0 {
			file.Close()
			return journal_start + int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i)
		}
	}
	return -1

}
