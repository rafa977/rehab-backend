package models

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email,max=128"`
	Address   string
	Amka      string `json:"amka" validate:"required"`
	Age       string `json:"age" validate:"required"`
	Job       string
	CreatedOn time.Time
	LastLogin time.Time `json:"lastlogin" default:"null"`
}
