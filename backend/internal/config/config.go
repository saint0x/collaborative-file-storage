package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLiteDBPath      string
	ClerkSecretKey    string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	OpenAIAPIKey      string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env.local"); err != nil {
		return nil, err
	}

	return &Config{
		SQLiteDBPath:      os.Getenv("SQLITE_DB_PATH"),
		ClerkSecretKey:    os.Getenv("CLERK_SECRET_KEY"),
		R2AccountID:       os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("CLOUDFLARE_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("CLOUDFLARE_SECRET_ACCESS_KEY"),
		R2BucketName:      os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
		OpenAIAPIKey:      os.Getenv("OPENAI_API_KEY"),
	}, nil
}
