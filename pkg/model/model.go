package model

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Table uint8

const (
	UserTable Table = iota
	URLTable
)

type Model interface {
	Fields() ([]string, error)
	Values() (map[string]interface{}, error)
	Value(field string) interface{}
	SetValue(field string, val interface{}) error
	JSON() ([]byte, error)
}

type CRUD interface {
	Model(table Table) Model
	Create(table Table) error
	Read(table Table) error
	Update(table Table) error
	Delete(table Table) error
}

type DB struct {
	*sqlx.DB
	user *User
	url  *URL
}

type Base struct {
	ID        string    `db:"id" json:"urlID"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		db,
		&User{db, Base{}, "", ""},
		&URL{db, Base{}, "", "", ""},
	}, nil
}

func (db *DB) Model(table Table) Model {
	switch table {
	case UserTable:
		return db.user
	case URLTable:
		return db.url
	default:
		return nil
	}
}

func (db *DB) Create(table Table) error {
	switch table {
	case UserTable:
		return db.user.create()
	case URLTable:
		return db.url.create()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}

func (db *DB) Read(table Table) error {
	switch table {
	case UserTable:
		return db.user.read()
	case URLTable:
		return db.url.read()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}

func (db *DB) Update(table Table) error {
	switch table {
	case UserTable:
		return db.user.update()
	case URLTable:
		return db.url.update()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}

func (db *DB) Delete(table Table) error {
	switch table {
	case UserTable:
		return db.user.delete()
	case URLTable:
		return db.url.delete()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}
