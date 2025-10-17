package models

import "gorm.io/gorm"

type Lnt struct {
	gorm.Model

	State LntState `gorm:"state"`
}

type LntState int

const (
	LntStateInit LntState = iota
	LntStatePending
	LntStateOpened
	LntStateUnknown = -2
)
