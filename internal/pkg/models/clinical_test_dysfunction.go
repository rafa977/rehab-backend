package models

import (
	"gorm.io/gorm"
)

type ClinicalTestDysfunction struct {
	gorm.Model

	DysfunctionID   uint
	Dysfunction     Dysfunction `validate:"-"`
	ClinicalTestsID uint
	ClinicalTests   ClinicalTests `validate:"-"`
	Position        string
	Score           string
	Note            string
}
