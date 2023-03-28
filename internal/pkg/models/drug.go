package models

import "gorm.io/gorm"

type Drug struct {
	gorm.Model

	DrugTitle       string          `json:"drug_title" validate:"required"`
	DrugDescription string          `json:"drug_description"`
	DrugNotes       string          `json:"drug_notes"`
	DrugTreatments  []DrugTreatment `json:"-"`
}
