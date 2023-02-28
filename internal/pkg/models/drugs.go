package models

import (
	"github.com/uptrace/bun"
)

type Drugs struct {
	bun.BaseModel `bun:"table:drugs,alias:drgs"`

	DrugID          string `bun:",pk,autoincrement"`
	DrugTitle       string `bun:"drug_title,notnull" json:"drug_title" validate:"required"`
	DrugDescription string `bun:"drug_description" json:"drug_description" `
	DrugNotes       string `bun:"drug_notes" json:"drug_notes" `
}
