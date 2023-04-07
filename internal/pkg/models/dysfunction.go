package models

import (
	"gorm.io/gorm"
)

type Dysfunction struct {
	gorm.Model

	PatientDetailsID uint `json:"patientDetailsId"`
	PatientDetails   PatientDetails
	CompanyID        uint    `json:"companyId,omitempty"`
	Company          Company `json:"-"`
	Dysfunction      string
	Description      string
	Note             string
	Date             string
	Frequency        string
	Timepervisit     string
	InjuryID         uint
	Protocols        []Protocol
}
