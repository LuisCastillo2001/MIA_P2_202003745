package Structs

type EBR struct {
	Part_mount byte
	Part_fit   byte
	Part_start int64
	Part_s     int64
	Part_next  int64
	Part_name  [16]byte
}

func NewEBR() EBR {
	ebr := EBR{
		Part_mount: 0,
		Part_fit:   0,
		Part_start: -1,
		Part_s:     0,
		Part_next:  -1,
		Part_name:  [16]byte{},
	}
	return ebr
}
