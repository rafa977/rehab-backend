package models

import "gorm.io/gorm"

type Injury struct {
	gorm.Model

	PatientID         uint
	InjuryTitle       string `json:"injury_title"`
	InjuryDescription string `json:"injury_description"`
	BpositionID       string `json:"bposition_id"`
	InjuryDate        string `json:"injury_date"`
}
