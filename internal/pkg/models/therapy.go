package models

import (
	"time"

	"gorm.io/gorm"
)

type Therapy struct {
	gorm.Model

	PatientID          string
	TherapyTitle       string `json:"therapy_title" validate:"required"`
	TherapyDescription string `json:"therapy_description" `
	Diagnosis          string `json:"diagnosis" `
	FromDate           string `json:"from_date" `
	ToDate             string `json:"to_date" `
	Quantity           string `json:"quantity" `
	Frequency          string `json:"frequency" `
	CreatedOn          time.Time
}
