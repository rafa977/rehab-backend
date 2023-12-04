package models

import (
	"gorm.io/gorm"
)

type PatientDetailsPermission struct {
	gorm.Model
	AccountID        uint
	Account          *Account `json:"-"`
	FromDate         string   `json:"fromDate"`
	ToDate           string   `json:"toDate"`
	Access           bool     `json:"access"`
	PatientDetailsID uint     `json:"patientDetailsId"`
	PatientDetails   *PatientDetails
}
