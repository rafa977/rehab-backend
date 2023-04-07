package models

import (
	"gorm.io/gorm"
)

type PhTherapy struct {
	gorm.Model

	PatientDetailsID uint
	PatientDetails   PatientDetails
	DysfunctionID    uint
	Dysfunction      Dysfunction
	Date             string
	EmployeeID       uint
	SupervisorID     uint
	Description      string
	Notes            []PhTherapyNote   `json:"phTherapyNotes,omitempty"`
	TherapyKeys      []PhTherapyKey    `json:"phTherapyKeys,omitempty"`
	Exercises        []PatientExercise `json:"patientExercise,omitempty"`
}
