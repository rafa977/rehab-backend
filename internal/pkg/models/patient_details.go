package models

import (
	"gorm.io/gorm"
)

type PatientDetails struct {
	gorm.Model

	PatientID       uint
	Patient         Patient   `json:"patient" validate:"-"`
	Title           string    `json:",omitempty"`
	Description     string    // Description
	Note            string    // Note
	AddedByID       uint      // New foreign key
	AddedBy         Account   `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	LastUpdatedByID uint      // New foreign key
	LastUpdatedBy   Account   `gorm:"foreignkey:LastUpdatedByID"` // LastUpdatedByID relationship
	Diseases        []Disease `json:"diseases,omitempty"`
}
