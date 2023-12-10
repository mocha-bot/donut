package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ApplicationConfig ApplicationConfig
	DatabaseConfig    DatabaseConfig
}

func Get() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
