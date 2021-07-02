package model

import (
	"github.com/AleksandrMac/ushort/pkg/connect"

	// Регистрация диалекта БД
	_ "github.com/lib/pq"
)

type User struct {
	Model
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

func New() User {
	return User{}
}

func (u *User) Insert(db *connect.DB) error {
	_, err := db.NamedExec(
		`INSERT INTO "public"."users" ("id","email","password") VALUES (:id, :email, :password);`, u)
	return err
}

// func Select(c *config.DB) (*[]User, error) {
// 	db, err := sqlx.Connect("postgres", c.URL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err = db.Close(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	users := []User{}
// 	err = db.Select(&users, `SELECT * FROM users;`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &users, nil
// }

// func (u *User) SelectWithID(c *config.DB) (*User, error) {
// 	db, err := sqlx.Connect("postgres", c.URL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err = db.Close(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	user := User{}
// 	err = db.Get(&user, `SELECT * FROM users WHERE id=$1;`, u.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }
func SelectWithEmail(email string, db *connect.DB) (*User, error) {
	user := new(User)
	err := db.DB.Get(user, `SELECT * FROM users WHERE email=$1;`, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func (u *User) Delete(c *config.DB) error {
// 	db, err := sqlx.Connect("postgres", c.URL)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err = db.Close(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	_, err = db.Exec(`DELETE FROM user WHERE id=$1;`, u.ID.String())
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
