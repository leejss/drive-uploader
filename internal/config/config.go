package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TokenFilePath string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	cfg := &Config{
		TokenFilePath: os.Getenv("TOKEN_FILE_PATH"),
	}

	return cfg, nil
}
