package models

import (
	"gorm.io/gorm"
)

type PhTherapyKey struct {
	gorm.Model

	PhTherapyID uint
	PhTherapy   PhTherapy
	AccountID   uint
	Account     Account
	Description string
	Date        string
	Note        string
}
