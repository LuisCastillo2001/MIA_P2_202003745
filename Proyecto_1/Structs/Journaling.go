package Structs

type Journaling struct {
	Operacion [8]byte
	Ruta      [40]byte
	Contenido [30]byte
	Fecha     [16]byte
	Active    byte
}

func NewJournaling() Journaling {
	var journal Journaling
	journal.Active = '0'
	for i := 0; i < 6; i++ {
		journal.Operacion[i] = 0
	}
	for i := 0; i < 16; i++ {
		journal.Ruta[i] = 0
		journal.Contenido[i] = 0
		journal.Fecha[i] = 0
	}
	return journal
}
