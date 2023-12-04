package models

import (
	"gorm.io/gorm"
)

type GenericNote struct {
	gorm.Model

	Description string   `json:"description"`
	Note        string   `json:"note"`
	AddedByID   uint     `json:"addedById"`
	AddedBy     *Account `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	Commentator string   `json:"commentator"`
	Highlight   bool     `json:"highlight"`
	Weight      float32  `json:"weight"`
}
