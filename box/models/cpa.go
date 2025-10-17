package models

import "gorm.io/gorm"

type Cpa struct {
	gorm.Model
	State CpaState `gorm:"state"`
}

type CpaState int

const (
	CpaStateInit CpaState = iota
	CpaStateExecuted
)
