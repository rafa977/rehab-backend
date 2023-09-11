package models

import (
	"gorm.io/gorm"
)

type PatientPermission struct {
	gorm.Model
	AccountID        uint
	Account          Account
	FromDate         string `json:"from_date"`
	ToDate           string `json:"to_date"`
	Access           bool
	AccessToHistory  bool
	SectionsOfAccess []string
	PatientID        uint
	Patient          Patient `json:"-"`
}
