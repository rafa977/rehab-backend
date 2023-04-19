package models

import (
	"gorm.io/gorm"
)

type ClinicalTestDysfunction struct {
	gorm.Model

	DysfunctionID  uint
	Dysfunction    Dysfunction
	ClinicalTestID uint
	ClinicalTests  ClinicalTests
	Position       string
	Score          string
	Note           string
}
