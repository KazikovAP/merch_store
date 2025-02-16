package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Тест на успешную загрузку конфигурации.
func TestLoadConfig_Success(t *testing.T) {
	configContent := `
server:
  port: 8080
  read_timeout: 5s
  write_timeout: 10s
db:
  host: localhost
  port: 5432
  username: user
  password: pass
  dbname: testdb
  sslmode: disable
auth:
  signing_key: secretkey
  token_ttl: 24h
`
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/config.yaml"
	err := os.WriteFile(tmpFile, []byte(configContent), 0o600)

	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}

	viper.Reset()
	viper.AddConfigPath(tmpDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	cfg, err := config.LoadConfig(tmpDir)
	assert.NoError(t, err)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 5*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "user", cfg.DB.Username)
	assert.Equal(t, "pass", cfg.DB.Password)
	assert.Equal(t, "testdb", cfg.DB.DBName)
	assert.Equal(t, "disable", cfg.DB.SSLMode)
	assert.Equal(t, "secretkey", cfg.Auth.SigningKey)
	assert.Equal(t, 24*time.Hour, cfg.Auth.TokenTTL)
}

// Тест на ошибку при чтении конфигурации.
func TestLoadConfig_ReadError(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config")
}

// Тест на ошибку при разборе конфигурации.
func TestLoadConfig_UnmarshalError(t *testing.T) {
	configContent := `
server:
  port: "not_a_number" # Некорректный тип данных
`
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/config.yaml"

	err := os.WriteFile(tmpFile, []byte(configContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}

	viper.Reset()
	viper.AddConfigPath(tmpDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	_, err = config.LoadConfig(tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config")
}

// Тест на корректность формирования строки DSN.
func TestDatabaseConfig_DSN(t *testing.T) {
	dbConfig := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "pass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expectedDSN := "host=localhost port=5432 user=user password=pass dbname=testdb sslmode=disable"
	assert.Equal(t, expectedDSN, dbConfig.DSN())
}
