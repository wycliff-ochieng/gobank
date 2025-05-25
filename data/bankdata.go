package data

import (
	"math/rand"
	"time"
)

type CreateAccountReq struct {
	Firstname string
	Lastname  string
}

type TransferReq struct {
	Ammount int `json:"ammount"`
}

type Account struct {
	ID        int
	Firstname string
	Lastname  string
	Balance   int
	CreatedAt time.Time
}

type Transaction struct {
	ID              int
	AccountID       int
	Type            string //deposit / withdrawal
	Transactiontime string
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
