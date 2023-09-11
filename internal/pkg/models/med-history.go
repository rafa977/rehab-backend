package models

import (
	"gorm.io/gorm"
)

type MedHistory struct {
	gorm.Model

	PatientID         uint
	CompanyID         uint               `gorm:"uniqueIndex:idx_companyid_amka" json:"companyId,omitempty"` // createonly (disabled read from db)
	Company           Company            `validate:"-"`
	AddedByID         uint               // New foreign key
	AddedBy           Account            `gorm:"foreignkey:AddedByID" validate:"-"` // AddedBy relationship
	Therapies         []Therapy          `json:"therapies,omitempty"`
	MedicalTherapies  []MedicalTherapy   `json:"medTherapies,omitempty"`
	DrugTreatments    []DrugTreatment    `json:"drugTreatments,omitempty"`
	Injuries          []Injury           `json:"injuries,omitempty"`
	PersonalAllergies []PersonalAllergy  `json:"personalAllergies,omitempty"`
	PersonalDisorders []PersonalDisorder `json:"persnoalDisorders,omitempty"`
}
