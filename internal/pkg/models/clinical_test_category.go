package models

import (
	"gorm.io/gorm"
)

type ClinicalTestCategory struct {
	gorm.Model

	CompanyID uint    `json:"companyId,omitempty"`
	Company   Company `json:"-"`
	Name      string
}
