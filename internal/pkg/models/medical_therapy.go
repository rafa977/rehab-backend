package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalTherapy struct {
	gorm.Model

	PatientDetailsID          uint
	MedicalTherapyTitle       string `json:"medicalTherapyTitle" validate:"required"`
	MedicalTherapyDescription string `json:"medicalTherapyDescription" `
	FromDate                  string `json:"fromDate" `
	ToDate                    string `json:"toDate" `
	Quantity                  string `json:"quantity" `
	Frequency                 string `json:"frequency" `
	CreatedOn                 time.Time
}
