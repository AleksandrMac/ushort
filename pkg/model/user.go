package model

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/AleksandrMac/ushort/pkg/utils"
	"github.com/jmoiron/sqlx"

	// Регистрация диалекта БД
	_ "github.com/lib/pq"
)

type User struct {
	*sqlx.DB
	Base
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

func (u *User) Fields() ([]string, error) {
	return utils.FieldsFromStruct(u)
}

func (u *User) Values() (map[Field]interface{}, error) {
	fields, err := u.Fields()
	if err != nil {
		return nil, err
	}
	out := make(map[Field]interface{}, len(fields))
	for _, val := range fields {
		out[Field(val)], err = u.Value(Field(val))
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (u *User) Value(field Field) (interface{}, error) {
	return reflect.ValueOf(u).Elem().FieldByName(string(field)).Interface(), nil
}

func (u *User) SetValues(mapValues map[string]interface{}) error {
	return utils.UpdateStruct(u, mapValues)
}
func (u *User) SetValue(field Field, value interface{}) error {
	return u.SetValues(map[string]interface{}{string(field): value})
}
func (u *User) JSON() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) FromJSON(body []byte) error {
	return json.Unmarshal(body, u)
}

func (u *User) Create() error {
	_, err := u.NamedExec(
		`INSERT INTO public.user (id, email,password) VALUES (:id, :email, :password);`, u)
	return err
}

func (u *User) Read() error {
	if u.ID != "" {
		err := u.Get(u, `SELECT * FROM public.user WHERE id=$1;`, u.ID)
		return err
	}

	if u.Email != "" {
		err := u.Get(u, `SELECT * FROM public.user WHERE email=$1;`, u.Email)
		return err
	}
	return fmt.Errorf("missing field 'id' or 'email' in %T", u)
}

func (u *User) ReadAll(userID string) ([]*User, error) {
	list := []*User{}
	if err := u.Select(list, `SELECT * FROM public.user`); err != nil {
		return nil, err
	}
	return list, nil
}

func (u *User) Update() error {
	_, err := u.NamedExec(`UPDATE public.user
SET email=:email,
password=:password
WHERE id=:id;`, u)
	return err
}

func (u *User) Delete() error {
	if u.ID != "" {
		err := u.Get(u, `DELETE FROM public.user WHERE id=$1;`, u.ID)
		return err
	}

	if u.Email != "" {
		err := u.Get(u, `DELETE FROM public.user WHERE email=$1;`, u.Email)
		return err
	}
	return fmt.Errorf("missing field 'id' or 'email' in %T", u)
}
