package models

type Response struct {
	Type     string `json:"type"`
	Status   string `json:"status" default:"null"`
	Date     string `json:"date"`
	Response string `json:"response"`
	Message  string `json:"message"`
}
