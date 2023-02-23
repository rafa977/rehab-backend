package models

import (
	"github.com/uptrace/bun"
)

type Allergy struct {
	bun.BaseModel `bun:"table:allergies ,alias:all"`

	AllergyID          string `bun:",pk,autoincrement"`
	AllergyTitle       string `bun:"allergy_title,notnull" json:"allergy_title" validate:"required"`
	AllergyDescription string `bun:"allergy_description" json:"allergy_description" `
	AllergyNotes       string `bun:"allergy_notes" json:"allergy_notes" `
}
