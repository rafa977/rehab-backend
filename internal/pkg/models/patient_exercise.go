package models

import "gorm.io/gorm"

type PatientExercise struct {
	gorm.Model

	PhTherapyID uint       `json:"phTherapyId"`
	PhTherapy   *PhTherapy `json:"phTherapy"`
	AccountID   uint       `json:"accountId"`
	Account     *Account   `json:"account"`
	Date        string     `json:"date"`
	Desription  string     `json:"description"`
	Completed   bool       `json:"completed"`
	Note        string     `json:"note"`
	Highlight   bool       `json:"highlight"`
	Weight      float32    `json:"weight"`
}
