package models

type Account struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	CurrentAmount int           `json:"currentAmount"`
	Transactions  []Transaction `json:"transactions"`
}

type CreateAccountBody struct {
	Name string `json:"name"`
}

type CurrentAmount struct {
	Amount int `json:"amount"`
}
