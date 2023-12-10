package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseDialect string

const (
	DialectMySQL    DatabaseDialect = "mysql"
	DialectPostgres DatabaseDialect = "postgres"
)

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT" envDefault:"3306"`
	Username string `env:"DATABASE_USERNAME" envDefault:"root"`
	Password string `env:"DATABASE_PASSWORD"`
	Schema   string `env:"DATABASE_SCHEMA"`
	Debug    bool   `env:"DATABASE_DEBUG" envDefault:"false"`
	LogLevel string `env:"DATABASE_LOG_LEVEL" envDefault:"info" enum:"silent,error,warn,info"`
	Dialect  string `env:"DATABASE_DIALECT"`
}

func (d DatabaseConfig) GetDialector() (gorm.Dialector, error) {
	switch d.Dialect {
	case string(DialectMySQL):
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			d.Username,
			d.Password,
			d.Host,
			d.Port,
			d.Schema,
		)
		return mysql.Open(dsn), nil
	case string(DialectPostgres):
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
			d.Host,
			d.Username,
			d.Password,
			d.Schema,
			d.Port,
		)
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", d.Dialect)
	}
}

func (d DatabaseConfig) GetLogLevel() logger.LogLevel {
	switch d.LogLevel {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}
