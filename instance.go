package main

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseInstance(conf *Config) (*gorm.DB, error) {
	gormConf := new(gorm.Config)

	if conf.DatabaseConfig.LogLevel != "" {
		gormConf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      conf.DatabaseConfig.GetLogLevel(),
				Colorful:      true,
			},
		)
	}

	dialector, err := conf.DatabaseConfig.GetDialector()
	if err != nil {
		return nil, err
	}

	instance, err := gorm.Open(dialector, gormConf)
	if err != nil {
		return nil, err
	}

	if conf.DatabaseConfig.Debug {
		return instance.Debug(), nil
	}

	return instance, nil
}
