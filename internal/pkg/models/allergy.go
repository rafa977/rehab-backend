package models

import "gorm.io/gorm"

type Allergy struct {
	gorm.Model

	AllergyTitle       string            `json:"allergy_title"`
	AllergyDescription string            `json:"allergy_description"`
	AllergyNotes       string            `json:"allergy_notes"`
	PersonalAllergies  []PersonalAllergy `json:"-"`
}
