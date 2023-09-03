package models

import (
	"gorm.io/gorm"
)

type PatientDetailsPermission struct {
	gorm.Model
	AccountID        uint
	Account          Account
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
	PatientDetailsID uint
	PatientDetails   PatientDetails `json:"-"`
}
