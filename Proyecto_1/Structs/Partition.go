package Structs

type Partition struct {
	Part_status      byte
	Part_type        byte
	Part_fit         byte
	Part_start       int64
	Part_s           int64
	Part_name        [16]byte
	Part_correlative int64
	Part_id          [4]byte
}

func NewPartition() Partition {
	partition := &Partition{
		Part_status:      '0',
		Part_type:        ' ',
		Part_fit:         ' ',
		Part_start:       -1,
		Part_s:           0,
		Part_name:        [16]byte{},
		Part_correlative: 0,
		Part_id:          [4]byte{},
	}
	return *partition
}
