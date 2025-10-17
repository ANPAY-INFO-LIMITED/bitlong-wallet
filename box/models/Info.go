package models

import "gorm.io/gorm"

type Info struct {
	gorm.Model
	BlkUUID        string `json:"blk_uuid" gorm:"type:varchar(255);uniqueIndex:uuid_mid"`
	MachineID      string `json:"machine_id" gorm:"type:varchar(255);uniqueIndex:uuid_mid"`
	MachineCoding  string `json:"machine_coding" gorm:"type:varchar(255);index"`
	IdentityPubkey string `json:"identity_pubkey" gorm:"type:varchar(255);index"`
}
