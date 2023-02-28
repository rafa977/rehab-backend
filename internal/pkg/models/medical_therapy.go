package models

import (
	"time"

	"github.com/uptrace/bun"
)

type MedicalTherapy struct {
	bun.BaseModel `bun:"table:medical_therapies,alias:medth"`

	MedicalTherapyID          string `bun:",pk,autoincrement"`
	UserID                    string `bun:"user_id,notnull" json:"user_id" validate:"required"`
	MedicalTherapyTitle       string `bun:"medical_therapy_title,notnull" json:"medical_therapy_title" validate:"required"`
	MedicalTherapyDescription string `bun:"medical_therapy_description" json:"medical_therapy_description" `
	FromDate                  string `bun:"from_date" json:"from_date" `
	ToDate                    string `bun:"to_date" json:"to_date" `
	Quantity                  string `bun:"quantity" json:"quantity" `
	Frequency                 string `bun:"frequency" json:"frequency" `
	CreatedOn                 time.Time
}
