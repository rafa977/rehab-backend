package models

type Notification struct {
	MailAddress string `json:"allergy_title"`
	URL         string `json:"allergy_description"`
	Subject     string `json:"allergy_notes"`
}
