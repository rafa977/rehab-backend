package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name       string     `json:"name" validate:"required"`
	Address    string     `json:"address" validate:"required"`
	City       string     `json:"city"`
	PostalCode string     `json:"postalcode"`
	TAXID      string     `json:"taxid"`
	Relations  []Relation `json:",omitempty"`
	CreatedOn  time.Time
	LastLogin  time.Time `json:"lastlogin"`
}
