package models

import (
	"gorm.io/gorm"
)

type Visit struct {
	gorm.Model

	PatientID uint
	Patient   Patient
	CompanyID uint
	Company   Company
}
