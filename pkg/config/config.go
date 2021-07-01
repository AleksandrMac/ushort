package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/AleksandrMac/ushort/pkg/models"
	"github.com/go-chi/jwtauth/v5"
)

type Server struct {
	Port    string `yaml:"Port" env:"PORT" envDefault:"8000"`
	TimeOut int64  `yaml:"TimeOut" env:"SERVER_TIMEOUT" envDefault:"30"`
}

type DB struct {
	URL string `yaml:"db_url" env:"DATABASE_URL" envDefault:"postgres://postgres:password@localhost:5432/ushort?sslmode=disable"`

	// Driver   string `yaml:"driver" env:"DB_DRIVER" envDefault:"postgres"`
	// Name     string `yaml:"name" env:"DB_NAME" envDefault:"ushort"`
	// Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	// Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
	// User     string `yaml:"user" env:"DB_USER" envDefault:"postgres"`
	// Password string `yaml:"password" env:"DB_PASSWORD,unset"`
	// SslMode  string `yaml:"sslMode" env:"DB_SSLMODE" envDefault:"disable"`
	// TimeZone string `yaml:"timeZone" env:"DB_TIMEZONE" envDefault:"Europe/Moscow"`
}

type Auth struct {
	PrivateKey []byte
	PublicKey  []byte
	ExpiresAt  int64 `yaml:"ExpiresAt" env:"AUTH_EXPIRESAT" envDefault:"30"` // minut
}

type Config struct {
	Server Server `yaml:"Server"`
	DB     DB     `yaml:"DB"`
	Auth   Auth   `yaml:"Auth"`
	// TmpURLLifeTime - время резервации сгенерированной ссылки
	TmpURLLifeTime int64 `yaml:"TmpURLLifeTime" env:"TMPURL_LIFE_TIME" envDefault:"60"` // second
	LengthURL      int64 `yaml:"LengthURL" env:"LENGTH_URL" envDefault:"10"`            // second
}

type Env struct {
	DB        *models.DB
	Config    *Config
	TokenAuth *jwtauth.JWTAuth
}

func New() (*Config, error) {
	cfg := new(Config)
	err := cfg.installKeys()
	return cfg, err
}

func (cfg *Config) installKeys() error {
	dir := filepath.Join("internal", "secrets")
	if err := cfg.Private(filepath.Join(dir, "myprivatekey.pem")); err != nil {
		return err
	}
	return cfg.Public(filepath.Join(dir, "mypublickey.pem"))
}

func (cfg *Config) Private(filename string) (err error) {
	log.Default().Printf("Читаю файл: %s", filename)
	cfg.Auth.PrivateKey, err = ioutil.ReadFile(filename)
	log.Default().Printf("Cчитан ivate key: длина %d байт", len(cfg.Auth.PrivateKey))
	return
}

func (cfg *Config) Public(filename string) (err error) {
	log.Default().Printf("Читаю файл: %s", filename)
	cfg.Auth.PublicKey, err = ioutil.ReadFile(filename)
	log.Default().Printf("Cчитан Public key: длина %d байт", len(cfg.Auth.PublicKey))
	return
}
