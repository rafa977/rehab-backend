package models

import (
	"gorm.io/gorm"
)

type Dysfunction struct {
	gorm.Model

	PatientDetailsID   uint            `json:"patientDetailsId" validate:"required"`
	PatientDetails     *PatientDetails `json:"patientDetails" validate:"-"`
	Dysfunction        string          `json:"dysfunction"`
	Description        string          `json:"description"`
	Note               string          `json:"note"`
	Date               string          `json:"date"`
	InjuryID           uint            `json:"injuryId"`
	Highlight          bool            `json:"highlight"`
	Weight             float32         `json:"weight"`
	DysfunctionHistory []DysfunctionHistory
}
