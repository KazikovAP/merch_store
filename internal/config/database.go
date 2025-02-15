package config

import (
	"time"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	RetryCount int
	RetryPause time.Duration
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:       viper.GetString("database.host"),
		Port:       viper.GetString("database.port"),
		User:       viper.GetString("database.user"),
		Password:   viper.GetString("database.password"),
		Name:       viper.GetString("database.name"),
		RetryCount: viper.GetInt("database.retry_count"),
		RetryPause: viper.GetDuration("database.retry_pause"),
	}
}
