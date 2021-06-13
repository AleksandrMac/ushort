package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Port    int64 `yaml:"port"`
	TimeOut int64 `yaml:"timeOut"`
}

func (s Server) String() string {
	return fmt.Sprintf("\n  Port: %d\n  TimeOut: %ds\n", s.Port, s.TimeOut)
}

type Config struct {
	Server Server `yaml:"server"`
}

func (c Config) String() string {
	return fmt.Sprintf("Server: %v", c.Server)
}

func New(path string) (*Config, error) {
	config := &Config{}
	configFileName := filepath.Join(path)
	if buf, err := ioutil.ReadFile(configFileName); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	return config, nil
}
