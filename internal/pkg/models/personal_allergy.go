package models

import "gorm.io/gorm"

type PersonalAllergy struct {
	gorm.Model

	PatientID     uint
	AllergyID     uint `json:"allergy_id"`
	Allergy       Allergy
	DiagnosedTime string `json:"diagnosed_time" `
}
