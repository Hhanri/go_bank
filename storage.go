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

func (ps *PostgresStore) Init() error {
	return ps.createAccountTable()
}

func (ps *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		ID serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at TIMESTAMP
	)`

	_, err := ps.db.Exec(query)
	return err
}

func (ps *PostgresStore) CreateAcccount(account *Account) error {
	query := `insert into account (first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)`

	resp, err := ps.db.Query(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

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
