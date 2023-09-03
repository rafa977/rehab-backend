package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Firstname      string `json:"firstname" validate:"required"`
	Lastname       string `json:"lastname" validate:"required"`
	Email          string `json:"email" validate:"required,email,max=128"`
	Address        string
	Amka           int        `gorm:"uniqueIndex:idx_companyid_amka"  json:"amka" validate:"required"`
	Birthdate      CustomDate `gorm:"embedded" json:"birthdate" validate:"required"`
	Job            string     `json:"phone"`
	Phone          string
	CompanyID      uint             `gorm:"uniqueIndex:idx_companyid_amka" json:"companyId,omitempty"` // createonly (disabled read from db)
	Company        Company          `validate:"-"`
	PatientDetails []PatientDetails `json:"-"`
}
