package models

import "time"

type Transaction struct {
	ID      int       `json:"id"`
	Date    time.Time `json:"date"`
	Account string    `json:"account"`

	Amount int `json:"amount"`
	Note   int `json:"note"`
}
