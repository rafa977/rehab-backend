package models

import "gorm.io/gorm"

type Injury struct {
	gorm.Model

	MedHistoryID      uint
	InjuryTitle       string  `json:"injury_title"`
	InjuryDescription string  `json:"injury_description"`
	BpositionID       string  `json:"bposition_id"`
	InjuryDate        string  `json:"injury_date"`
	Highlight         bool    `json:"highlight"`
	Weight            float32 `json:"weight"`
}
