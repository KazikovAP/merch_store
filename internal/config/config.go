package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetDefault("ServerPort", "8080")
	viper.SetDefault("DB_PORT", "5432")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Не удалось прочитать файл конфигурации, будут использованы переменные окружения: %v", err)
	}

	cfg := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
		ServerPort: viper.GetString("ServerPort"),
	}

	return cfg, nil
}
