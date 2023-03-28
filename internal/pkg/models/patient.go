package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Firstname         string `json:"firstname" validate:"required"`
	Lastname          string `json:"lastname" validate:"required"`
	Email             string `json:"email" validate:"required,email,max=128"`
	Address           string
	Amka              int `json:"amka" validate:"required"`
	Age               int `json:"age" validate:"required"`
	Job               string
	CompanyID         int
	Therapies         []Therapy          `json:"therapies,omitempty"`
	MedicalTherapies  []MedicalTherapy   `json:"medTherapies,omitempty"`
	DrugTreatments    []DrugTreatment    `json:"drugTreatments,omitempty"`
	Injuries          []Injury           `json:"injuries,omitempty"`
	PersonalAllergies []PersonalAllergy  `json:"personalAllergies,omitempty"`
	PersonalDisorders []PersonalDisorder `json:"persnoalDisorders,omitempty"`
}
