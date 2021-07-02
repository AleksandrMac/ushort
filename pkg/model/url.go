package model

import (
	"github.com/AleksandrMac/ushort/pkg/connect"
)

type URL struct {
	Model
	RedirectTo  string `db:"redirect_to" json:"redirectTo"`
	Description string `db:"description" json:"description"`
	UserID      string `db:"user_id"`
}

func (u *URL) Insert(db *connect.DB) error {
	_, err := db.NamedExec(`INSERT INTO "public"."url" (id,redirect_to,description,user_id)
	VALUES (:id, :redirect_to, :description, :user_id);`, u)
	if err != nil {
		return err
	}
	return nil
}

func (u *URL) Update(db *connect.DB) (err error) {
	_, err = db.NamedExec(`UPDATE public.url
    SET redirect_to=:redirect_to,
		description=:description
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}

func (u *URL) Delete(db *connect.DB) (err error) {
	_, err = db.NamedExec(`DELETE FROM public.url
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}

func (u *URL) Select(db *connect.DB) (*[]URL, error) {
	urls := &[]URL{}
	err := db.Select(urls, `SELECT * FROM url WHERE user_id=$1`, u.ID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (u *URL) SelectWithID(db *connect.DB) (url *URL, err error) {
	url = &URL{}
	switch u.ID {
	case "*":
		err = db.Get(url, `SELECT * FROM url WHERE id=$1;`, u.ID)
	default:
		err = db.Get(url, `SELECT * FROM url WHERE id=$1 AND user_id=$2;`, u.ID, u.UserID)
	}
	if err != nil {
		return nil, err
	}
	return url, nil
}
