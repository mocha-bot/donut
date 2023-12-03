package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseInstance(conf *Config) (*gorm.DB, error) {
	gormConf := new(gorm.Config)

	if conf.DatabaseConfig.Debug {
		gormConf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			},
		)
	}

	dialector := conf.DatabaseConfig.GetDialector()
	if dialector == nil {
		return nil, fmt.Errorf("unsupported database dialect, %s", conf.DatabaseConfig.Dialect)
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
