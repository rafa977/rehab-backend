package models

import (
	"gorm.io/gorm"
)

type PatientDetails struct {
	gorm.Model

	PatientID           uint
	Patient             *Patient              `json:"patient,omitempty" validate:"-" `
	Title               string                `json:",omitempty"`
	Description         string                `json:"description"`
	Note                string                `json:"note"` // Note
	AddedByID           uint                  // New foreign key
	AddedBy             *Account              `json:",omitempty" gorm:"foreignkey:AddedByID"` // AddedBy relationship
	LastUpdatedByID     uint                  // New foreign key
	LastUpdatedBy       *Account              `json:",omitempty" gorm:"foreignkey:LastUpdatedByID"` // LastUpdatedByID relationship
	Disease             string                `json:"disease"`
	Date                string                `json:"data"`
	SessionsPerWeek     uint                  `json:"sessionsPerWeek"`
	TotalTherapies      uint                  `json:"totalTherapies"`
	StartDate           string                `json:"startDate"`
	Dysfunctions        []Dysfunction         `json:"dysfunctions,omitempty"`
	ClinicalTestDisease []ClinicalTestDisease `json:"clinicalTestDisease,omitempty"`
	DiseaseHistory      []DiseaseHistory      `json:"diseaseHistory,omitempty"`
	PhTherapies         []PhTherapy           `json:"phTherapies,omitempty"`
}
