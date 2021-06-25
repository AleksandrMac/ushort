package config

import (
	"fmt"
)

type Server struct {
	Port    string `yaml:"Port" env:"SERVER_PORT" envDefault:"8000"`
	TimeOut int64  `yaml:"TimeOut" env:"SERVER_TIMEOUT" envDefault:"30"`
}

func (s Server) String() string {
	return fmt.Sprintf("\n  Port: %s\n  TimeOut: %ds", s.Port, s.TimeOut)
}

type DB struct {
	Driver   string `yaml:"driver" env:"DB_DRIVER" envDefault:"postgres"`
	Name     string `yaml:"name" env:"DB_NAME" envDefault:"ushort"`
	Host     string `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" envDefault:"5432"`
	User     string `yaml:"user" env:"DB_USER" envDefault:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD,unset"`
	SslMode  string `yaml:"sslMode" env:"DB_SSLMODE" envDefault:"disable"`
	TimeZone string `yaml:"timeZone" env:"DB_TIMEZONE" envDefault:"Europe/Moscow"`
}

func (db DB) String() string {
	return fmt.Sprintf(`
  Driver: %v
  Name: %v
  Host: %v
  Port: %v
  User: %v
  Password: %v
  SslMode: %v
  TimeZone: %v`, db.Driver, db.Name, db.Host, db.Port, db.User, db.Password, db.SslMode, db.TimeZone)
}

type Config struct {
	Server Server `yaml:"Server"`
	DB     DB     `yaml:"DB"`
}

func (c Config) String() string {
	return fmt.Sprintf("Server: %v\nDB: %+v", c.Server, c.DB)
}
