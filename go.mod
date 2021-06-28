module github.com/AleksandrMac/ushort

go 1.16

require (
	github.com/caarlos0/env/v6 v6.6.2
	github.com/go-chi/chi/v5 v5.0.3
	github.com/google/uuid v1.2.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	gopkg.in/yaml.v2 v2.4.0
)

// +heroku goVersion go1.16 install ./cmd/ushort /bin/ushort