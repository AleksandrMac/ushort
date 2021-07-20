package model

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Table uint8

const (
	TableUser Table = iota
	TableURL
)

type DBField string

const (
	DBFieldID          DBField = "id"
	DBFieldCreatedAt   DBField = "created_at"
	DBFieldUpdateddAt  DBField = "updated_at"
	DBFieldEmail       DBField = "email"
	DBFieldPassword    DBField = "password"
	DBFieldRedirectTo  DBField = "redirect_to"
	DBFieldDescription DBField = "description"
	DBFieldUserID      DBField = "user_id"
)

type Model interface {
	Fields() ([]string, error)
	Values() (map[DBField]interface{}, error)
	Value(field DBField) interface{}
	SetValue(field DBField, val interface{}) error
	JSON() ([]byte, error)
	FromJSON(body []byte) error
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
		nil,
		nil,
	}, nil
}

func (db *DB) Model(table Table) Model {
	switch table {
	case TableUser:
		return db.user
	case TableURL:
		return db.url
	default:
		return nil
	}
}

func (db *DB) Create(table Table) error {
	switch table {
	case TableUser:
		return db.user.Create()
	case TableURL:
		return db.url.Create()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}

func (db *DB) Read(table Table) error {
	switch table {
	case TableUser:
		return db.user.Read()
	case TableURL:
		return db.url.Read()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}

func (db *DB) Update(table Table) error {
	switch table {
	case TableUser:
		return db.user.Update()
	case TableURL:
		return db.url.Update()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}

func (db *DB) Delete(table Table) error {
	switch table {
	case TableUser:
		return db.user.Delete()
	case TableURL:
		return db.url.Delete()
	default:
		return fmt.Errorf("'%T' table not found", table)
	}
}
