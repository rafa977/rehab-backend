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
	Quantity           string  `json:"quantity"`
	Frequency          string  `json:"frequency"`
	Doctor             string  `json:"doctor"`
	Highlight          bool    `json:"highlight"`
	Weight             float32 `json:"weight"`
}
