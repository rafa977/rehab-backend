package models

import (
	"gorm.io/gorm"
)

type PhTherapyKey struct {
	gorm.Model

	PhTherapyID uint       `json:"phTherapyId"`
	PhTherapy   *PhTherapy `json:"phTherapy"`
	AccountID   uint       `json:"accountId"`
	Account     *Account   `json:"account"`
	Description string     `json:"description"`
	Date        string     `json:"date"`
	Note        string     `json:"note"`
	Highlight   bool       `json:"highlight"`
	Weight      float32    `json:"weight"`
}
