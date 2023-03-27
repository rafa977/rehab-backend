package models

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Firstname        string `json:"firstname" validate:"required"`
	Lastname         string `json:"lastname" validate:"required"`
	Email            string `json:"email" validate:"required,email,max=128"`
	Address          string
	Amka             int `json:"amka" validate:"required"`
	Age              int `json:"age" validate:"required"`
	Job              string
	CreatedOn        time.Time
	LastLogin        time.Time        `json:"lastlogin" default:"null"`
	CompanyID        int              `gorm:"foreignkey:CompanyID"`
	Therapies        []Therapy        `json:"therapies,omitempty"`
	MedicalTherapies []MedicalTherapy `json:"medTherapies,omitempty"`
	DrugTreatments   []DrugTreatment  `json:"drugTreatments,omitempty"`
}
