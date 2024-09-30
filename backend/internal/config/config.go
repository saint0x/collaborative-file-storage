package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLiteDBPath string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env.local"); err != nil {
		return nil, err
	}

	return &Config{
		SQLiteDBPath: os.Getenv("SQLITE_DB_PATH"),
	}, nil
}
