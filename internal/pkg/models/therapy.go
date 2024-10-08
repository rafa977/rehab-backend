package models

import (
	"gorm.io/gorm"
)

type Therapy struct {
	gorm.Model

	MedHistoryID       uint
	TherapyTitle       string  `json:"therapyTitle" validate:"required"`
	TherapyDescription string  `json:"therapyDescription"`
	Diagnosis          string  `json:"diagnosis"`
	FromDate           string  `json:"fromDate"`
	ToDate             string  `json:"toDate"`
	Quantity           string  `json:"quantity"`
	Frequency          string  `json:"frequency"`
	Highlight          bool    `json:"highlight"`
	Weight             float32 `json:"weight"`
}
