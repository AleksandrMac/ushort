package model

import "time"

type Model struct {
	ID        string    `db:"id" json:"urlID"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
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
