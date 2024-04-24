package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"unsafe"
)

func ReturnDirs(disk string, path string, existe *bool) []string {
	var dirs []string
	hola := ""
	partition := getMount(Logged.Id, &hola)
	if partition.Part_start == -1 {

		fmt.Println("MKUSER" + "No se encontró la partición montada con el id: " + Logged.Id)
		return dirs
	}

	file, err := os.Open(strings.ReplaceAll(hola, "\"", ""))
	if err != nil {
		fmt.Println("MKDIR" + "No se ha encontrado el disco.")
		return dirs
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)

	if err_ != nil {
		return dirs
	}

	if err_ != nil {
		Concatenar("Error al leer el archivo")
		return dirs
	}

	path1 := strings.Split(path, "/")
	fmt.Println(path1)
	if len(path1) == 2 && path1[1] == "" {
		for i := 0; i < 16; i++ {
			file.Seek(super.S_inode_start, 0)
			inoderaiz := Structs.NewInodos()
			data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &inoderaiz)
			if inoderaiz.I_block[i] != -1 {
				// Escribir en la raíz
				//Actualizar el inodo raíz
				var fbaux Structs.BloquesCarpetas
				//escribir bloque de carpetas
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*inoderaiz.I_block[i], 0)
				reader := bufio.NewReader(file)
				err = binary.Read(reader, binary.BigEndian, &fbaux)
				for j := 0; j < 4; j++ {
					if fbaux.B_content[j].B_inodo != -1 && j != 0 && j != 1 {
						nameblock := ""
						for name := 0; name < 12; name++ {
							if fbaux.B_content[j].B_name[name] == 0 {
								break
							}
							nameblock += string(fbaux.B_content[j].B_name[name])

						}
						*existe = true
						dirs = append(dirs, nameblock)

					}
				}

			}

		}
	} else {
		fnd1 := false

		file.Seek(super.S_inode_start, 0)
		inodeaux := Structs.NewInodos()
		data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inodeaux)

		path1 = path1[1:]

		for i := 0; i < len(path1); i++ {
			fnd1 = false

			for j := 0; j < 16; j++ {
				if inodeaux.I_block[j] == -1 {
					continue
				}
				var fb Structs.BloquesCarpetas
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inodeaux.I_block[j]), 0)
				data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &fb)

				for h := 0; h < 4; h++ {
					nameblock := ""
					for name := 0; name < 12; name++ {
						if fb.B_content[h].B_name[name] == 0 {
							break
						}
						nameblock += string(fb.B_content[h].B_name[name])

					}

					if nameblock == path1[i] {

						file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*fb.B_content[h].B_inodo, 0)
						inodeaux = Structs.NewInodos()
						data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &inodeaux)

						fnd1 = true
						*existe = true
						break
					}
				}

			}

		}

		if fnd1 == false {

			*existe = false
			return dirs
		}

		for i := 0; i < 16; i++ {
			if inodeaux.I_block[i] != -1 {
				// Escribir en la raíz
				//Actualizar el inodo raíz
				var fbaux Structs.BloquesCarpetas
				//escribir bloque de carpetas
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*inodeaux.I_block[i], 0)
				reader := bufio.NewReader(file)
				err = binary.Read(reader, binary.BigEndian, &fbaux)
				for j := 0; j < 4; j++ {
					if fbaux.B_content[j].B_inodo != -1 && j != 0 && j != 1 {
						nameblock := ""
						for name := 0; name < 12; name++ {
							if fbaux.B_content[j].B_name[name] == 0 {
								break
							}
							nameblock += string(fbaux.B_content[j].B_name[name])

						}
						dirs = append(dirs, nameblock)

					}
				}

			}

		}

	}

	return dirs
}
