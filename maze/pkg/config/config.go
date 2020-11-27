package config

import (
	"encoding/json"
	"github.com/avanticaTest/maze/pkg/consul"
)

type Config struct {
	KeyConfig string `json:"-"`
	DBName    string `json:"db_name"`
	DBURI     string `json:"db_URI"`
}

func NewConfig(keyConfig string) Config {
	return Config{KeyConfig: keyConfig}
}

func (c *Config) Load(consulAddr string) error {

	consulKV, err := consul.OpenConsul(consulAddr)
	if err != nil {
		return err
	}

	config, err := consul.SafeConsulGet(consulKV, c.KeyConfig)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, c)
	if err != nil {
		return err
	}

	return nil
}
