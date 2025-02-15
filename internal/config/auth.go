package config

import "github.com/spf13/viper"

type AuthConfig struct {
	JWTSecret string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		JWTSecret: viper.GetString("auth.jwt_secret"),
	}
}
