package models

import "github.com/uptrace/bun"

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:ac"`

	UserID    string `bun:",autoincrement" `
	Username  string `bun:"username, notnull" json:"username" validate:"required"`
	Password  string `bun:"password, notnull" json:"password" validate:"required"`
	Firstname string `bun:"firstname, notnull" json:"firstname" validate:"required"`
	Lastname  string `bun:"lastname, notnull" json:"lastname" validate:"required"`
	Email     string `bun:"email, notnull" json:"email" validate:"required,email,max=128"`
	Amka      string `bun:"amka, notnull" json:"amka" validate:"required"`
	Age       string `bun:"age, notnull" json:"age" validate:"required"`
	Job       string
	CreatedOn string
	LastLogin string
}
