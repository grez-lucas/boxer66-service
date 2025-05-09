package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	configInstance *Config
	loadConfigOnce sync.Once
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	SMTPConfig  SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func LoadConfig() *Config {
	loadConfigOnce.Do(func() {
		configInstance = load()
	})
	return configInstance
}

func load() *Config {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			slog.Warn("No .env file found, defaulting to system environment variables")
		}
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DB_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		SMTPConfig: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			User:     os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}

	return cfg
}
