package models

import (
	"gorm.io/gorm"
)

type PatientDetails struct {
	gorm.Model

	PatientID       uint
	Patient         Patient       `json:"patient" validate:"-"`
	Title           string        // Patient details card title
	Description     string        // Description
	Note            string        // Note
	AddedByID       uint          // New foreign key
	AddedBy         Account       `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	LastUpdatedByID uint          // New foreign key
	LastUpdatedBy   Account       `gorm:"foreignkey:LastUpdatedByID"` // LastUpdatedByID relationship
	Dysfunctions    []Dysfunction `json:"dysfunctions,omitempty"`
}
