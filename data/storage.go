package data

import (
	"database/sql"
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
	connStr := "user=postgres,dname=postgres,password=gobank,sslmode=disable"
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

func (p *Postgrestore) CreateAccount(acc *Account) error {
	query := p.db.Query(
		`INSERT INTO accounts(id,firstname,lastname,balance,createdat) VALUES(&1,&2,&3,&4,&5)`
	)
	resp,err:= p.db.Exec(query, acc.Firstname,acc.Lastname,acc.Balance,acc.Createdat)
}

func (p *Postgrestore) GetAccount(int) error {
	return nil
}

func (p *Postgrestore) DeleteAccount(*Account) error {
	return nil
}

func (p *Postgrestore) UpdateAccount(int) error {
	return nil
}
