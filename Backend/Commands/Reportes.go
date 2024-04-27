package Commands

import (
	"Proyecto_1/Structs"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

type Rep struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func NewRep(parameters []string) *Rep {
	rep := &Rep{
		Name: "",
		Path: "",
		Id:   "",
		Ruta: "",
	}
	rep.readParameters(parameters)
	return rep
}

func (rep *Rep) readParameters(parameters []string) {
	for _, param := range parameters {
		rep.identifyParameter(param)
	}
	rep.Makereports()
}

func (rep *Rep) identifyParameter(parameter string) {
	parameterIdentifier := strings.Split(parameter, "=")
	if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "name" {
		rep.Name = strings.ToLower(Stringmake(parameterIdentifier[1]))
	} else if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "path" {
		rep.Path = Stringmake(parameterIdentifier[1])
	} else if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "id" {
		rep.Id = Stringmake(parameterIdentifier[1])
	} else if strings.ToLower(strings.TrimSpace(parameterIdentifier[0])) == "ruta" {
		rep.Ruta = Stringmake(parameterIdentifier[1])
	}
}

func (rep *Rep) Makereports() {
	path := ""
	partition := getMount(rep.Id, &path)
	Concatenar("Haciendo el reporte de tipo : " + rep.Name + " de la particion o disco con el id: " + rep.Id)
	if partition.Part_start == -1 {
		Concatenar("El id de la partición no existe, vuelva a intentarlo")
		return
	}
	//Concatenar(path)
	if rep.Name == "mbr" {
		rep.Mkrepmbr(path)
	}
	if rep.Name == "disk" {
		rep.Mkrepdisk(path)
	}

	if rep.Name == "sb" {
		rep.Mkrepsuperblock(path, partition)
	}
	if rep.Name == "bm_inode" {
		rep.Mkrepbminodos(path, partition)
	}

	if rep.Name == "bm_bloque" {
		rep.Mkrepbmbloques(path, partition)
	}
	if rep.Name == "inode" {
		rep.ReportInodes(path, partition)
	}

	if rep.Name == "bloque" {
		rep.repBlocks(path, partition)
	}

	if rep.Name == "journaling" {
		rep.journalingreport(path, partition)
	}
	if rep.Name == "tree" {
		rep.makeTree(path, partition)
	}

}

func (rep *Rep) Mkrepmbr(path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	var mbr Structs.MBR

	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	if err != nil {
		return
	}
	var seek int64
	seek = 0
	graphDOT := "digraph G {a0 [shape=none label=<\n" +
		"<TABLE border=\"1\" cellborder=\"1\" cellspacing=\"4\" bgcolor=\"lightgrey\" cellpadding = \"10\">\n"

	graphDOT += " <TR>\n  <TD bgcolor=\"yellow\" colspan =\"2\" >REPORTE MBR</TD>\n  </TR>\n "
	graphDOT += " <TR>\n  <TD bgcolor=\"yellow\">mbr_tamano</TD>\n  <TD bgcolor=\"yellow\">" + strconv.FormatInt(mbr.MBR_tamano, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"yellow\">mbr_signature</TD>\n  <TD bgcolor=\"yellow\">" + strconv.FormatInt(mbr.MBR_tamano, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"yellow\">mbr_fecha</TD>\n  <TD bgcolor=\"yellow\">" + string(mbr.MBR_fecha_creacion[:]) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"yellow\">dsk_fit</TD>\n  <TD bgcolor=\"yellow\">" + string(mbr.DSK_fit) + "</TD>\n  </TR>"

	for i := range mbr.Partitions {

		if mbr.Partitions[i].Part_status == '0' {
			graphDOT += " <TR>\n  <TD bgcolor=\"turquoise\" colspan =\"2\">Particion " + strconv.Itoa(i) + "</TD>\n  \n  </TR>"
			graphDOT += " \n    <TR>\n  <TD bgcolor=\"turquoise\">Part status</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part type</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part fit</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part start</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part size</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part name</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part correlative</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part ID</TD>\n   " +
				" <TD bgcolor=\"turquoise\">" + "0" + "</TD>\n  \n  </TR>"
			continue
		}

		graphDOT += " <TR>\n  <TD bgcolor=\"turquoise\" colspan =\"2\">Particion " + strconv.Itoa(i) + "</TD>\n  \n  </TR>"
		graphDOT += " \n    <TR>\n  <TD bgcolor=\"turquoise\">Part status</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + string(mbr.Partitions[i].Part_status) + "</TD>\n  \n  </TR>"
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part type</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + string(mbr.Partitions[i].Part_type) + "</TD>\n  \n  </TR>"
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part fit</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + string(mbr.Partitions[i].Part_fit) + "</TD>\n  \n  </TR>"
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part start</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + strconv.FormatInt(mbr.Partitions[i].Part_start, 10) + "</TD>\n  \n  </TR>"
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part size</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + strconv.FormatInt(mbr.Partitions[i].Part_s, 10) + "</TD>\n  \n  </TR>"
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part name</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + string(mbr.Partitions[i].Part_name[:]) + "</TD>\n  \n  </TR>"
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part correlative</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + strconv.FormatInt(mbr.Partitions[i].Part_correlative, 10) + "</TD>\n  \n  </TR>"
		part_id := string(mbr.Partitions[i].Part_id[:])

		if mbr.Partitions[i].Part_id[0] == 0 {
			part_id = "0000"
		}
		graphDOT += "    <TR>\n  <TD bgcolor=\"turquoise\">Part ID</TD>\n   " +
			" <TD bgcolor=\"turquoise\">" + part_id + "</TD>\n  \n  </TR>"

		if mbr.Partitions[i].Part_type == 'E' {
			seek = mbr.Partitions[i].Part_start
		}
	}

	if seek != 0 {
		var ebr Structs.EBR

		for {
			_, err := file.Seek(seek, 0)
			if err != nil {
				Concatenar("Error al establecer la posición de escritura:")
				return
			}

			reader := bufio.NewReader(file)
			err = binary.Read(reader, binary.BigEndian, &ebr)
			if ebr.Part_start == 0 {
				Concatenar("Hubo un error")
				return
			}
			if err != nil {
				Concatenar("Error al leer el EBR:")
				return
			}
			graphDOT += " <TR>\n  <TD bgcolor=\"orange\" colspan =\"2\">EBR " + "</TD>\n  \n  </TR>"
			partmount := ""
			if ebr.Part_mount == '0' {
				partmount = "0"
			} else {
				partmount = "1"
			}
			graphDOT += " \n    <TR>\n  <TD bgcolor=\"orange\">EBR Part mount</TD>\n   " +
				" <TD bgcolor=\"orange\">" + partmount + "</TD>\n  \n  </TR>"

			graphDOT += " \n    <TR>\n  <TD bgcolor=\"orange\">EBR Part fit</TD>\n   " +
				" <TD bgcolor=\"orange\">" + string(ebr.Part_fit) + "</TD>\n  \n  </TR>"
			graphDOT += " \n    <TR>\n  <TD bgcolor=\"orange\">EBR Part start</TD>\n   " +
				" <TD bgcolor=\"orange\">" + strconv.FormatInt(ebr.Part_start, 10) + "</TD>\n  \n  </TR>"
			graphDOT += " \n    <TR>\n  <TD bgcolor=\"orange\">EBR Part size</TD>\n   " +
				" <TD bgcolor=\"orange\">" + strconv.FormatInt(ebr.Part_s, 10) + "</TD>\n  \n  </TR>"
			graphDOT += " \n    <TR>\n  <TD bgcolor=\"orange\">EBR Part next</TD>\n   " +
				" <TD bgcolor=\"orange\">" + strconv.FormatInt(ebr.Part_next, 10) + "</TD>\n  \n  </TR>"
			graphDOT += " \n    <TR>\n  <TD bgcolor=\"orange\">EBR Part name</TD>\n   " +
				" <TD bgcolor=\"orange\">" + string(ebr.Part_name[:]) + "</TD>\n  \n  </TR>"

			if ebr.Part_next == -1 {
				break
			}

			seek = ebr.Part_next
		}

	}
	graphDOT += "</TABLE>>];\n\n}"
	err2 := rep.generateAndSaveReport(graphDOT)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}

}

func (rep *Rep) Mkrepdisk(path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	var mbr Structs.MBR

	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &mbr)
	if err != nil {
		Concatenar("Hubo un error al abrir el disco")
		return
	}
	graphDot := "digraph D {\n    subgraph cluster_0 {\n       " +
		" bgcolor=\"#68d9e2\"\n        " +
		"node [style=\"rounded\" style=filled];" +
		" node_A [shape=record    label=\"MBR"
	size := 0.0

	for i := range mbr.Partitions {
		if mbr.Partitions[i].Part_type == 'P' {
			graphDot += "|{Primaria|"

			percentage := (float64(mbr.Partitions[i].Part_s) / float64(mbr.MBR_tamano)) * 100

			graphDot += "{" + strconv.FormatFloat(percentage, 'f', 4, 64) + "%" + "}}"
			size += percentage
		}

		if mbr.Partitions[i].Part_type == 'E' {
			graphDot += "|{Extendida|"
			percentage := (float64(mbr.Partitions[i].Part_s) / float64(mbr.MBR_tamano)) * 100

			graphDot += "{" //+ strconv.FormatFloat(percentage, 'f', 2, 64) + "%" + "}}"
			size += percentage
			var ebr Structs.EBR
			seek := mbr.Partitions[i].Part_start
			file.Seek(seek, 0)
			size_ebr := 0
			for {

				_, err := file.Seek(seek, 0)
				if err != nil {
					Concatenar("Error al establecer la posición de escritura:")
					os.Exit(1)
				}

				reader := bufio.NewReader(file)
				err = binary.Read(reader, binary.BigEndian, &ebr)

				if err != nil {
					Concatenar("Error al leer el EBR:")
					os.Exit(1)
				}

				if ebr.Part_start == 0 {
					graphDot += "LIBRE "
					graphDot += strconv.FormatFloat(percentage, 'f', 4, 64) + "%"

					break
				}

				if ebr.Part_next == -1 {
					size_ebr += int(ebr.Part_s)
					graphDot += "EBR|LOGICA "
					percentage := (float64(ebr.Part_s) / float64(mbr.MBR_tamano)) * 100

					graphDot += strconv.FormatFloat(percentage, 'f', 4, 64) + "%"

					graphDot += "|LIBRE "
					freespace := (float64(size_ebr) / float64(mbr.MBR_tamano)) * 100
					graphDot += strconv.FormatFloat(freespace, 'f', 4, 64) + "%"
					break
				}
				graphDot += "EBR|LOGICA "
				percentage := (float64(ebr.Part_s) / float64(mbr.MBR_tamano)) * 100

				graphDot += "  " + strconv.FormatFloat(percentage, 'f', 4, 64) + "%"
				graphDot += "|"
				size_ebr += int(ebr.Part_s)
				seek = ebr.Part_next

			}

			graphDot += "}}"
		}

		if i != 3 && mbr.Partitions[i].Part_start+mbr.Partitions[i].Part_s != mbr.Partitions[i+1].Part_start && mbr.Partitions[i+1].Part_status != '0' {
			libre := mbr.Partitions[i+1].Part_start - mbr.Partitions[i].Part_s - mbr.Partitions[i].Part_start
			graphDot += "|{Libre|"
			percentage := (float64(libre) / float64(mbr.MBR_tamano)) * 100
			graphDot += "{" + strconv.FormatFloat(percentage, 'f', 4, 64) + "%" + "}}"
			size += percentage
		}

		if mbr.Partitions[i].Part_status == '0' {
			graphDot += "|{Libre}"
		}

	}
	size = 100 - size
	graphDot += "|{Libre|"
	graphDot += "{" + strconv.FormatFloat(size, 'f', 2, 64) + "%" + "}}"
	graphDot += "\"];\n    }\n   \n}"
	err2 := rep.generateAndSaveReport(graphDot)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}
}

func (rep *Rep) Mkrepsuperblock(path string, partition *Structs.Partition) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	if err != nil {
		Concatenar("Hubo un error al realizar el reporte")
		return
	}
	graphDOT := "digraph G {a0 [shape=none label=<\n" +
		"<TABLE border=\"1\" cellborder=\"1\" cellspacing=\"4\" bgcolor=\"lightgrey\" cellpadding = \"10\">\n"

	graphDOT += " <TR>\n  <TD bgcolor=\"orange\" colspan =\"2\" >REPORTE SUPERBLOQUE</TD>\n  </TR>\n "
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">file_system</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_filesystem_type, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">inodes_count</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_inodes_count, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">blocks_count</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_blocks_count, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">free_blocks_count</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_free_blocks_count, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">free_inodes_count</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_free_inodes_count, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">mtime</TD>\n  <TD bgcolor=\"orange\">" + string(superbloque.S_mtime[:]) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">mnt_count</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_mnt_count, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">magic</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_magic, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">inode_s</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_inode_s, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">block_s</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_block_s, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">first_ino</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_firts_ino, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">first_blo</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_first_blo, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">bm_inode_start</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_bm_inode_start, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">bm_block_start</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_bm_block_start, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">inode_start</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_inode_start, 10) + "</TD>\n  </TR>"
	graphDOT += " <TR>\n  <TD bgcolor=\"orange\">block_start</TD>\n  <TD bgcolor=\"orange\">" + strconv.FormatInt(superbloque.S_block_start, 10) + "</TD>\n  </TR>"

	graphDOT += "</TABLE>>];\n\n}"

	err2 := rep.generateAndSaveReport(graphDOT)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}

}

func (rep *Rep) Mkrepbminodos(path string, partition *Structs.Partition) {
	graphdot := "digraph G {\n  fontname=\"Helvetica,Arial,sans-serif\"\n  node [fontname=\"Helvetica,Arial,sans-serif\"]\n " +
		" edge [fontname=\"Helvetica,Arial,sans-serif\"]\n  a0 [shape=none label=<\n <TABLE border=\"0\" cellspacing=\"6\" cellpadding=\"6\" " +
		"style=\"rounded\" bgcolor=\"white\" >"
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	file.Seek(superbloque.S_bm_inode_start, 0)
	if err != nil {
		Concatenar("Hubo un error")
		return
	}
	for i := 0; i < int(200); i++ {
		var b byte
		if err := binary.Read(file, binary.BigEndian, &b); err != nil {
			if err == io.EOF {
				break
			}
			Concatenar("Error al leer byte:")
			return
		}

		if i%20 == 0 {

			if i != 0 {
				graphdot += "</TR>"
			}
			graphdot += "<TR>"

		}

		if b == '1' {
			graphdot += "<TD bgcolor=\"green\">1</TD>"
		} else {
			graphdot += "<TD bgcolor=\"red\">0</TD>"
		}
	}

	// Cerrar la última fila si es necesario

	graphdot += "</TR>"

	graphdot += "</TABLE>>];\n\n}"
	err2 := rep.generateAndSaveReport(graphdot)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}
}

func (rep *Rep) Mkrepbmbloques(path string, partition *Structs.Partition) {
	graphdot := "digraph G {\n  fontname=\"Helvetica,Arial,sans-serif\"\n  node [fontname=\"Helvetica,Arial,sans-serif\"]\n " +
		" edge [fontname=\"Helvetica,Arial,sans-serif\"]\n  a0 [shape=none label=<\n <TABLE border=\"0\" cellspacing=\"6\" cellpadding=\"6\" " +
		"style=\"rounded\" bgcolor=\"white\" >"
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	file.Seek(superbloque.S_bm_block_start, 0)
	if err != nil {
		Concatenar("Hubo un error")
		return
	}
	for i := 0; i < int(200); i++ {
		var b byte
		if err := binary.Read(file, binary.BigEndian, &b); err != nil {
			if err == io.EOF {
				break
			}
			Concatenar("Error al leer byte:")
			return
		}

		if i%20 == 0 {

			if i != 0 {
				graphdot += "</TR>"
			}
			graphdot += "<TR>"
		}

		if b == '1' {
			graphdot += "<TD bgcolor=\"green\">1</TD>"
		} else {
			graphdot += "<TD bgcolor=\"red\">0</TD>"
		}
	}

	// Cerrar la última fila si es necesario

	graphdot += "</TR>"

	graphdot += "</TABLE>>];\n\n}"
	err2 := rep.generateAndSaveReport(graphdot)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}
}

func (rep *Rep) ReportInodes(path string, partition *Structs.Partition) {
	graphdot := "digraph TablasConectadas {\n    rankdir = \"LR\"\n    node [shape=plaintext];"
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	if err != nil {
		Concatenar("Hubo un error")
		return
	}
	var arrinodos []string
	for i := 0; i < int(superbloque.S_inodes_count); i++ {

		file.Seek(superbloque.S_inode_start+int64(i)*int64(unsafe.Sizeof(Structs.Inodos{})), 0)
		var inodo Structs.Inodos
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &inodo)

		if err != nil {
			Concatenar("Hubo un error")
			return
		}
		if inodo.I_type == -1 {
			continue
		}
		graphdot += "    inodo" + strconv.Itoa(i) + " [label = <<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">"
		graphdot += "<TR><TD COLSPAN=\"2\">Inodo " + strconv.Itoa(i) + "</TD></TR>"
		graphdot += "<TR><TD>Size</TD><TD>" + strconv.FormatInt(inodo.I_s, 10) + "</TD></TR>"
		graphdot += "<TR><TD>UID</TD><TD>" + strconv.FormatInt(inodo.I_uid, 10) + "</TD></TR>"
		graphdot += "<TR><TD>GID</TD><TD>" + strconv.FormatInt(inodo.I_gid, 10) + "</TD></TR>"
		graphdot += "<TR><TD>Atime</TD><TD>" + string(inodo.I_atime[:]) + "</TD></TR>"
		graphdot += "<TR><TD>Ctime</TD><TD>" + string(inodo.I_ctime[:]) + "</TD></TR>"
		graphdot += "<TR><TD>Mtime</TD><TD>" + string(inodo.I_mtime[:]) + "</TD></TR>"
		graphdot += "<TR><TD>Perm</TD><TD>" + strconv.FormatInt(inodo.I_perm, 10) + "</TD></TR>"
		arrinodos = append(arrinodos, "inodo"+strconv.Itoa(i))
		for j := 0; j < 16; j++ {
			if inodo.I_block[j] != -1 {
				graphdot += "<TR><TD>Block" + strconv.Itoa(j) + "</TD><TD>" + strconv.FormatInt(inodo.I_block[j], 10) + "</TD></TR>"
			} else {
				graphdot += "<TR><TD>Block" + strconv.Itoa(j) + "</TD><TD>" + strconv.FormatInt(inodo.I_block[j], 10) + "</TD></TR>"
			}
		}
		graphdot += "</TABLE>>];\n"
	}
	for i := 0; i < len(arrinodos); i++ {
		if i == 0 {
			graphdot += "    " + arrinodos[i]
		} else {
			graphdot += " -> " + arrinodos[i]
		}
	}

	graphdot += "}"
	err2 := rep.generateAndSaveReport(graphdot)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}

}

func (rep *Rep) repBlocks(path string, partition *Structs.Partition) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	graphDOT := "digraph TablasConectadas {\n    rankdir = \"LR\"\n    node [shape=plaintext];"
	if err != nil {
		Concatenar("Hubo un error al abrir el archivo")
		return
	}
	var arribloques []string
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)

	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		txt := ""
		var inode Structs.Inodos
		file.Seek(superbloque.S_inode_start+int64(i)*int64(unsafe.Sizeof(Structs.Inodos{})), 0)
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &inode)
		if err != nil {
			Concatenar("Hubo un error al leer el inodo")
			return
		}
		if inode.I_type == -1 {
			continue
		}

		for j := 0; j < 16; j++ {
			if inode.I_block[j] != -1 && inode.I_type == 0 {

				var bloque Structs.BloquesCarpetas
				file.Seek(superbloque.S_block_start+int64(inode.I_block[j])*int64(unsafe.Sizeof(Structs.BloquesCarpetas{})), 0)
				reader := bufio.NewReader(file)
				err = binary.Read(reader, binary.BigEndian, &bloque)
				if err != nil {
					Concatenar("Hubo un error al leer el bloque")
					return
				}
				graphDOT += "    bloque" + strconv.Itoa(int(inode.I_block[j])) + " [label = <<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">"
				arribloques = append(arribloques, "bloque"+strconv.Itoa(int(inode.I_block[j])))
				graphDOT += " <TR>\n <TD colspan=\"2\">" + "bloque " + strconv.Itoa(int(inode.I_block[j])) + "</TD>\n</TR>"
				graphDOT += " <TR>\n <TD>b_name</TD>\n <TD>b_inodo</TD>\n</TR>"
				for h := 0; h < 4; h++ {
					nameblock := ""
					for name := 0; name < 12; name++ {
						if bloque.B_content[h].B_name[name] == 0 {
							break
						}
						nameblock += string(bloque.B_content[h].B_name[name])

					}
					graphDOT += " <TR>\n <TD>" + nameblock + "</TD>\n <TD>" + strconv.FormatInt(bloque.B_content[h].B_inodo, 10) + "</TD>\n</TR>"
				}
				graphDOT += "</TABLE>>];\n"

			} else if inode.I_type == 1 && inode.I_block[j] != -1 {
				graphDOT += "    bloque" + strconv.Itoa(int(inode.I_block[j])) + " [label = <<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">"
				arribloques = append(arribloques, "bloque"+strconv.Itoa(int(inode.I_block[j])))
				var bloque Structs.BloquesArchivos
				file.Seek(superbloque.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(inode.I_block[j]-1), 0)
				reader := bufio.NewReader(file)
				err = binary.Read(reader, binary.BigEndian, &bloque)
				for k := 0; k < 64; k++ {
					if bloque.B_content[k] != 0 {
						txt += string(bloque.B_content[k])
					}
				}
				graphDOT += " <TR>\n <TD colspan=\"2\">" + "bloque" + strconv.Itoa(int(inode.I_block[j])) + "</TD>\n</TR>"
				graphDOT += " <TR>\n <TD>" + txt + "</TD>\n</TR>"
				graphDOT += "</TABLE>>];\n"

			}

		}

	}

	//aqui
	for j := 0; j < len(arribloques); j++ {

		if j == 0 {
			graphDOT += "    " + arribloques[j]
		} else {
			graphDOT += " -> " + arribloques[j]
		}
	}
	graphDOT += "}"
	err2 := rep.generateAndSaveReport(graphDOT)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}

}

func (rep *Rep) journalingreport(path string, partition *Structs.Partition) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	if err != nil {
		Concatenar("Hubo un error")
		return
	}

	file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), 0)

	if err != nil {
		Concatenar("Hubo un error")
		return
	}
	graphDOT := "digraph G {a0 [shape=none label=<\n" +
		"<TABLE border=\"1\" cellborder=\"1\" cellspacing=\"4\" bgcolor=\"lightgrey\" cellpadding = \"10\">\n"
	graphDOT += "<TR>"
	graphDOT += "\n  <TD bgcolor=\"orange\" colspan =\"5\" >REPORTE JOURNALING</TD>\n  \n "
	graphDOT += "</TR>"
	graphDOT += "<TR>"
	graphDOT += "\n  <TD bgcolor=\"orange\"  >OPERACION</TD>\n  \n "
	graphDOT += " \n  <TD bgcolor=\"orange\" >RUTA</TD>\n \n "
	graphDOT += " \n  <TD bgcolor=\"orange\"  >CONTENIDO</TD>\n  \n "
	graphDOT += " \n  <TD bgcolor=\"orange\"  >FECHA</TD>\n  \n "
	graphDOT += "</TR>"

	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		var journaling Structs.Journaling
		file.Seek(partition.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{}))+int64(unsafe.Sizeof(Structs.Journaling{}))*int64(i), 0)
		reader := bufio.NewReader(file)
		err = binary.Read(reader, binary.BigEndian, &journaling)
		if err != nil {
			Concatenar("Hubo un error")
			return
		}
		if journaling.Active == '1' {

			operacionStr := string(bytes.TrimRight(journaling.Operacion[:], "\x00"))

			rutaStr := string(bytes.TrimRight(journaling.Ruta[:], "\x00"))

			contenidoStr := string(bytes.TrimRight(journaling.Contenido[:], "\x00"))

			fechaStr := string(bytes.TrimRight(journaling.Fecha[:], "\x00"))
			graphDOT += "<TR>"
			graphDOT += "   <TD bgcolor=\"orange\" >" + operacionStr + "</TD>\n   "
			graphDOT += "  <TD bgcolor=\"orange\" >" + rutaStr + "</TD>\n   "
			graphDOT += "  <TD bgcolor=\"orange\" >" + contenidoStr + "</TD>\n   "
			graphDOT += "   <TD bgcolor=\"orange\" >" + fechaStr + "</TD>\n   "
			graphDOT += "</TR>"
		}

	}
	graphDOT += "</TABLE>>];\n\n}"
	err2 := rep.generateAndSaveReport(graphDOT)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}

}

func (rep *Rep) makeTree(path string, partition *Structs.Partition) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	file.Seek(partition.Part_start, 0)
	var superbloque Structs.SuperBloque
	reader := bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &superbloque)
	if err != nil {
		Concatenar("Hubo un error al leer el inodo")
		return
	}
	raiz := Structs.NewInodos()
	file.Seek(int64(superbloque.S_inode_start), 0)
	reader = bufio.NewReader(file)
	err = binary.Read(reader, binary.BigEndian, &raiz)
	graphdot := "digraph H {\n"
	graphdot += "node [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"1\"];\n"
	graphdot += "node [shape=plaintext];\n"
	graphdot += "graph [bb=\"0,0,352,154\"];\n"
	graphdot += "rankdir=LR;\n"
	graphdot += treerecursive(raiz, superbloque, file, 0)
	graphdot += "}"

	err2 := rep.generateAndSaveReport(graphdot)
	if err2 != nil {
		Concatenar("Hubo un error")
		return
	}

}

func treerecursive(inodo Structs.Inodos, sb Structs.SuperBloque, archivo *os.File, numeroInodo int) string {

	grahpdot := "Inodo" + strconv.Itoa(numeroInodo) + "[label = <\n"
	grahpdot += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	grahpdot += "<tr><td bgcolor=\"lightgrey\">Inodo" + strconv.Itoa(numeroInodo) + "</td></tr>\n"
	grahpdot += "<tr><td>i_uid</td><td>" + strconv.FormatInt(inodo.I_uid, 10) + "</td></tr>\n"
	grahpdot += "<tr><td>i_gid</td><td>" + strconv.FormatInt(inodo.I_gid, 10) + "</td></tr>\n"
	grahpdot += "<tr><td>i_size</td><td>" + strconv.Itoa(int(inodo.I_s)) + "</td></tr>\n"
	grahpdot += "<tr><td>i_atime</td><td>" + string(inodo.I_atime[:]) + "</td></tr>\n"
	grahpdot += "<tr><td>i_ctime</td><td>" + string(inodo.I_ctime[:]) + "</td></tr>\n"
	grahpdot += "<tr><td>i_mtime</td><td>" + string(inodo.I_mtime[:]) + "</td></tr>\n"
	luxk := ""
	if inodo.I_type == '0' {
		luxk = "0"
	} else {
		luxk = "1"
	}

	grahpdot += "<tr><td>i_type</td><td>" + luxk + "</td></tr>\n"
	grahpdot += "<tr><td>i_perm</td><td>" + strconv.Itoa(int(inodo.I_perm)) + "</td></tr>\n"
	Contador := 0
	for _, i := range inodo.I_block {
		grahpdot += "<tr><td>i_block" + strconv.Itoa(Contador+1) + "</td><td port='" + strconv.Itoa(Contador+1) + "'>" + strconv.Itoa(int(i)) + "</td></tr>\n"
		Contador++
	}
	grahpdot += "</table>>];\n"
	Contador = 0
	for _, i := range inodo.I_block {
		if i != -1 {
			//Leer el bloque
			grahpdot += "Inodo" + strconv.Itoa(numeroInodo) + ":" + strconv.Itoa(Contador+1) + " -> Bloque" + strconv.Itoa(int(i)) + ":0;\n"
			grahpdot += "Bloque" + strconv.Itoa(int(i)) + "[label = <\n"
			grahpdot += "<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
			DesplazamientoBloque := int(sb.S_block_start) + (int(i) * int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
			carpeta := Structs.NewBloquesCarpetas()
			archivo.Seek(int64(DesplazamientoBloque), 0)
			reader := bufio.NewReader(archivo)
			binary.Read(reader, binary.BigEndian, &carpeta)

			if inodo.I_type == 0 {

				grahpdot += "<tr><td colspan=\"2\" port='0'>Bloque" + strconv.Itoa(int(i)) + "</td></tr>\n"
				Contador2 := 0
				for _, j := range carpeta.B_content {

					nam := strings.TrimRight(string(j.B_name[:]), string(rune(0)))

					if Contador2 == 0 {
						nam = "."
					}
					if Contador2 == 1 {
						nam = ".."
					}
					if j.B_inodo == -1 {
						nam = ""
					}

					grahpdot += "<tr><td>" + nam + "</td><td port='" + strconv.Itoa(Contador2+1) + "'>" + strconv.Itoa(int(j.B_inodo)) + "</td></tr>\n"
					Contador2++
				}
				grahpdot += "</table>>];\n"
				Contador2 = 0

				for j := 0; j < 4; j++ {
					if carpeta.B_content[j].B_inodo != -1 {
						if j != 0 && j != 1 {
							grahpdot += "Bloque" + strconv.Itoa(int(i)) + ":" + strconv.Itoa(Contador2+1) + " -> Inodo" + strconv.Itoa(int(carpeta.B_content[j].B_inodo)) + ":0;\n"

							DesplazamientoInodo := int(sb.S_inode_start) + (int(carpeta.B_content[j].B_inodo) * binary.Size(Structs.Inodos{}))
							inodoSiguiente := Structs.NewInodos()

							archivo.Seek(int64(DesplazamientoInodo), 0)
							reader := bufio.NewReader(archivo)
							binary.Read(reader, binary.BigEndian, &inodoSiguiente)
							grahpdot += treerecursive(inodoSiguiente, sb, archivo, int(carpeta.B_content[j].B_inodo))

						}
					}
					Contador2++
				}

			} else {
				file := Structs.BloquesArchivos{}
				archivo.Seek(int64(sb.S_block_start)+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i-1), 0)
				reader := bufio.NewReader(archivo)
				binary.Read(reader, binary.BigEndian, &file)
				txt := ""
				for k := 0; k < 64; k++ {
					if file.B_content[k] != 0 {
						txt += string(file.B_content[k])
					}
				}

				grahpdot += "<tr><td colspan=\"1\" port='0'>Bloque" + strconv.Itoa(int(i)) + "</td></tr>\n"
				grahpdot += "<tr><td port='1'>" + txt + "</td></tr>\n"
				grahpdot += "</table>>];\n"
			}
		}
		Contador++
	}

	return grahpdot
}

func (rep *Rep) generateAndSaveReport(graphDOT string) error {
	// Extraer el nombre del reporte y el formato de imagen del parámetro path
	fileName := filepath.Base(rep.Path)

	// Remover la extensión del archivo
	//reportName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Construir la ruta completa para el archivo DOT
	//dotFilePath := filepath.Join(filepath.Dir(rep.Path), reportName+".dot")
	x := strings.Split(fileName, ".")

	dotFilePath := "MIA/Reportes/" + x[0] + ".dot"

	err := rep.saveGraphToFile(dotFilePath, graphDOT)

	if err != nil {
		Concatenar("Hubo un error")
		return fmt.Errorf("Error al guardar el archivo DOT: %v", err)
	}

	// Construir la ruta completa para el archivo de imagen
	//filePath := "MIA/Img-rep/" + x[0] + ".pdf"

	return nil
}

func (rep *Rep) saveGraphToFile(filePath, content string) error {

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func (rep *Rep) convertDotToImage(dotFilePath, imageFilePath string) error {
	cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", imageFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error al convertir DOT a imagen: %v\n%s", err, output)
	}

	return nil
}

func (rep *Rep) convertDotToPdf(dotFilePath, pdfFilePath string) error {
	cmd := exec.Command("dot", "-Tpdf", dotFilePath, "-o", pdfFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error al convertir DOT a PDF: %v\n%s", err, output)
	}

	return nil
}
