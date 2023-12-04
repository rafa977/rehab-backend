package models

import "gorm.io/gorm"

type Protocol struct {
	gorm.Model

	PhTherapyID uint
	PhTherapy   PhTherapy
	Timetable   string
	Expectation string
	Result      string
	Highlight   bool    `json:"highlight"`
	Weight      float32 `json:"weight"`
}
