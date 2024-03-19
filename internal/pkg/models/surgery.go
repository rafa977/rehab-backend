package models

import (
	"gorm.io/gorm"
)

type Surgery struct {
	gorm.Model

	MedHistoryID       uint
	SurgeryTitle       string  `json:"surgeryTitle" validate:"required"`
	SurgeryDescription string  `json:"surgeryDescription"`
	Diagnosis          string  `json:"diagnosis"`
	Note               string  `json:"note"`
	Date               string  `json:"date"`
	Doctor             string  `json:"doctor"`
	Highlight          bool    `json:"highlight"`
	Weight             float32 `json:"weight"`
}
