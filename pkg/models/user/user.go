package user

import (
	"fmt"
	"log"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/models"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	// Регистрация диалекта БД
	_ "github.com/lib/pq"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
}

func New() User {
	return User{}
}

func CreateTable(c *config.DB) error {
	// nolint:lll // ложное срабатывание из за большого количества параметров при формировании url
	db, err := sqlx.Connect(c.Driver, fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", c.User, c.Password, c.Host, c.Port, c.Name, models.DataSourceParam(c)))
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	// если таблица отсутствует, то создаем новую, иначе добавляем новые столбцы, либо обновлям
	rows, err := db.Query(`SELECT * FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'public'
AND  TABLE_NAME = 'users'`)
	if err != nil {
		return err
	}
	if rows.Next() {
		if _, err = db.Exec(`ALTER TABLE users
ADD COLUMN IF NOT EXISTS id uuid CONSTRAINT user_id PRIMARY KEY,
ADD COLUMN IF NOT EXISTS email text CONSTRAINT email UNIQUE,
ADD COLUMN IF NOT EXISTS password text;
ALTER TABLE users
ALTER COLUMN email SET DATA TYPE text,
ALTER COLUMN email SET NOT NULL,
ALTER COLUMN password SET DATA TYPE text,
ALTER COLUMN password SET NOT NULL;`); err != nil {
			return err
		}
	} else {
		if _, err = db.Exec(`
CREATE TABLE users (
id uuid CONSTRAINT user_id PRIMARY KEY,
email text NOT NULL,
password text NOT NULL);`); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Insert(c *config.DB) (uuid.UUID, error) {
	// nolint:lll // ложное срабатывание из за большого количества параметров при формировании url
	db, err := sqlx.Connect(c.Driver, fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", c.User, c.Password, c.Host, c.Port, c.Name, models.DataSourceParam(c)))
	if err != nil {
		return uuid.UUID{}, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	u.ID = uuid.New()
	_, err = db.NamedExec(`INSERT INTO "public"."users" ("id","email","password")
	VALUES (:id, :email, :password);`, u)
	if err != nil {
		return uuid.UUID{}, err
	}
	return u.ID, nil
}

func Select(c *config.DB) (*[]User, error) {
	// nolint:lll // ложное срабатывание из за большого количества параметров при формировании url
	db, err := sqlx.Connect(c.Driver, fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", c.User, c.Password, c.Host, c.Port, c.Name, models.DataSourceParam(c)))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	users := []User{}
	err = db.Select(&users, `SELECT * FROM users;`)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (u *User) SelectWithID(c *config.DB) (*User, error) {
	// nolint:lll // ложное срабатывание из за большого количества параметров при формировании url
	db, err := sqlx.Connect(c.Driver, fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", c.User, c.Password, c.Host, c.Port, c.Name, models.DataSourceParam(c)))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	user := User{}
	err = db.Get(&user, `SELECT * FROM users WHERE id=$1;`, u.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Delete(c *config.DB) error {
	// nolint:lll // ложное срабатывание из за большого количества параметров при формировании url
	db, err := sqlx.Connect(c.Driver, fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s", c.User, c.Password, c.Host, c.Port, c.Name, models.DataSourceParam(c)))
	if err != nil {
		return err
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	_, err = db.Exec(`DELETE FROM user WHERE id=$1;`, u.ID.String())
	if err != nil {
		return err
	}
	return nil
}
