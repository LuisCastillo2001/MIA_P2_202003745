package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Mkusr struct {
	User string
	Pass string
	Grp  string
}

func NewMkusr(parameters []string) *Mkusr {
	mkusr := &Mkusr{
		User: "",
		Pass: "",
		Grp:  "",
	}
	mkusr.readParameters(parameters)
	return mkusr
}

func (mkusr *Mkusr) readParameters(parameters []string) {
	for _, param := range parameters {
		mkusr.identifyParameter(param)
	}
	mkusr.Makeuser()
}

func (mkusr *Mkusr) identifyParameter(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "user" {
		mkusr.User = Stringmake(parameterIdentifier[1])
	} else if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "pass" {
		mkusr.Pass = Stringmake(parameterIdentifier[1])
	} else if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "grp" {
		mkusr.Grp = Stringmake(parameterIdentifier[1])
	}
}

func (mkusr *Mkusr) Makeuser() {
	if !Compare(Logged.User, "root") {
		Concatenar("MKUSER" + "Solo el usuario \"root\" puede acceder a estos comandos.")
		return
	}

	var path string
	partition := getMount(Logged.Id, &path)
	if partition.Part_start == -1 {
		Concatenar("MKUSER" + "No se encontró la partición montada con el id: " + Logged.Id)
		return
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Concatenar("MKUSER" + "No se ha encontrado el disco.")
		return
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Concatenar("MKUSER" + "Error al leer el archivo")
		return
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Concatenar("MKUSER" + "Error al leer el archivo")
		return
	}

	var fb Structs.BloquesArchivos
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(bloque-1), 0)

		data = ReadBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			Concatenar("MKUSER" + "Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	existe := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'G' || linea[2] == 'g') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[2] == mkusr.Grp {
				existe = true
				break
			}
		}
	}
	if !existe {
		Concatenar("MKUSER" + "No se encontró el grupo \"" + mkusr.Grp + "\".")
		return
	}

	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'U' || linea[2] == 'u' {
			c++
			in := strings.Split(linea, ",")
			if in[3] == mkusr.User {
				if linea[0] != '0' {
					Concatenar("MKUSER" + "EL nombre " + mkusr.User + ", ya está en uso.")
					return
				}
			}
		}
	}
	txt += strconv.Itoa(c+1) + ",U," + mkusr.Grp + "," + mkusr.User + "," + mkusr.Pass + "\n"
	tam := len(txt)
	var contenido []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			contenido = append(contenido, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			contenido = append(contenido, txt)
		}
	} else {
		contenido = append(contenido, txt)
	}
	if len(contenido) > 16 {
		Concatenar("MKUSER" + "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Concatenar("MKUSER" + "No se ha encontrado el disco.")
		return
	}

	for i := 0; i < len(contenido); i++ {

		var fbAux Structs.BloquesArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			WriteBytes(file, binAux.Bytes())

			//Actualizar bitmap de bloques
			diferencia := super.S_blocks_count - super.S_free_blocks_count
			file.Seek(super.S_bm_block_start+diferencia, 0)
			WriteBytes(file, []byte{byte('1')})
			super.S_free_blocks_count--
			inode.I_block[i] = diferencia

		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], contenido[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		WriteBytes(file, bin6.Bytes())

	}

	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	WriteBytes(file, inodos.Bytes())

	file.Seek(partition.Part_start, 0)
	var super2 bytes.Buffer
	binary.Write(&super2, binary.BigEndian, super)
	WriteBytes(file, super2.Bytes())
	Concatenar("MKUSER" + "Usuario " + mkusr.User + ", creado correctamente!")

	if super.S_filesystem_type == 3 {
		var journalingbytes bytes.Buffer
		var journaling Structs.Journaling
		copy(journaling.Operacion[:], "mkusr")
		copy(journaling.Ruta[:], "/users.txt")
		fecha := time.Now().String()
		copy(journaling.Fecha[:], fecha)
		copy(journaling.Contenido[:], "user: "+mkusr.User+"grp"+mkusr.Grp)
		journaling.Active = '1'
		file.Seek(findfreejournaling(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), path, super), 0)
		binary.Write(&journalingbytes, binary.BigEndian, journaling)
		WriteBytes(file, journalingbytes.Bytes())
	}

	file.Close()
}
