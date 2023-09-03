package models

import (
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	AccountID uint
	Account   Account
	Companies []Company `gorm:"many2many:relation_companies;"`
	AddedByID uint      // New foreign key
	AddedBy   Account   `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	Title     string
	Type      string
}
