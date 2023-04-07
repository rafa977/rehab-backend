package models

import "gorm.io/gorm"

type Protocol struct {
	gorm.Model

	DysfunctionID uint
	Dysfunction   Dysfunction
	Timetable     string
	Expectation   string
	Result        string
}
