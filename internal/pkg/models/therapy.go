package models

import (
	"gorm.io/gorm"
)

type Therapy struct {
	gorm.Model

	PatientDetailsID   uint
	TherapyTitle       string `json:"therapyTitle" validate:"required"`
	TherapyDescription string `json:"therapyDescription"`
	Diagnosis          string `json:"diagnosis"`
	FromDate           string `json:"fromDate"`
	ToDate             string `json:"toDate"`
	Quantity           string `json:"quantity"`
	Frequency          string `json:"frequency"`
}
