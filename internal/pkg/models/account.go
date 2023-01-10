package models

import "github.com/uptrace/bun"

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:ac"`

	UserID    string `bun:",autoincrement"`
	Username  string `bun:"username,notnull"`
	Password  string `bun:"password,notnull"`
	Firstname string
	Lastname  string
	Email     string
	Amka      string
	Age       string
	Job       string
	CreatedOn string
	LastLogin string
}
