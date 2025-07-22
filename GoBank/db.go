package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) CreateAccountsTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		account_id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		number BIGINT NOT NULL,
		balance DOUBLE PRECISION NOT NULL,
		created TIMESTAMP DEFAULT NOW()
	)`
	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStorage) CreateAccount(a *account) error {
	query := `INSERT INTO accounts (first_name, last_name, number, balance)
			  VALUES ($1, $2, $3, $4) RETURNING account_id`

	return p.db.QueryRow(query, a.FirstName, a.LastName, a.Number, a.Balance).Scan(&a.AccountID)
}

func (p *PostgresStorage) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE account_id = $1`

	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) UpdateAccount(a *account) error {
	query := `UPDATE accounts SET first_name = $1,
								  last_name = $2,
								  number = $3,
								  balance = $4
			WHERE account_id = $5`

	_, err := p.db.Exec(query, a.FirstName, a.LastName, a.Number, a.Balance, a.AccountID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) GetAccountByID(id int) (*account, error) {
	query := `SELECT account_id, first_name, last_name, number, balance FROM accounts WHERE account_id = $1`

	var a account
	err := p.db.QueryRow(query, id).Scan(&a.AccountID, &a.FirstName, &a.LastName, &a.Number, &a.Balance)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
