package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PersonalAllergy struct {
	bun.BaseModel `bun:"table:personal_allergies ,alias:persall"`

	PersonalAllergyID string `bun:",pk,autoincrement"`
	AllergyID         string `bun:"allergy_id,notnull" json:"allergy_id" validate:"required"`
	UserID            string `bun:"user_id,notnull" json:"user_id" validate:"required"`
	DiagnosedDate     string `bun:"diagnosed_time" json:"diagnosed_time" `
	CreatedOn         time.Time
}
