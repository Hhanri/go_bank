package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Int(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Int()),
		CreatedAt: time.Now().UTC(),
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
