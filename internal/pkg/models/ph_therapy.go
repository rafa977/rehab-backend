package models

import (
	"gorm.io/gorm"
)

type PhTherapy struct {
	gorm.Model

	PatientDetailsID  uint              `json:"patientDetailsId" validate:"required"`
	PatientDetails    *PatientDetails   `json:"patientDetails" validate:"-"`
	Date              string            `json:"date"`
	EmployeeID        uint              `json:"employeeId"`
	AccountEmployee   *Account          `gorm:"foreignKey:EmployeeID"`
	SupervisorID      uint              `json:"supervisorId"`
	AccountSuperVisor *Account          `gorm:"foreignKey:SupervisorID"`
	Description       string            `json:"description"`
	TherapyNumber     int64             `json:"therapyNumber"`
	Notes             []PhTherapyNote   `json:"phTherapyNotes,omitempty"`
	TherapyKeys       []PhTherapyKey    `json:"phTherapyKeys,omitempty"`
	Exercises         []PatientExercise `json:"patientExercise,omitempty"`
	Protocols         []Protocol        `json:"protocols,omitempty"`
}
