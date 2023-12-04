package models

import (
	"gorm.io/gorm"
)

type DysfunctionHistory struct {
	gorm.Model

	DysfunctionID uint
	Dysfunction   *Dysfunction
	Note          string
}
