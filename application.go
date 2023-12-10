package main

import "fmt"

type ApplicationConfig struct {
	Host string `env:"APPLICATION_HOST" envDefault:"localhost"`
	Port int    `env:"APPLICATION_PORT" envDefault:"8080"`
	TZ   string `env:"APPLICATION_TZ" envDefault:"Asia/Jakarta"`
}

func (cfg ApplicationConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
