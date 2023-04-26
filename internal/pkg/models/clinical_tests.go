package models

import (
	"gorm.io/gorm"
)

type ClinicalTests struct {
	gorm.Model

	CompanyID              uint    `json:"companyId,omitempty"`
	Company                Company `json:"-" validate:"-"`
	Title                  string  `gorm:"unique"`
	ClinicalTestCategoryID uint
	ClinicalTestCategory   ClinicalTestCategory `json:"-"`
	Description            string
}
