package models

import (
	"gorm.io/gorm"
)

type DiseaseHistory struct {
	gorm.Model

	PatientDetailsID uint            `json:"patientDetailsId" validate:"required"`
	PatientDetails   *PatientDetails `json:"patientDetails" validate:"-"`
	Highlight        bool            `json:"highlight"`
	Weight           float32         `json:"weight"`
	Note             string
}
