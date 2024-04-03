package Commands

import (
	"Proyecto_1/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"
)

type Mkfs struct {
	Id         string
	Type       string
	Fs         string
	Parameters []string
}

func NewMkfs(parameters []string) *Mkfs {
	mkfs := &Mkfs{
		Id:         "",
		Type:       "",
		Fs:         "2fs",
		Parameters: parameters,
	}
	mkfs.readParameters()
	return mkfs
}

func (mkfs *Mkfs) readParameters() {
	for _, parametro := range mkfs.Parameters {
		mkfs.identifyParameters(parametro)

	}
	mkfs.MakeSystem()

}

func (mkfs *Mkfs) identifyParameters(parameter string) {
	parameter_identifier := strings.Split(parameter, "=")
	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "id" {
		mkfs.Id = strings.ToUpper(Stringmake(parameter_identifier[1]))
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "type" {
		mkfs.Type = strings.ToUpper(Stringmake(parameter_identifier[1]))
	}

	if strings.ToLower(strings.TrimSpace(parameter_identifier[0])) == "fs" {
		mkfs.Fs = strings.ToLower(Stringmake(parameter_identifier[1]))

	}
}

func (mkfs *Mkfs) MakeSystem() {
	p := ""
	partition := getMount(mkfs.Id, &p)
	var n float64
	x := 2
	if mkfs.Fs == "2fs" {

		n = math.Floor(float64(partition.Part_s-int64(unsafe.Sizeof(Structs.SuperBloque{}))) /
			float64(4+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.BloquesCarpetas{})))
	}

	if mkfs.Fs == "3fs" {
		x = 3
		n = math.Floor(float64(partition.Part_s-int64(unsafe.Sizeof(Structs.SuperBloque{}))) /
			float64(4+unsafe.Sizeof(Structs.Inodos{})+(unsafe.Sizeof(Structs.Journaling{})+3*unsafe.Sizeof(Structs.BloquesCarpetas{}))))
	}

	//Creo el superbloque, usando n como parametro, y usando las especificaciones del enunciado
	spr := Structs.NewSuperBloque()
	spr.S_magic = 0xEF53
	spr.S_inode_s = int64(unsafe.Sizeof(Structs.Inodos{}))
	spr.S_block_s = int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))
	spr.S_inodes_count = int64(n)
	spr.S_free_inodes_count = int64(n) - 2
	spr.S_blocks_count = int64(3 * n)
	spr.S_free_blocks_count = int64(3*n) - 2
	fecha := time.Now().String()
	copy(spr.S_mtime[:], fecha)
	spr.S_mnt_count = spr.S_mnt_count + 1
	spr.S_filesystem_type = int64(x)

	mkfs.ext2_ext3(spr, partition, int64(n), p)

}

func (mkfs *Mkfs) ext2_ext3(spr Structs.SuperBloque, partition *Structs.Partition, n int64, path string) {
	var journal_start int64
	journal_start = 0
	var journaling_size int64
	journaling_size = 0

	for i := 0; i < int(n); i++ {
		journaling_size += int64(unsafe.Sizeof(Structs.Journaling{}))
	}

	if partition.Part_type == 'L' {
		partition.Part_start = partition.Part_start + int64(unsafe.Sizeof(Structs.EBR{}))
	}

	if mkfs.Fs == "3fs" {
		journal_start = partition.Part_start + int64(unsafe.Sizeof(Structs.SuperBloque{}))
	}
	if journal_start == 0 {

		spr.S_bm_inode_start = partition.Part_start + int64(unsafe.Sizeof(Structs.SuperBloque{}))
	} else {
		spr.S_bm_inode_start = journal_start + journaling_size

	}

	spr.S_bm_block_start = spr.S_bm_inode_start + n
	spr.S_inode_start = spr.S_bm_block_start + (3 * n)
	spr.S_block_start = spr.S_inode_start + (n * int64(unsafe.Sizeof(Structs.Inodos{})))

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("No se encontro el disco")
		return
	}

	//Escribir el superbloque en el part_start de la partición

	file.Seek(partition.Part_start, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, spr)
	WriteBytes(file, binario2.Bytes())

	//zero := '0'
	file.Seek(spr.S_bm_inode_start, 0)
	for i := 0; i < int(n); i++ {
		WriteBytes(file, []byte{byte('0')})
	}

	file.Seek(spr.S_bm_block_start, 0)

	for i := 0; i < 3*int(n); i++ {
		WriteBytes(file, []byte{byte('0')})
	}
	inode := Structs.NewInodos()
	//Aqui crea los bloques de apuntadores del inodo y todos los atributos de este inodo
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_s = -1
	for i := 0; i < len(inode.I_block); i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1
	//Aqui escribe el inodo en el archivo binario
	file.Seek(spr.S_inode_start, 0)

	for i := 0; i < int(n); i++ {
		var binarioInodos bytes.Buffer
		binary.Write(&binarioInodos, binary.BigEndian, inode)
		WriteBytes(file, binarioInodos.Bytes())
	}
	if journal_start != 0 {

		file.Seek(journal_start, 0)
		var journaling Structs.Journaling
		for i := 0; i < int(n); i++ {
			var binariojournaling bytes.Buffer
			binary.Write(&binariojournaling, binary.BigEndian, journaling)
			WriteBytes(file, binariojournaling.Bytes())
		}
	}

	folder := Structs.NewBloquesCarpetas()

	//Escribir el bitmap de bloques, creo

	file.Seek(spr.S_block_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioFolder bytes.Buffer
		binary.Write(&binarioFolder, binary.BigEndian, folder)
		WriteBytes(file, binarioFolder.Bytes())
	}
	file.Close()

	recuperado := Structs.NewSuperBloque()

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {

		return
	}

	file.Seek(partition.Part_start, 0)
	data := ReadBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &recuperado)
	if err_ != nil {

		return
	}
	file.Close()

	//inicia un nuevo inodo, como en la raíz se tiene que crear un archivo llamado users.txt
	//entonces empieza a llenar el inodo
	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_s = 0
	fecha := time.Now().String()
	copy(inode.I_atime[:], fecha)
	copy(inode.I_ctime[:], fecha)
	copy(inode.I_mtime[:], fecha)
	inode.I_type = 0
	inode.I_perm = 664
	inode.I_block[0] = 0

	fb := Structs.NewBloquesCarpetas()
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = 0
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = 0
	copy(fb.B_content[2].B_name[:], "users.txt")
	fb.B_content[2].B_inodo = 1

	dataArchivo := "1,G,root\n1,U,root,root,123\n"

	var fileb Structs.BloquesArchivos
	copy(fileb.B_content[:], dataArchivo)

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_RDWR, os.ModeAppend)

	if err != nil {
		return
	}

	//PONE UN CARACTER 1 PORQUE SE CREO USER.TXT
	file.Seek(spr.S_bm_inode_start, 0)
	WriteBytes(file, []byte{byte('1')})
	WriteBytes(file, []byte{byte('1')})

	file.Seek(spr.S_bm_block_start, 0)
	WriteBytes(file, []byte{byte('1')})
	WriteBytes(file, []byte{byte('1')})

	file.Seek(spr.S_inode_start, 0)
	var bin3 bytes.Buffer
	binary.Write(&bin3, binary.BigEndian, inode)
	WriteBytes(file, bin3.Bytes())

	inodetmp := Structs.NewInodos()
	inodetmp.I_uid = 1
	inodetmp.I_gid = 1
	inodetmp.I_s = int64(unsafe.Sizeof(dataArchivo) + unsafe.Sizeof(Structs.BloquesCarpetas{}))

	copy(inodetmp.I_atime[:], fecha)
	copy(inodetmp.I_ctime[:], fecha)
	copy(inodetmp.I_mtime[:], fecha)
	inodetmp.I_type = 1
	inodetmp.I_perm = 664
	inodetmp.I_block[0] = 1

	inode.I_s = inodetmp.I_s + int64(unsafe.Sizeof(Structs.BloquesCarpetas{})) + int64(unsafe.Sizeof(Structs.Inodos{}))

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var bin4 bytes.Buffer
	binary.Write(&bin4, binary.BigEndian, inodetmp)
	WriteBytes(file, bin4.Bytes())

	file.Seek(spr.S_block_start, 0)

	var bin5 bytes.Buffer
	binary.Write(&bin5, binary.BigEndian, fb)
	WriteBytes(file, bin5.Bytes())

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{})), 0)
	var bin6 bytes.Buffer
	binary.Write(&bin6, binary.BigEndian, fileb)
	WriteBytes(file, bin6.Bytes())

	if mkfs.Fs == "3fs" {

		file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), 0)
		var journaling Structs.Journaling
		copy(journaling.Ruta[:], path)
		copy(journaling.Contenido[:], dataArchivo)
		copy(journaling.Fecha[:], fecha)
		copy(journaling.Operacion[:], "mkfile")
		journaling.Active = '1'
		var binjournal bytes.Buffer
		binary.Write(&binjournal, binary.BigEndian, journaling)
		WriteBytes(file, binjournal.Bytes())

	}

	/*
		file.Seek(spr.S_bm_inode_start, 0)

		diferencia := spr.S_inodes_count - spr.S_free_inodes_count
		fmt.Println(diferencia)
		file.Seek(spr.S_bm_inode_start+diferencia, 0)
		WriteBytes(file, []byte{byte('1')})


	*/
	file.Close()
	cadena := "MKFS de la particion " + string(partition.Part_name[:]) + " realizado correctamente"
	fmt.Println(cadena)

}

/*
file2, err := os.OpenFile(path, os.O_RDWR, 0644)
	file2.Seek(spr.S_bm_block_start, 0)
	for i := 0; i < int(n); i++ {
		buffer := make([]byte, 1) // Lee un byte en cada iteración
		_, err := file2.Read(buffer)
		if err != nil {
			fmt.Println("Error al leer los bytes:", err)
			return
		}

		// Imprimir el byte leído
		fmt.Printf("Byte %d: %v\n", i, buffer)
		if buffer[0] == 1 {
			fmt.Println("Aquí hay un uno")
		}
	}


*/
