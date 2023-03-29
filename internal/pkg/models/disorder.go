package models

import "gorm.io/gorm"

type Disorder struct {
	gorm.Model

	DisorderTitle       string             `json:"disorderTitle"`
	DisorderDescription string             `json:"disorderDescription"`
	DisorderNotes       string             `json:"disorderNotes"`
	PersonalDisorders   []PersonalDisorder `json:",omitempty"`
}
