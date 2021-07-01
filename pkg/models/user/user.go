package user

import (
	"log"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	// Регистрация диалекта БД
	_ "github.com/lib/pq"
)

type User struct {
	ID       uuid.UUID `db:"id" json:"uuid"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"password"`
}

func New() User {
	return User{}
}

// func CreateTable(c *config.DB) error {
// 	db, err := sqlx.Connect("postgres", c.URL)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err = db.Close(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	// если таблица отсутствует, то создаем новую, иначе добавляем новые столбцы, либо обновлям
// 	rows, err := db.Query(`SELECT * FROM INFORMATION_SCHEMA.TABLES
// WHERE TABLE_SCHEMA = 'public'
// AND  TABLE_NAME = 'users'`)
// 	if err != nil {
// 		return err
// 	}
// 	if rows.Next() {
// 		if _, err = db.Exec(`ALTER TABLE users
// ADD COLUMN IF NOT EXISTS id uuid CONSTRAINT user_id PRIMARY KEY,
// ADD COLUMN IF NOT EXISTS email text CONSTRAINT email UNIQUE,
// ADD COLUMN IF NOT EXISTS password text;
// ALTER TABLE users
// ALTER COLUMN email SET DATA TYPE text,
// ALTER COLUMN email SET NOT NULL,
// ALTER COLUMN password SET DATA TYPE text,
// ALTER COLUMN password SET NOT NULL;`); err != nil {
// 			return err
// 		}
// 	} else {
// 		if _, err = db.Exec(`
// CREATE TABLE users (
// id uuid CONSTRAINT user_id PRIMARY KEY,
// email text NOT NULL,
// password text NOT NULL);`); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (u *User) Insert(db *models.DB) (uuid.UUID, error) {
	u.ID = uuid.New()
	_, err := db.NamedExec(`INSERT INTO "public"."users" ("id","email","password")
	VALUES (:id, :email, :password);`, u)
	if err != nil {
		return uuid.UUID{}, err
	}
	return u.ID, nil
}

func Select(c *config.DB) (*[]User, error) {
	db, err := sqlx.Connect("postgres", c.URL)
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
	db, err := sqlx.Connect("postgres", c.URL)
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
func SelectWithEmail(email string, db *models.DB) (*User, error) {
	user := new(User)
	err := db.DB.Get(user, `SELECT * FROM users WHERE email=$1;`, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Delete(c *config.DB) error {
	db, err := sqlx.Connect("postgres", c.URL)
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
