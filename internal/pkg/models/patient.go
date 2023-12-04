package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Firstname      string           `json:"firstname" validate:"required"`
	Lastname       string           `json:"lastname" validate:"required"`
	Email          string           `json:"email" validate:"required,email,max=128"`
	Address        string           `json:"address"`
	Amka           int              `gorm:"uniqueIndex:idx_companyid_amka"  json:"amka" validate:"required"`
	Birthdate      CustomDate       `gorm:"embedded" json:"birthdate" validate:"required"`
	Job            string           `json:"job"`
	Phone          string           `json:"phone"`
	PhotoUrl       string           `json:"photoUrl"`
	RecommendedBy  string           `json:"recommendedBy"`
	Drugs          bool             `json:"drugs"`
	Thyroid        bool             `json:"thyroid"`
	Diabetes       bool             `json:"diabetes"`
	Smoking        bool             `json:"smoking"`
	CompanyID      uint             `gorm:"uniqueIndex:idx_companyid_amka" json:"companyId,omitempty"` // createonly (disabled read from db)
	Company        *Company         `validate:"-"`
	AddedByID      uint             // New foreign key
	AddedBy        *Account         `gorm:"foreignkey:AddedByID" validate:"-" json:"addedBy,omitempty"` // AddedBy relationship
	PatientDetails []PatientDetails `json:"-"`
	TherapyHistory []TherapyHistory `json:"therapyHistory,omitempty"`
}

type PatientEmployee struct {
	ID        uint
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email,max=128"`
	Phone     string
	CompanyID uint    `gorm:"uniqueIndex:idx_companyid_amka" json:"-"` // createonly (disabled read from db)
	Company   Company `json:"-" validate:"-"`
	AddedByID uint    `json:"-"`                             // New foreign key
	AddedBy   Account `gorm:"foreignkey:AddedByID" json:"-"` // AddedBy relationship
}
