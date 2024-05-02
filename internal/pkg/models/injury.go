package models

import "gorm.io/gorm"

type Injury struct {
	gorm.Model

	MedHistoryID      uint
	InjuryTitle       string  `json:"injuryTitle"`
	InjuryDescription string  `json:"injuryDescription"`
	BpositionID       string  `json:"bpositionId"`
	InjuryDate        string  `json:"injuryDate"`
	Highlight         bool    `json:"highlight"`
	Weight            float32 `json:"weight"`
}
