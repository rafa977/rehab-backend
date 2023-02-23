package models

import (
	"github.com/uptrace/bun"
)

type Disorders struct {
	bun.BaseModel `bun:"table:disorders ,alias:dsrdr"`

	DisorderID          string `bun:",pk,autoincrement"`
	DisorderTitle       string `bun:"disorder_title,notnull" json:"disorder_title" validate:"required"`
	DisorderDescription string `bun:"disorder_description" json:"disorder_description" `
	DisorderNotes       string `bun:"disorder_notes" json:"disorder_notes" `
}
