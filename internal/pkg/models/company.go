package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name       string     `json:"name" validate:"required"`
	Address    string     `json:"address" validate:"required"`
	City       string     `json:"city"`
	PostalCode string     `json:"postalcode"`
	TAXID      string     `gorm:"unique" json:"taxid"`
	Relations  []Relation `gorm:"many2many:relation_companies;"`
}
