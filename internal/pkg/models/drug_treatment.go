package models

import (
	"time"

	"github.com/uptrace/bun"
)

type DrugTreatment struct {
	bun.BaseModel `bun:"table:drug_treatment ,alias:drgtr"`

	DrugTreatmentID string `bun:",pk,autoincrement"`
	Drug            string `bun:"drug_id,notnull" json:"drug_id" validate:"required"`
	UserID          string `bun:"user_id,notnull" json:"user_id" validate:"required"`
	FromDate        string `bun:"from_date" json:"from_date" `
	ToDate          string `bun:"to_date" json:"to_date" `
	Quantity        string `bun:"quantity" json:"quantity" `
	Frequency       string `bun:"frequency" json:"frequency" `
	CreatedOn       time.Time
}
