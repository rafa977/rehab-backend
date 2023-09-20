package models

import (
	"gorm.io/gorm"
)

type Dysfunction struct {
	gorm.Model

	DiseaseID          uint    `json:"diseaseId" validate:"required"`
	Disease            Disease `json:"disease" validate:"-"`
	Dysfunction        string
	Description        string
	Note               string
	Date               string
	InjuryID           uint
	DysfunctionHistory []DysfunctionHistory
}
