package models

type AccountInt struct {
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Address   string     `json:"address"`
	Birthdate CustomDate `json:"birthdate"`
	Job       string     `json:"job"`
}
