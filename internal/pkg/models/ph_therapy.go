package models

import (
	"gorm.io/gorm"
)

type PhTherapy struct {
	gorm.Model

	DiseaseID         uint
	Disease           Disease
	Date              string
	EmployeeID        uint
	AccountEmployee   Account `gorm:"foreignKey:EmployeeID"`
	SupervisorID      uint
	AccountSuperVisor Account `gorm:"foreignKey:SupervisorID"`
	Description       string
	TherapyNumber     int64
	Notes             []PhTherapyNote   `json:"phTherapyNotes,omitempty"`
	TherapyKeys       []PhTherapyKey    `json:"phTherapyKeys,omitempty"`
	Exercises         []PatientExercise `json:"patientExercise,omitempty"`
	Protocols         []Protocol        `json:"protocols,omitempty"`
}
