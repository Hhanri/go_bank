package main

import "math/rand"

type Account struct {
	ID        int
	FirstName string
	LastName  string
	Number    int64
	Balance   int64
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Int(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Int()),
	}
}
