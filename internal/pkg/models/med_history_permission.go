package models

import (
	"gorm.io/gorm"
)

type MedHistoryPermission struct {
	gorm.Model
	AccountID    uint        `json:"accountId"`
	Account      *Account    `json:"-"`
	FromDate     string      `json:"fromDate"`
	ToDate       string      `json:"toDate"`
	Access       bool        `json:"access"`
	MedHistoryID uint        `json:"medHistoryId"`
	MedHistory   *MedHistory `json:"medHistory"`
}
