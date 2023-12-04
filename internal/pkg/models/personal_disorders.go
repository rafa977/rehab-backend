package models

import (
	"gorm.io/gorm"
)

type PersonalDisorder struct {
	gorm.Model

	MedHistoryID uint      `json:"medHistoryId"`
	DisorderID   uint      `json:"disorderId"`
	Disorder     *Disorder `json:"disorder"`
	FromDate     string    `json:"from_date"`
	ToDate       string    `json:"to_date"`
	Quantity     string    `json:"quantity"`
	Frequency    string    `json:"frequency"`
	Highlight    bool      `json:"highlight"`
	Weight       float32   `json:"weight"`
}
