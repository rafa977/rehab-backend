package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Therapy struct {
	bun.BaseModel `bun:"table:therapies ,alias:th"`

	TherapyID          string `bun:",pk,autoincrement"`
	UserID             string `bun:"user_id,notnull" json:"user_id" validate:"required"`
	TherapyTitle       string `bun:"therapy_title,notnull" json:"therapy_title" validate:"required"`
	TherapyDescription string `bun:"therapy_description" json:"therapy_description" `
	Diagnosis          string `bun:"diagnosis" json:"diagnosis" `
	FromDate           string `bun:"from_date" json:"from_date" `
	ToDate             string `bun:"to_date" json:"to_date" `
	Quantity           string `bun:"quantity" json:"quantity" `
	Frequency          string `bun:"frequency" json:"frequency" `
	CreatedOn          time.Time
}
