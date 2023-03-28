package models

import (
	"gorm.io/gorm"
)

type PersonalDisorder struct {
	gorm.Model

	PatientID  uint
	DisorderID string `json:"disorder_id" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Quantity   string `json:"quantity"`
	Frequency  string `json:"frequency"`
}
