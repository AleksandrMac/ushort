package url

import (
	"time"

	"github.com/AleksandrMac/ushort/pkg/models"
)

type URL struct {
	ID          string    `db:"id" json:"urlID"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	RedirectTo  string    `db:"redirect_to" json:"redirectTo"`
	Description string    `db:"description" json:"description"`
	UserID      string    `db:"user_id"`
}

func (u *URL) Insert(db *models.DB) error {
	_, err := db.NamedExec(`INSERT INTO "public"."url" (id,redirect_to,description,user_id)
	VALUES (:id, :redirect_to, :description, :user_id);`, u)
	if err != nil {
		return err
	}
	return nil
}

func (u *URL) Update(db *models.DB) (err error) {
	_, err = db.NamedExec(`UPDATE public.url
    SET redirect_to=:redirect_to,
		description=:description
    WHERE id=:id AND user_id=:user_id;`, u)
	return err
}

func Select(userID string, db *models.DB) (*[]URL, error) {
	urls := &[]URL{}
	err := db.Select(urls, `SELECT * FROM url WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func SelectWithID(urlID, userID string, db *models.DB) (url *URL, err error) {
	url = &URL{}
	switch userID {
	case "*":
		err = db.Get(url, `SELECT * FROM url WHERE id=$1;`, urlID)
	default:
		err = db.Get(url, `SELECT * FROM url WHERE id=$1 AND user_id=$2;`, urlID, userID)
	}
	if err != nil {
		return nil, err
	}
	return url, nil
}
