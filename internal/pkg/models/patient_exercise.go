package models

import "gorm.io/gorm"

type PatientExercise struct {
	gorm.Model

	PhTherapyID uint
	PhTherapy   PhTherapy
	AccountID   uint
	Account     Account
	Date        string
	Desription  string
	Completed   bool
	Note        string
}
