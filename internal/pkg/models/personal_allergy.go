package models

import "gorm.io/gorm"

type PersonalAllergy struct {
	gorm.Model

	PatientDetailsID uint
	AllergyID        uint
	Allergy          Allergy
	DiagnosedTime    string
}
