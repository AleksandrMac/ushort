package env

import (
	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/caarlos0/env/v6"
)

func New() (*config.Config, error) {
	config := new(config.Config)
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
