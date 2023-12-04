package models

import (
	"gorm.io/gorm"
)

type ClinicalTests struct {
	gorm.Model

	CompanyID   uint     `gorm:"uniqueIndex:idx_companyid_title" json:"companyId,omitempty"`
	Company     *Company `json:"-" validate:"-"`
	Title       string   `gorm:"uniqueIndex:idx_companyid_title"`
	Description string
}
