package models

import (
	"time"

	"gorm.io/gorm"
)

type DrugTreatment struct {
	gorm.Model

	MedHistoryID uint
	DrugID       uint
	Drug         Drug
	UserID       string `json:"user_id" validate:"required"`
	FromDate     string `json:"from_date" `
	ToDate       string `json:"to_date" `
	Quantity     string `json:"quantity" `
	Frequency    string `json:"frequency" `
	CreatedOn    time.Time
}
