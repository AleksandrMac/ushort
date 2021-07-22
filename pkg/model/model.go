package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"

	// используется по для чтения миграций из файла
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	ReadAll(table Table, userID string) ([]Model, error)
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

	log.Default().Printf("БД подключена")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	log.Default().Printf("Начинаем миграции")
	// nolint: gomnd	// меньше нуля migration.down, иначе migration.up
	err = m.Steps(2)
	if err != nil {
		switch err {
		case os.ErrNotExist:
			log.Default().Printf("Новых миграций не найдено")
		default:
			return nil, err
		}
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		db,
		&User{},
		&URL{},
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

func (db *DB) ReadAll(table Table, userID string) ([]Model, error) {
	var out []Model
	switch table {
	case TableUser:
		users, err := db.user.ReadAll(userID)
		if err != nil {
			return nil, err
		}
		out = make([]Model, 0, len(users))
		for _, val := range users {
			out = append(out, val)
		}
	case TableURL:
		urls, err := db.url.ReadAll(userID)
		if err != nil {
			return nil, err
		}
		out = make([]Model, 0, len(urls))
		for _, val := range urls {
			out = append(out, val)
		}
	default:
		return nil, fmt.Errorf("'%T' table not found", table)
	}
	return out, nil
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
