package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" validate:"required,email,max=128"`
	Address   string
	Amka      string `json:"amka"`
	Age       string `json:"age"`
	Job       string
	CreatedOn time.Time
	LastLogin time.Time `json:"lastlogin"`
}
