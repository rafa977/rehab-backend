package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" validate:"required,email,max=128"`
	Address   string
	Amka      string `json:"amka"`
	Age       string `json:"age"`
	Job       string `gorm:"default:null"`
	Relations []Relation
	CreatedOn time.Time
	LastLogin time.Time `json:"lastlogin"`
}
