package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PersonalDisorder struct {
	bun.BaseModel `bun:"table:personal_disorders,alias:prsdsrd"`

	PersonalDisorderID string `bun:"personal_disorder_id,pk,autoincrement"`
	DisorderID         string `bun:"disorder_id,notnull" json:"disorder_id" validate:"required"`
	UserID             string `bun:"user_id,notnull" json:"user_id" validate:"required"`
	FromDate           string `bun:"from_date" json:"from_date" `
	ToDate             string `bun:"to_date" json:"to_date" `
	Quantity           string `bun:"quantity" json:"quantity" `
	Frequency          string `bun:"frequency" json:"frequency" `
	CreatedOn          time.Time
}
