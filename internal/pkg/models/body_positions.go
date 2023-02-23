package models

import (
	"github.com/uptrace/bun"
)

type BodyPositions struct {
	bun.BaseModel `bun:"table:body_positions ,alias:bdpst"`

	DisorderID           string `bun:",pk,autoincrement"`
	BpositionTitle       string `bun:"bposition_title,notnull" json:"bposition_title" validate:"required"`
	BpositionDescription string `bun:"bposition_description" json:"bposition_description" `
	BpositionNotes       string `bun:"bposition_notes" json:"bposition_notes" `
}
