package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAcccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
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

func (ps *PostgresStore) CreateAcccount(account *Account) (*Account, error) {
	query := `insert into account (first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)
	returning *`

	resp, err := ps.db.Query(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", resp)

	for resp.Next() {
		return scanIntoAccount(resp)
	}

	return nil, fmt.Errorf("bad sql row")
}

func (ps *PostgresStore) DeleteAccount(id int) error {
	_, err := ps.db.Query("delete from account where id = $1", id)

	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (ps *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := ps.db.Query("select * from account where ID = $1", id)

	if err != nil {
		fmt.Println("Error in db query")
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (ps *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := ps.db.Query("select * from account")

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	return account, err
}
