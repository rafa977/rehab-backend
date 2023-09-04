package models

import (
	"gorm.io/gorm"
)

type Signature struct {
	gorm.Model
	PatientID uint
	Patient   Patient
	Signature string
}
