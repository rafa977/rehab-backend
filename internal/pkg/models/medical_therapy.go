package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalTherapy struct {
	gorm.Model

	PatientID                 string
	MedicalTherapyTitle       string `json:"medical_therapy_title" validate:"required"`
	MedicalTherapyDescription string `json:"medical_therapy_description" `
	FromDate                  string `json:"from_date" `
	ToDate                    string `json:"to_date" `
	Quantity                  string `json:"quantity" `
	Frequency                 string `json:"frequency" `
	CreatedOn                 time.Time
}
