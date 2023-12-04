package models

import (
	"time"

	"gorm.io/gorm"
)

type DrugTreatment struct {
	gorm.Model

	MedHistoryID uint   `json:"medHistoryId"`
	DrugID       uint   `json:"drugId"`
	Drug         *Drug  `json:"drug"`
	UserID       string `json:"userId" validate:"required"`
	FromDate     string `json:"fromDate" `
	ToDate       string `json:"toDate" `
	Quantity     string `json:"quantity" `
	Frequency    string `json:"frequency" `
	CreatedOn    time.Time
	Highlight    bool    `json:"highlight"`
	Weight       float32 `json:"weight"`
}
