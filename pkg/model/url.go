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

func (u *URL) Values() (map[DBField]interface{}, error) {
	fields, err := u.Fields()
	if err != nil {
		return nil, err
	}
	out := make(map[DBField]interface{}, len(fields))
	for _, val := range fields {
		out[DBField(val)] = u.Value(DBField(val))
	}
	return out, nil
}

func (u *URL) Value(field DBField) interface{} {
	return reflect.ValueOf(u).FieldByName(string(field)).Interface()
}

func (u *URL) SetValues(mapValues map[string]interface{}) error {
	return utils.UpdateStruct(u, mapValues)
}

func (u *URL) SetValue(field DBField, value interface{}) error {
	return u.SetValues(map[string]interface{}{string(field): value})
}
func (u *URL) JSON() ([]byte, error) {
	return json.Marshal(u)
}

func (u *URL) FromJSON(body []byte) error {
	return json.Unmarshal(body, u)
}

func (u *URL) Create() error {
	_, err := u.NamedExec(`INSERT INTO "public"."url" (id,redirect_to,description,user_id)
 	VALUES (:id, :redirect_to, :description, :user_id);`, u)
	return err
}

func (u *URL) Read() error {
	return u.Get(u, `SELECT * FROM public.url WHERE id=$1;`, u.ID)
}

func (u *URL) ReadAll(userID string) ([]*URL, error) {
	list := []*URL{}
	if userID == "" {
		if err := u.Select(list, `SELECT * FROM public.url`); err != nil {
			return nil, err
		}
	} else {
		if err := u.Select(list, `SELECT * FROM public.url where user_id = $1`, userID); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (u *URL) Update() error {
	_, err := u.NamedExec(`UPDATE public.url
    SET redirect_to=:redirect_to,
		description=:description
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}

func (u *URL) Delete() error {
	_, err := u.NamedExec(`DELETE FROM public.url
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}
