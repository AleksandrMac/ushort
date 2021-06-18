package models

import "github.com/AleksandrMac/ushort/pkg/config"

// DataSourceParam формирует строку с параметрами
func DataSourceParam(db *config.DB) string {
	param := "?"
	if len(db.SslMode) > 0 {
		param += "sslmode=" + db.SslMode + "&"
	} else {
		param += "sslmode=disable&"
	}
	if len(db.TimeZone) > 0 {
		param += "timezone=" + db.TimeZone + "&"
	}
	return param[:len(param)-1]
}
