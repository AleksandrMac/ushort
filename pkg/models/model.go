package models

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DataStore interface {
	Insert() (uuid.UUID, error)
}

type DB struct {
	*sqlx.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// // DataSourceParam формирует строку с параметрами
// func DataSourceParam(db *config.DB) string {
// 	param := "?"
// 	if len(db.SslMode) > 0 {
// 		param += "sslmode=" + db.SslMode + "&"
// 	} else {
// 		param += "sslmode=disable&"
// 	}
// 	if len(db.TimeZone) > 0 {
// 		param += "timezone=" + db.TimeZone + "&"
// 	}
// 	return param[:len(param)-1]
// }
