package models

import (
	"gorm.io/gorm"
)

type PersonalDisorder struct {
	gorm.Model

	PatientDetailsID uint
	DisorderID       uint
	Disorder         Disorder
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
	Quantity         string `json:"quantity"`
	Frequency        string `json:"frequency"`
}
