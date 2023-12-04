package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalTherapy struct {
	gorm.Model

	MedHistoryID              uint
	MedicalTherapyTitle       string `json:"medicalTherapyTitle" validate:"required"`
	MedicalTherapyDescription string `json:"medicalTherapyDescription" `
	FromDate                  string `json:"fromDate" `
	ToDate                    string `json:"toDate" `
	Quantity                  string `json:"quantity" `
	Frequency                 string `json:"frequency" `
	CreatedOn                 time.Time
	Highlight                 bool    `json:"highlight"`
	Weight                    float32 `json:"weight"`
}
