package models

import (
	"time"

	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	AccountID int
	CompanyID int
	Title     string
	Type      string
	CreatedOn time.Time
	LastLogin time.Time
}
