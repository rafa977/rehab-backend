package models

import (
	"gorm.io/gorm"
)

type ClinicalTestCategory struct {
	gorm.Model

	CompanyID uint     `gorm:"uniqueIndex:idx_companyid_name" json:"companyId,omitempty"`
	Company   *Company `json:"-" validate:"-"`
	Name      string   `gorm:"uniqueIndex:idx_companyid_name"`
}
