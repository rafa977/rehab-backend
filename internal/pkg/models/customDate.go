package models

import "time"

type CustomDate struct {
	Birthdate time.Time `json:"birthdate"`
}

func (t *CustomDate) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	t.Birthdate = date
	return
}
