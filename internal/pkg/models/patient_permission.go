package models

import (
	"gorm.io/gorm"
)

type PatientPermission struct {
	gorm.Model
	AccountID        uint
	Account          *Account
	FromDate         string   `json:"fromDate"`
	ToDate           string   `json:"toDate"`
	Access           bool     `json:"access"`
	AccessToHistory  bool     `json:"accessToHistory"`
	SectionsOfAccess []string `json:"sectionsOfAccess"`
	PatientID        uint
	Patient          *Patient `json:"-"`
}
