package models

import (
	"gorm.io/gorm"
)

type Key struct {
	gorm.Model
	PriKey string `json:"pri_key" gorm:"type:varchar(255);index"`
	PubKey string `json:"pub_key" gorm:"type:varchar(255);index"`
	NPub   string `json:"n_pub" gorm:"type:varchar(255);index"`
}
