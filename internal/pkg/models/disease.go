package models

import (
	"gorm.io/gorm"
)

type Disease struct {
	gorm.Model

	PatientDetailsID    uint           `json:"patientDetailsId" validate:"required"`
	PatientDetails      PatientDetails `json:"patientDetails" validate:"-"`
	Disease             string
	Description         string
	Note                string
	Date                string
	InjuryID            uint
	Dysfunctions        []Dysfunction         `json:"dysfunctions,omitempty"`
	ClinicalTestDisease []ClinicalTestDisease `json:"clinicalTestDisease,omitempty"`
	DiseaseHistory      []DiseaseHistory      `json:"diseaseHistory,omitempty"`
	PhTherapies         []PhTherapy           `json:"phTherapies,omitempty"`
}
