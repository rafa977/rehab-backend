package models

type Notification struct {
	MailAddress string `json:"address"`
	URL         string `json:"url"`
	Subject     string `json:"subject"`
}
