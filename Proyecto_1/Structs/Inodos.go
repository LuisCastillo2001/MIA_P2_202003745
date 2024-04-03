package Structs

import "time"

type Inodos struct {
	I_uid   int64
	I_gid   int64
	I_s     int64
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [16]int64
	I_type  int64
	I_perm  int64
}

func NewInodos() Inodos {
	var inode Inodos
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_s = 112
	for i := 0; i < 16; i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1
	copy(inode.I_atime[:], time.Now().String())
	copy(inode.I_ctime[:], time.Now().String())
	copy(inode.I_mtime[:], time.Now().String())

	return inode
}
