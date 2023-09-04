package models

import (
	"gorm.io/gorm"
)

type PatientDetails struct {
	gorm.Model

	PatientID         uint
	Patient           Patient            `json:"patient" validate:"-"`
	AddedByID         uint               // New foreign key
	AddedBy           Account            `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	LastUpdatedByID   uint               // New foreign key
	LastUpdatedBy     Account            `gorm:"foreignkey:LastUpdatedByID"` // LastUpdatedByID relationship
	Therapies         []Therapy          `json:"therapies,omitempty"`
	MedicalTherapies  []MedicalTherapy   `json:"medTherapies,omitempty"`
	DrugTreatments    []DrugTreatment    `json:"drugTreatments,omitempty"`
	Injuries          []Injury           `json:"injuries,omitempty"`
	PersonalAllergies []PersonalAllergy  `json:"personalAllergies,omitempty"`
	PersonalDisorders []PersonalDisorder `json:"persnoalDisorders,omitempty"`
	Dysfunctions      []Dysfunction      `json:"dysfunctions,omitempty"`
	PhTherapies       []PhTherapy        `json:"phTherapies,omitempty"`
}
