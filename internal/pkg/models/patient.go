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
	Amka           int `gorm:"unique"  json:"amka" validate:"required"`
	Age            int `json:"age" validate:"required"`
	Job            string
	CompanyID      uint             `json:"companyId,omitempty"` // createonly (disabled read from db)
	Company        Company          `json:"-" validate:"-"`
	PatientDetails []PatientDetails `json:"-"`
}
