package env

import (
	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/caarlos0/env/v6"
)

func New() (*config.Config, error) {
	conf := new(config.Config)
	if err := env.Parse(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
