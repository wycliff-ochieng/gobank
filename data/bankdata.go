package data

import "math/rand"

type Account struct {
	ID        int
	Firstname string
	Lastname  string
	Balance   int
}

func NewAccount(firstname string, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(10),
		Firstname: firstname,
		Lastname:  lastname,
		Balance:   rand.Intn(1000),
	}
}
