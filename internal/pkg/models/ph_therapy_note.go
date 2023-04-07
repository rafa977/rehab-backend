package models

import "gorm.io/gorm"

type PhTherapyNote struct {
	gorm.Model

	PhTherapyID uint
	PhTherapy   PhTherapy
	AccountID   uint
	Account     Account
	Date        string
	Note        string
}
