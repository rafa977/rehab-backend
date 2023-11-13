package models

import (
	"gorm.io/gorm"
)

type GenericNote struct {
	gorm.Model

	Description string  // Description
	Note        string  // Note
	AddedByID   uint    // New foreign key
	AddedBy     Account `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	Commentator string  `json:"commentator"`
}
