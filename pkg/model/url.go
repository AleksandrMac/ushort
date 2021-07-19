package model

import (
	"encoding/json"
	"reflect"

	"github.com/AleksandrMac/ushort/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type URL struct {
	*sqlx.DB
	Base
	RedirectTo  string `db:"redirect_to" json:"redirectTo"`
	Description string `db:"description" json:"description"`
	UserID      string `db:"user_id"`
}

func (u *URL) Fields() ([]string, error) {
	return utils.FieldsFromStruct(u)
}

func (u *URL) Values() (map[string]interface{}, error) {
	fields, err := u.Fields()
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{}, len(fields))
	for _, val := range fields {
		out[val] = u.Value(val)
	}
	return out, nil
}

func (u *URL) Value(field string) interface{} {
	return reflect.ValueOf(u).FieldByName(field).Interface()
}

func (u *URL) SetValues(mapValues map[string]interface{}) error {
	return utils.UpdateStruct(u, mapValues)
}

func (u *URL) SetValue(field string, value interface{}) error {
	return u.SetValues(map[string]interface{}{field: value})
}
func (u *URL) JSON() ([]byte, error) {
	return json.Marshal(u)
}

func (u *URL) create() error {
	_, err := u.NamedExec(`INSERT INTO "public"."url" (id,redirect_to,description,user_id)
 	VALUES (:id, :redirect_to, :description, :user_id);`, u)
	return err
}

func (u *URL) read() error {
	return u.Get(u, `SELECT * FROM url WHERE id=$1;`, u.ID)
}

func (u *URL) update() error {
	_, err := u.NamedExec(`UPDATE public.url
    SET redirect_to=:redirect_to,
		description=:description
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}

func (u *URL) delete() error {
	_, err := u.NamedExec(`DELETE FROM public.url
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}
