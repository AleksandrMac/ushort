package url

import (
	"time"

	"github.com/AleksandrMac/ushort/pkg/models"
)

type URL struct {
	ID          string    `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updates_at"`
	RedirectTo  string    `db:"redirect_to" json:"redirect_to"`
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

func Select(userID string, db *models.DB) (*[]URL, error) {
	urls := &[]URL{}
	err := db.Select(urls, `SELECT * FROM url WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func SelectWithID(id string, db *models.DB) (*URL, error) {
	url := &URL{}
	err := db.Get(url, `SELECT * FROM url WHERE id=$1;`, id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
