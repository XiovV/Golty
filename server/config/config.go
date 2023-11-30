package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port        string `env:"PORT" env-default:"8080"`
	Environment string `env:"ENVIRONMENT" env-default:"PRODUCTION"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}
