package models

import "gorm.io/gorm"

type PersonalAllergy struct {
	gorm.Model

	MedHistoryID  uint     `json:"medHistoryId"`
	AllergyID     uint     `json:"allergyId"`
	Allergy       *Allergy `json:"allergy"`
	DiagnosedTime string   `json:"diagnosedTime"`
	Highlight     bool     `json:"highlight"`
	Weight        float32  `json:"weight"`
}
