package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	DatabaseConfig DatabaseConfig
}

func Get() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
