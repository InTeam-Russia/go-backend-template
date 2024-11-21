package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	PostgresUrl         string
	RedisUrl            string
	SessionCookieSecure bool
	SessionCookieDomain string
}

func LoadConfigFromEnv(logger *zap.Logger) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn(".env file not found")
	}

	sessionCookieSecure, err := strconv.ParseBool(os.Getenv("SESSION_COOKIE_SECURE"))
	if err != nil {
		return nil, err
	}

	return &Config{
		PostgresUrl:         os.Getenv("POSTGRES_URL"),
		RedisUrl:            os.Getenv("REDIS_URL"),
		SessionCookieSecure: sessionCookieSecure,
		SessionCookieDomain: os.Getenv("SESSION_COOKIE_DOMAIN"),
	}, nil
}
