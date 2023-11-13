package models

import (
	"gorm.io/gorm"
)

type MedHistoryPermission struct {
	gorm.Model
	AccountID    uint
	Account      Account `json:"-"`
	FromDate     string  `json:"from_date"`
	ToDate       string  `json:"to_date"`
	Access       bool
	MedHistoryID uint
	MedHistory   MedHistory
}
