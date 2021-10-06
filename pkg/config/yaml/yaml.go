package yaml

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/AleksandrMac/ushort/pkg/config"
)

func New(path string) (*config.Config, error) {
	conf := new(config.Config)
	configFileName := path
	if buf, err := ioutil.ReadFile(configFileName); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}
