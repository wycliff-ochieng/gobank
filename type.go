package main

import (
	"math/rand"
)

type Account struct {
	ID        int
	Firstname string
	Lastname  string
	Number    int
	Balance   int64
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(10),
		Firstname: firstname,
		Lastname:  lastname,
		Number:    rand.Intn(100),
	}
}
