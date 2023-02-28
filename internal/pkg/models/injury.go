package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Injury struct {
	bun.BaseModel `bun:"table:injuries,alias:inj"`

	InjuryID          string `bun:"injury_id,pk,autoincrement"`
	UserID            string `bun:"user_id,notnull" json:"user_id" validate:"required"`
	InjuryTitle       string `bun:"injury_title,notnull" json:"injury_title" validate:"required"`
	InjuryDescription string `bun:"injury_description" json:"injury_description" `
	BpositionID       string `bun:"bposition_id" json:"bposition_id" `
	InjuryDate        string `bun:"injury_date" json:"injury_date" `
	CreatedOn         time.Time
}
