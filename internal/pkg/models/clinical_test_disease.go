package models

import (
	"gorm.io/gorm"
)

type ClinicalTestDisease struct {
	gorm.Model

	DiseaseID       uint
	Disease         Disease `validate:"-"`
	ClinicalTestsID uint
	ClinicalTests   ClinicalTests `json:"clinicalTests" validate:"-"`
	Position        string
	Score           string
	Note            string
}
