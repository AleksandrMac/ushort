module github.com/AleksandrMac/ushort

go 1.17

require (
	github.com/caarlos0/env/v6 v6.6.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/httplog v0.2.0
	github.com/go-chi/jwtauth/v5 v5.0.1
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v3 v3.0.0 // indirect
	github.com/goccy/go-json v0.7.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/jwx v1.2.4 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/lestrrat-go/pdebug/v3 v3.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// +heroku goVersion go1.17 install ./cmd/ushort ./bin/ushort
