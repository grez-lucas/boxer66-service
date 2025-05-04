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
	}

	return cfg
}
