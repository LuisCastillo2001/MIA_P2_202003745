package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
)

type Mkdir struct {
	Path string
	R    int
}

func NewMkdir(parameters []string) *Mkdir {
	mkdir := &Mkdir{
		Path: "",
		R:    0,
	}
	mkdir.readParameters(parameters)
	return mkdir
}

func (mkdir *Mkdir) readParameters(parameters []string) {
	for _, parametro := range parameters {
		mkdir.identifyParameters(parametro)
	}
	if mkdir.R == 1 {
		mkdir.Createalldirs()
		return
	}
	mkdir.Makedirectory(mkdir.Path, 0)
}

func (mkdir *Mkdir) identifyParameters(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	parameterName := strings.ToLower(strings.TrimSpace(parameterIdentifier[0]))

	if parameterName == "path" {
		mkdir.Path = Stringmake(parameterIdentifier[1])
	}

	if parameterName == "r" {
		mkdir.R = 1
	}
}

func (mkdir *Mkdir) Makedirectory(pathaux string, numero int) bool {
	var path string
	partition := getMount(Logged.Id, &path)
	if partition.Part_start == -1 {
		fmt.Println("MKUSER", "No se encontró la partición montada con el id: "+Logged.Id)
		return false
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		fmt.Println("MKDIR", "No se ha encontrado el disco.")
		return false
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		return false
	}

	if err_ != nil {
		fmt.Println("Error al leer el archivo")
		return false
	}
	path1 := strings.Split(pathaux, "/")
	fnd := false
	if len(path1) == 2 {

		file.Seek(super.S_inode_start, 0)
		inoderaiz := Structs.NewInodos()
		data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inoderaiz)
		var fb Structs.BloquesCarpetas

		for i := 0; i < 16; i++ {

			if inoderaiz.I_block[i] == -1 {
				continue
			}
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inoderaiz.I_block[i]), 0)
			data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &fb)
			for j := 0; j < 4; j++ {
				if fb.B_content[j].B_inodo == -1 {
					continue
				}
				nameblock := ""
				for name := 0; name < 12; name++ {
					if fb.B_content[j].B_name[name] == 0 {
						break
					}
					nameblock += string(fb.B_content[j].B_name[name])

				}
				if nameblock == path1[1] && numero == 0 {
					fmt.Println("El nombre de la carpeta ya existe")
					return false

				} else if nameblock == path1[1] && numero == 1 {
					return false
				}
			}
		}

		for i := 0; i < 16; i++ {
			if fnd {
				break
			}

			if inoderaiz.I_block[i] == -1 {

				var fbaux Structs.BloquesCarpetas
				copy(fbaux.B_content[0].B_name[:], ".")
				copy(fbaux.B_content[1].B_name[:], "..")
				copy(fbaux.B_content[2].B_name[:], path1[1])
				diferenciabloque := super.S_blocks_count - super.S_free_blocks_count
				diferenciainodo := super.S_inodes_count - super.S_free_inodes_count

				fbaux.B_content[0].B_inodo = diferenciabloque
				fbaux.B_content[1].B_inodo = 0
				fbaux.B_content[2].B_inodo = diferenciainodo

				fbaux.B_content[3].B_inodo = -1

				//Escribir inodo

				inodeaux := Structs.NewInodos()
				inodeaux.I_uid = int64(Logged.Uid)
				inodeaux.I_gid = int64(Logged.Gid)
				inodeaux.I_s = 112
				inodeaux.I_type = 0
				inodeaux.I_perm = 664

				inoderaiz.I_block[i] = diferenciabloque

				// Escribir en la raíz
				//Actualizar el inodo raíz
				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				file.Seek(super.S_inode_start, 0)
				var binarioInodos bytes.Buffer
				binary.Write(&binarioInodos, binary.BigEndian, inoderaiz)
				WriteBytes(file, binarioInodos.Bytes())
				//escribir inodo nuevo
				file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(diferenciainodo), 0)
				var binarioInodonuevo bytes.Buffer
				binary.Write(&binarioInodonuevo, binary.BigEndian, inodeaux)
				WriteBytes(file, binarioInodonuevo.Bytes())

				//escribir bloque de carpetas
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*diferenciabloque, 0)
				var binariobloque bytes.Buffer
				binary.Write(&binariobloque, binary.BigEndian, fbaux)
				WriteBytes(file, binariobloque.Bytes())

				//actualizar el superbloque

				//actualizar el bitmap de inodos
				file.Seek(super.S_bm_block_start+diferenciabloque, 0)
				WriteBytes(file, []byte{byte('1')})

				//actualizar el bitmap de bloques
				file.Seek(super.S_bm_inode_start+diferenciainodo, 0)
				WriteBytes(file, []byte{byte('1')})
				fnd = true

				super.S_free_blocks_count--
				super.S_free_inodes_count--

				file.Seek(partition.Part_start, 0)
				var superbloquebinario bytes.Buffer
				binary.Write(&superbloquebinario, binary.BigEndian, super)
				WriteBytes(file, superbloquebinario.Bytes())

				path2 := strings.Split(pathaux, "/")
				pathtoshow := path2[len(path2)-1]
				if len(path2) > 0 {
					path2 = path2[:len(path2)-1]
				}
				stack := ""
				for v := 0; v < len(path2); v++ {
					stack += "/" + path2[v]
				}

				fmt.Println("Se creo la carpeta " + pathtoshow + " en la ubicación " + stack)

				if super.S_filesystem_type == 3 {
					var journalingbytes bytes.Buffer
					var journaling Structs.Journaling
					copy(journaling.Operacion[:], "mkdir")
					copy(journaling.Ruta[:], stack)
					fecha := time.Now().String()
					copy(journaling.Fecha[:], fecha)
					copy(journaling.Contenido[:], "mkdir")
					journaling.Active = '1'
					file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
					binary.Write(&journalingbytes, binary.BigEndian, journaling)
					WriteBytes(file, journalingbytes.Bytes())
				}

				return true

			} else {
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inoderaiz.I_block[i]), 0)
				data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &fb)
				for j := 0; j < 4; j++ {

					if fb.B_content[j].B_inodo == -1 {

						file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)

						inodeaux := Structs.NewInodos()
						inodeaux.I_uid = int64(Logged.Uid)
						inodeaux.I_gid = int64(Logged.Gid)
						inodeaux.I_s = 112
						inodeaux.I_type = 0
						inodeaux.I_perm = 664
						diferenciainodo := super.S_inodes_count - super.S_free_inodes_count
						file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*(diferenciainodo), 0)
						var binarioInodonuevo bytes.Buffer
						binary.Write(&binarioInodonuevo, binary.BigEndian, inodeaux)
						WriteBytes(file, binarioInodonuevo.Bytes())

						fb.B_content[3].B_inodo = diferenciainodo
						copy(fb.B_content[j].B_name[:], path1[1])

						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inoderaiz.I_block[i]), 0)
						var binariobloquecarpetas bytes.Buffer
						binary.Write(&binariobloquecarpetas, binary.BigEndian, fb)
						WriteBytes(file, binariobloquecarpetas.Bytes())

						file.Seek(super.S_bm_inode_start+diferenciainodo, 0)
						WriteBytes(file, []byte{byte('1')})

						super.S_free_inodes_count--

						file.Seek(partition.Part_start, 0)
						var superbloquebinario bytes.Buffer
						binary.Write(&superbloquebinario, binary.BigEndian, super)
						WriteBytes(file, superbloquebinario.Bytes())
						path2 := strings.Split(pathaux, "/")
						pathtoshow := path2[len(path2)-1]
						if len(path2) > 0 {
							path2 = path2[:len(path2)-1]
						}
						stack := ""
						for v := 0; v < len(path2); v++ {
							stack += "/" + path2[v]
						}

						fmt.Println("Se creo la carpeta " + pathtoshow + " en la ubicación " + stack)
						if super.S_filesystem_type == 3 {
							var journalingbytes bytes.Buffer
							var journaling Structs.Journaling
							copy(journaling.Operacion[:], "mkdir")
							copy(journaling.Ruta[:], stack)
							fecha := time.Now().String()
							copy(journaling.Fecha[:], fecha)
							copy(journaling.Contenido[:], "mkdir")
							journaling.Active = '1'
							file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
							binary.Write(&journalingbytes, binary.BigEndian, journaling)
							WriteBytes(file, journalingbytes.Bytes())
						}
						fnd = true
						return true

					}
				}

			}

		}

	} else {
		fnd1 := false
		var number int64
		number = 0
		file.Seek(super.S_inode_start, 0)
		inodeaux := Structs.NewInodos()
		data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inodeaux)
		namefolder := path1[len(path1)-1]
		path1 = path1[:len(path1)-1]
		path1 = path1[1:]

		//Buscar el padre
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

						number = fb.B_content[h].B_inodo

						file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*fb.B_content[h].B_inodo, 0)
						inodeaux = Structs.NewInodos()
						data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &inodeaux)

						fnd1 = true
						break
					}
				}

			}

		}

		//creacion
		if fnd1 == false {
			fmt.Println("No se encontro el directorio padre")
			return false
		}

		/*
			for h := 0; h < 15; h++ {
				fmt.Println(inodeaux.I_block[h])
			}

		*/

		//verificar si existe el nombre
		for h := 0; h < 16; h++ {
			if inodeaux.I_block[h] == -1 {
				continue
			}
			var fb Structs.BloquesCarpetas
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inodeaux.I_block[h]), 0)
			data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
			buffer = bytes.NewBuffer(data)

			err_ = binary.Read(buffer, binary.BigEndian, &fb)
			for j := 0; j < 4; j++ {
				if fb.B_content[j].B_inodo == -1 {
					continue
				}
				nameblock := ""
				for name := 0; name < 12; name++ {
					if fb.B_content[j].B_name[name] == 0 {
						break
					}
					nameblock += string(fb.B_content[j].B_name[name])

				}
				if nameblock == namefolder && numero == 0 {
					fmt.Println("El nombre de la carpeta ya existe")
					return false

				} else if numero == 1 && nameblock == namefolder {
					return false
				}
			}

		}

		file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
		/*
			for h := 0; h < 15; h++ {
				fmt.Println(inodeaux.I_block[h])
			}


		*/
		fnd2 := false
		var fb Structs.BloquesCarpetas
		fnd = false
		for i := 0; i < 16; i++ {
			if fnd2 {
				break
			}
			if inodeaux.I_block[i] == -1 {

				var fbaux Structs.BloquesCarpetas
				copy(fbaux.B_content[0].B_name[:], ".")
				copy(fbaux.B_content[1].B_name[:], "..")
				copy(fbaux.B_content[2].B_name[:], namefolder)

				fbaux.B_content[0].B_inodo = 0
				fbaux.B_content[1].B_inodo = number

				fbaux.B_content[3].B_inodo = -1

				//Escribir inodo

				inodeaux2 := Structs.NewInodos()
				inodeaux2.I_uid = int64(Logged.Uid)
				inodeaux2.I_gid = int64(Logged.Gid)
				inodeaux2.I_s = 112
				inodeaux2.I_type = 0
				inodeaux2.I_perm = 664

				diferenciabloque := super.S_blocks_count - super.S_free_blocks_count
				inodeaux.I_block[i] = diferenciabloque
				fbaux.B_content[0].B_inodo = diferenciabloque

				diferenciainodo := super.S_inodes_count - super.S_free_inodes_count
				fbaux.B_content[2].B_inodo = diferenciainodo
				// Escribir el padre
				//Actualizar el inodo padre
				file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
				file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*number, 0)
				var binarioInodos bytes.Buffer
				binary.Write(&binarioInodos, binary.BigEndian, inodeaux)
				WriteBytes(file, binarioInodos.Bytes())
				//escribir inodo nuevo
				file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(diferenciainodo), 0)
				var binarioInodonuevo bytes.Buffer
				binary.Write(&binarioInodonuevo, binary.BigEndian, inodeaux2)
				WriteBytes(file, binarioInodonuevo.Bytes())

				//escribir bloque de carpetas
				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*diferenciabloque, 0)
				var binariobloque bytes.Buffer
				binary.Write(&binariobloque, binary.BigEndian, fbaux)
				WriteBytes(file, binariobloque.Bytes())

				//actualizar el superbloque

				//actualizar el bitmap de inodos
				file.Seek(super.S_bm_block_start+diferenciabloque, 0)
				WriteBytes(file, []byte{byte('1')})

				//actualizar el bitmap de bloques
				file.Seek(super.S_bm_inode_start+diferenciainodo, 0)
				WriteBytes(file, []byte{byte('1')})

				super.S_free_blocks_count--
				super.S_free_inodes_count--

				file.Seek(partition.Part_start, 0)
				var superbloquebinario bytes.Buffer
				binary.Write(&superbloquebinario, binary.BigEndian, super)
				WriteBytes(file, superbloquebinario.Bytes())

				path2 := strings.Split(pathaux, "/")
				pathtoshow := path2[len(path2)-1]
				if len(path2) > 0 {
					path2 = path2[:len(path2)-1]
				}
				stack := ""
				for v := 0; v < len(path2); v++ {
					stack += "/" + path2[v]
				}

				fmt.Println("Se creo la carpeta " + pathtoshow + " en la ubicación " + stack)
				if super.S_filesystem_type == 3 {
					var journalingbytes bytes.Buffer
					var journaling Structs.Journaling
					copy(journaling.Operacion[:], "mkdir")
					copy(journaling.Ruta[:], stack)
					fecha := time.Now().String()
					copy(journaling.Fecha[:], fecha)
					copy(journaling.Contenido[:], "mkdir")
					journaling.Active = '1'
					file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
					binary.Write(&journalingbytes, binary.BigEndian, journaling)
					WriteBytes(file, journalingbytes.Bytes())
				}
				return true
				break

			} else {

				file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inodeaux.I_block[i]), 0)
				data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
				buffer = bytes.NewBuffer(data)
				err_ = binary.Read(buffer, binary.BigEndian, &fb)
				for j := 0; j < 4; j++ {
					if fb.B_content[j].B_inodo == -1 {
						file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)

						inodeaux2 := Structs.NewInodos()
						inodeaux2.I_uid = int64(Logged.Uid)
						inodeaux2.I_gid = int64(Logged.Gid)
						inodeaux2.I_s = 112
						inodeaux2.I_type = 0
						inodeaux2.I_perm = 664
						diferenciainodo := super.S_inodes_count - super.S_free_inodes_count

						file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(diferenciainodo), 0)
						var binarioInodonuevo bytes.Buffer
						binary.Write(&binarioInodonuevo, binary.BigEndian, inodeaux2)
						WriteBytes(file, binarioInodonuevo.Bytes())

						copy(fb.B_content[j].B_name[:], namefolder)
						fb.B_content[j].B_inodo = diferenciainodo

						file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*int64(inodeaux.I_block[i]), 0)
						var binariobloquecarpetas bytes.Buffer
						binary.Write(&binariobloquecarpetas, binary.BigEndian, fb)
						WriteBytes(file, binariobloquecarpetas.Bytes())

						file.Seek(super.S_bm_inode_start+diferenciainodo, 0)
						WriteBytes(file, []byte{byte('1')})

						super.S_free_inodes_count--

						file.Seek(partition.Part_start, 0)
						var superbloquebinario bytes.Buffer
						binary.Write(&superbloquebinario, binary.BigEndian, super)
						WriteBytes(file, superbloquebinario.Bytes())
						path2 := strings.Split(pathaux, "/")
						pathtoshow := path2[len(path2)-1]
						if len(path2) > 0 {
							path2 = path2[:len(path2)-1]
						}
						stack := ""
						for v := 0; v < len(path2); v++ {
							stack += "/" + path2[v]
						}

						fmt.Println("Se creo la carpeta " + pathtoshow + " en la ubicación " + stack)
						fnd2 = true
						if super.S_filesystem_type == 3 {
							var journalingbytes bytes.Buffer
							var journaling Structs.Journaling
							copy(journaling.Operacion[:], "mkdir")
							copy(journaling.Ruta[:], stack)
							fecha := time.Now().String()
							copy(journaling.Fecha[:], fecha)
							copy(journaling.Contenido[:], "mkdir")
							journaling.Active = '1'
							file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
							binary.Write(&journalingbytes, binary.BigEndian, journaling)
							WriteBytes(file, journalingbytes.Bytes())
						}
						return true

					}
				}
			}

		}

	}
	fmt.Println("Hubo un error y algo se creo incorrectamente")
	return false
}

func (mkdir *Mkdir) Createalldirs() {

	path := strings.Split(mkdir.Path, "/")
	path = path[1:]

	stack := ""
	pathaux := path

	for i := 0; i < len(path); i++ {
		stack += "/" + path[i]
		create := mkdir.Makedirectory(stack, 1)
		pathaux = pathaux[1:]

		if create == false && len(pathaux) == 0 {
			fmt.Println("No se pudo crear el directorio")
			return
		}

	}
	fmt.Println("Directorios " + mkdir.Path + " creados con éxito")
	//Crear los directorios si estos no existen

}

//file, err = os.Open(strings.ReplaceAll(path, "\"", ""))

/*
	if inode.I_type == 1 {
		fmt.Println("Esto es un bloque de archivos")
		fmt.Println(inode.I_block[i] - 1)
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*(inode.I_block[i]-1), 0)
		var fb Structs.BloquesArchivos
		data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)
		fmt.Println()
	}
*/
