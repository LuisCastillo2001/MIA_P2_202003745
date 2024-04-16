package Structs

import (
	"math/rand"
	"time"
)

type MBR struct {
	MBR_tamano         int64
	MBR_fecha_creacion [16]byte
	MBR_dsk_signature  int64
	DSK_fit            byte
	Partitions         [4]Partition
}

func NewMBR(tamano int64, fit byte) MBR {
	var mbr MBR
	mbr.MBR_tamano = tamano
	copy(mbr.MBR_fecha_creacion[:], Date())
	mbr.MBR_dsk_signature = Random()
	mbr.DSK_fit = fit
	for i := range mbr.Partitions {
		mbr.Partitions[i] = NewPartition()
	}

	return mbr
}

func Random() int64 {
	rand.Seed(time.Now().UnixNano())
	aleatorio := rand.Intn(9000) + 1000
	return int64(aleatorio)
}

func Date() string {
	fechaActual := time.Now()

	formato := "2006-01-02 15:04:05"
	fechaFormateada := fechaActual.Format(formato)

	return fechaFormateada
}
