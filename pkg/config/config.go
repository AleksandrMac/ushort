package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Port    string `yaml:"Port"`
	TimeOut int64  `yaml:"TimeOut"`
}

func (s Server) String() string {
	return fmt.Sprintf("\n  Port: %s\n  TimeOut: %ds", s.Port, s.TimeOut)
}

type DB struct {
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslMode"`
	TimeZone string `yaml:"timeZone"`
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

func New(path string) (*Config, error) {
	config := new(Config)
	configFileName := path
	if buf, err := ioutil.ReadFile(configFileName); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	return config, nil
}
