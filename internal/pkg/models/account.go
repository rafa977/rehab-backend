package models

import (
	"time"

	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username  string     `gorm:"unique" json:"username" validate:"required"`
	Password  string     `json:"password,omitempty" validate:"required"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Email     string     `gorm:"unique" json:"email" validate:"required,email,max=128"`
	Address   string     `json:"address"`
	Amka      int        `json:"amka"`
	Birthdate CustomDate `gorm:"embedded" json:"birthdate"`
	Job       string     `gorm:"default:null"`
	RoleID    uint       `json:"roleId"`
	Role      *Role      `validate:"-"`
	Relations []Relation `json:"relations"`
	LastLogin time.Time  `json:"lastlogin"`
}

func NewSmartRegisterLink() SmartRegisterLink {
	return SmartRegisterLink{
		Token: generateRandomString(),
	}
}

func generateRandomString() string {
	code := randstr.String(20)
	return code
}

type SmartRegisterLink struct {
	gorm.Model
	CompanyID uint
	Company   Company `gorm:"foreignkey:CompanyID" json:"-"`
	AddedByID uint    // New foreign key
	AddedBy   Account `gorm:"foreignkey:AddedByID"` // AddedBy relationship
	Email     string  `gorm:"size:255;unique" json:"email"`
	Token     string
}
