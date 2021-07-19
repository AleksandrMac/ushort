module github.com/AleksandrMac/ushort

go 1.16

require (
	github.com/caarlos0/env/v6 v6.6.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/httplog v0.2.0
	github.com/go-chi/jwtauth/v5 v5.0.1
	github.com/google/uuid v1.2.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/rs/zerolog v1.18.1-0.20200514152719-663cbb4c8469
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// +heroku goVersion go1.16 install ./cmd/ushort ./bin/ushort
