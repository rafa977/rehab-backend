package models

import (
	"gorm.io/gorm"
)

type ClinicalTests struct {
	gorm.Model

	CompanyID              uint    `json:"companyId,omitempty"`
	Company                Company `json:"-"`
	Title                  string
	ClinicalTestCategoryID uint
	ClinicalTestCategory   ClinicalTestCategory
	Description            string
}
