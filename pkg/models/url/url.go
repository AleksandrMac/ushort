package url

import (
	"time"

	"github.com/AleksandrMac/ushort/pkg/models"

	"github.com/google/uuid"
)

type URL struct {
	ID          string    `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updates_at"`
	RedirectTo  string    `db:"redirect_to" json:"redirect_to"`
	Description string    `db:"description" json:"description"`
	UserID      uuid.UUID `db:"user_id"`
}

func SelectWithID(id string, db *models.DB) (*URL, error) {
	url := &URL{}
	err := db.Get(url, `SELECT * FROM url WHERE id=$1;`, id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
