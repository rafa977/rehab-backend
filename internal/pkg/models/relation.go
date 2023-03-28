package models

import (
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	AccountID uint
	Account   Account
	Companies []Company `gorm:"many2many:relation_companies;"`
	Title     string
	Type      string
}
