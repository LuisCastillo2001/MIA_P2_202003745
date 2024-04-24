package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type UsuarioActivo struct {
	User     string
	Password string
	Id       string
	Uid      int
	Gid      int
}

var Logged UsuarioActivo

type Login struct {
	User       string
	Pass       string
	Id         string
	Parameters []string
}

func NewLogin(parameters []string) *Login {
	login := &Login{
		User:       "",
		Pass:       "",
		Id:         "",
		Parameters: parameters,
	}
	login.readParameters()
	return login

}

func (login *Login) readParameters() {
	for _, parametro := range login.Parameters {
		login.IdentifyParameters(parametro)

	}

	login.Makelogin()
}

func (login *Login) IdentifyParameters(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "user" {
		login.User = Stringmake(parameterIdentifier[1])
	}

	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "pass" {
		login.Pass = Stringmake(parameterIdentifier[1])
	}

	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "id" {
		login.Id = strings.ToUpper(Stringmake(parameterIdentifier[1]))
	}
}

func (login *Login) Makelogin() bool {
	var path string
	partition := getMount(login.Id, &path)
	if Logged.User != "" {
		Concatenar("Ya hay un usuario logueado :(")
		return false
	}
	if partition.Part_start == -1 {
		Concatenar("No se encontró el id de la partición a la cual hace referencia :(")
		return false
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {

		return false
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Concatenar("Error al leer el archivo")
		return false
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = ReadBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Concatenar("Error al leer el archivo")
		return false
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
			Concatenar("Error al leer el archivo")
			return false
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")

	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'U' || linea[2] == 'u' {
			inU := strings.Split(linea, ",")
			if Compare(inU[3], login.User) && Compare(inU[4], login.Pass) && inU[0] != "0" {
				//aqui compare que la contraseña y el usuario coincidan con el archivo
				idGrupo := "0"
				existe := false
				// aqui se encarga de que el grupo coincida y exista con el usuario
				for j := 0; j < len(vctr)-1; j++ {
					group := vctr[j]
					if (group[2] == 'G' || group[2] == 'g') && group[0] != '0' {
						inG := strings.Split(group, ",")
						// como el usuario devolvio algo, osea el grupo  y la estructura es la siguiente:
						//GID, Tipo, Grupo
						//UID, Tipo, Grupo, Usuario, Contraseña
						//0   , 1  , 2    , 3      ,  4
						//comparo si los grupos son iguales, y listo, si coincide con algún grupo
						//significa que si existe ese grupo
						if inG[2] == inU[2] {

							idGrupo = inG[0]
							existe = true
							break
						}
					}
				}
				if !existe {
					Concatenar("Login" + "No se encontró el grupo \"" + inU[2] + "\".")
					return false
				}

				Concatenar("Usuario " + login.User + " logeado con éxito")
				Logged.Id = login.Id
				Logged.User = login.User
				Logged.Password = login.Pass
				Logged.Uid, _ = strconv.Atoi(inU[0])
				Logged.Gid, _ = strconv.Atoi(idGrupo)
				return true
			}
		}
	}
	Concatenar("No se encontro al usuario")
	return false
}

func Logout() bool {
	if Logged.User == "" {
		Concatenar("No hay ningún usuario logueado")
		return false
	}

	Logged = UsuarioActivo{}
	Concatenar("Logout del usuario hecho correctamente :)")
	return true
}
