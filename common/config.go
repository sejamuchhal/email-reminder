package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger           *logrus.Entry
	DSN              string
	RedisHost        string
	RedisPassword    string
	MailersendAPIKey string
}

func ConfigureOrDie() *Config {

	redisHost := GetEnvDefault("REDIS_HOST", "127.0.0.1:6379")

	if !strings.ContainsAny(redisHost, ":") {
		redisHost = redisHost + ":6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	// mailersendAPIKey := os.Getenv("MAILERSEND_API_KEY")
	mailersendAPIKey := "mlsn.98f87d8d4cead677ebe2940b5c885ee809f07afa3dfa1eb36eb74968ad0b7dbb"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			GetEnvDefault("POSTGRES_DB_HOST", "localhost"),
			GetEnvDefault("POSTGRES_DB_PORT", "5432"),
			GetEnvDefault("POSTGRES_DB_USER", "postgres"),
			GetEnvDefault("POSTGRES_DB_NAME", "postgres"),
			GetEnvDefault("POSTGRES_DB_PASSWORD", "postgres"))
	config := &Config{
		Logger:           Logger,
		RedisHost:        redisHost,
		RedisPassword:    redisPassword,
		MailersendAPIKey: mailersendAPIKey,
		DSN:              dsn,
	}
	return config
}

func GetEnvDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}
