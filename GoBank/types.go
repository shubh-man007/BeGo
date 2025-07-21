package main

import "math/rand/v2"

type account struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	AccountID int64   `json:"accountID"`
	Number    int64   `json:"number"`
	Balance   float64 `json:"balance"`
}

func NewAccount(FirstName, LastName string) *account {
	return &account{
		FirstName: FirstName,
		LastName:  LastName,
		AccountID: rand.Int64(),
		Number:    rand.Int64(),
		Balance:   rand.Float64(),
	}
}
