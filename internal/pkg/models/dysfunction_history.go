package models

import (
	"gorm.io/gorm"
)

type DiseaseHistory struct {
	gorm.Model

	DiseaseID uint
	Disease   Disease
	Note      string
}
