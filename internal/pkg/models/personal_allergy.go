package models

import "gorm.io/gorm"

type PersonalAllergy struct {
	gorm.Model

	MedHistoryID  uint
	AllergyID     uint
	Allergy       Allergy
	DiagnosedTime string
}
