package data

import (
	"math/rand"
	"time"
)

type CreateAccountReq struct {
	Firstname string
	Lastname  string
}

type Account struct {
	ID        int
	Firstname string
	Lastname  string
	Balance   int
	CreatedAt time.Time
}

func NewAccount(firstname string, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(10),
		Firstname: firstname,
		Lastname:  lastname,
		Balance:   rand.Intn(1000),
		CreatedAt: time.Now().UTC(),
	}
}
