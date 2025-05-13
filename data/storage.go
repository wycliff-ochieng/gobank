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
	query := ""
	_, err := p.db.Exec(query)
	return err
}

func (p *Postgrestore) CreateAccount(*Account) error {
	return nil
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
