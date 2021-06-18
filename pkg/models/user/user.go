package user

import (
	"fmt"
	"log"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	id       uint64
	email    string
	password string
}

var shemas = map[string]string{
	"create": `CREATE TABLE user (
	id uuid,
	email var,
	password varchar);`,
	"insert": `INSERT INTO user
	(id, email, password)
	VALUES ($1, $2, $3);`,
	"update": `UPDATE user
	SET  email = value1, column2 = value2, ...
	WHERE id = ;`,
	"delete": ``,
	"select": ``,
}

func CreateTable(c *config.DB) error {
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
