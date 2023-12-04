package models

import (
	"gorm.io/gorm"
)

type ClinicalTestDisease struct {
	gorm.Model

	PatientDetailsID uint            `json:"patientDetailsId" validate:"required"`
	PatientDetails   *PatientDetails `json:"patientDetails" validate:"-"`
	ClinicalTestsID  uint
	ClinicalTests    *ClinicalTests `json:"clinicalTests" validate:"-"`
	Position         string
	Score            string
	Note             string
	Highlight        bool    `json:"highlight"`
	Weight           float32 `json:"weight"`
}
