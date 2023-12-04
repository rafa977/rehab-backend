package models

import "gorm.io/gorm"

type Highlight struct {
	gorm.Model

	Weight    int    `json:"weight"`
	ModelName string `json:"modelName"`
	Value     string `json:"value"`
	FieldID   int    `json:"fieldId"`
	PatientId uint   `json:"patientId"`
}
