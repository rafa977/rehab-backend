package models

import "time"

type PatientPersonal struct {
	Account          Account            `json:"account"`
	Therapy          []Therapy          `json:"therapies"`
	MedicalTherapy   []MedicalTherapy   `json:"medtherapies"`
	Injury           []Injury           `json:"injuries"`
	PersonalAllergy  []PersonalAllergy  `json:"persallergies"`
	PersonalDisorder []PersonalDisorder `json:"persdisorders"`
	DrugTreatment    []DrugTreatment    `json:"drugtreatments"`
	CreatedOn        time.Time
}
