package models

import (
	"gorm.io/gorm"
)

type Dysfunction struct {
	gorm.Model

	PatientDetailsID uint           `json:"patientDetailsId" validate:"required"`
	PatientDetails   PatientDetails `json:"patientDetails" validate:"-"`
	CompanyID        uint           `json:"companyId,omitempty" validate:"required"`
	Company          Company        `json:"-" validate:"-"`
	Dysfunction      string
	Description      string
	Note             string
	Date             string
	Frequency        string
	Timepervisit     string
	InjuryID         uint
	Protocols        []Protocol
}
