package data

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccount(int) error
	DeleteAccount(*Account) error
	UpdateAccount(int) error
}

type Postgrestore struct {
	db *sql.DB
}

func NewPostgrestore() (*Postgrestore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgrestore{
		db: db,
	}, nil
}

func (p *Postgrestore) Init() error {
	return p.CreateAccountTable()
}

func (p *Postgrestore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts(
	id SERIAL PRIMARY KEY,
	firstname VARCHAR(10),
	lastname VARCHAR(10),
	balance SERIAL,
	createdat TIMESTAMP)`

	_, err := p.db.Exec(query)
	return err
}

func (p *Postgrestore) CreateTransactionTable() error {
	query := `CREATE TABLE IF NOT EXISTs transactions(
	id SERIAL PRIMARY KEY,
	amount_id INT, 
	amount INT,
	type VARCHAR(10),
	transactionTime TIMESTAMP DEFAULT CURRENT TIMESTAMP)`

	_, err := p.db.Exec(query)
	return err
}

func (p *Postgrestore) CreateAccount(account *Account) error {
	_, err := p.db.Exec(
		`INSERT INTO accounts(firstname, lastname, balance, createdat) VALUES($1, $2, $3, $4)`,
		account.Firstname, account.Lastname, account.Balance, account.CreatedAt,
	)
	if err != nil {
		return errors.New("cannot insert into table: " + err.Error())
	}
	return nil
}

func (p *Postgrestore) GetAccount(int) error {
	acc := &Account{}
	query := `SELECT id,firstname,lastname,balance,createdat`
	err := p.db.QueryRow(query).Scan(&acc.ID, &acc.Firstname, &acc.Lastname, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgrestore) DeleteAccount(*Account) error {
	return nil
}

func (p *Postgrestore) UpdateAccount(int) error {
	return nil
}
