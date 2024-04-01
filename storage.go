package main

import "database/sql"

type Storage interface {
	CreateAcccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db}, nil
}

func (ps *PostgresStore) CreateAcccount(account *Account) error {
	return nil
}

func (ps *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (ps *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (ps *PostgresStore) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
