package models

import "time"

type Transaction struct {
	ID   int       `json:"id"`
	Date time.Time `json:"date"`

	Amount int    `json:"amount"`
	Note   string `json:"note"`

	Type string `json:"type"`
}
