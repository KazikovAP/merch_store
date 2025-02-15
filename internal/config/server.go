package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: viper.GetString("server.port"),
	}
}
