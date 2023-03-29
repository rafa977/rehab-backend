package models

import (
	"gorm.io/gorm"
)

type PatientDetails struct {
	gorm.Model

	PatientID         uint
	Patient           Patient
	CompanyID         uint    `json:"companyId,omitempty"`
	Company           Company `json:"-"`
	CreatedBy         uint
	LastUpdatedBy     uint
	Therapies         []Therapy          `json:"therapies,omitempty"`
	MedicalTherapies  []MedicalTherapy   `json:"medTherapies,omitempty"`
	DrugTreatments    []DrugTreatment    `json:"drugTreatments,omitempty"`
	Injuries          []Injury           `json:"injuries,omitempty"`
	PersonalAllergies []PersonalAllergy  `json:"personalAllergies,omitempty"`
	PersonalDisorders []PersonalDisorder `json:"persnoalDisorders,omitempty"`
}
