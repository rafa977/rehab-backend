package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username  string `gorm:"unique" json:"username" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `gorm:"unique" json:"email" validate:"required,email,max=128"`
	Address   string
	Amka      int        `json:"amka"`
	Birthdate CustomDate `gorm:"embedded" json:"birthdate"`
	Job       string     `gorm:"default:null"`
	Relations []Relation
	LastLogin time.Time `json:"lastlogin"`
}
