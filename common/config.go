package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger        *logrus.Entry
	DSN           string
	RedisHost     string
	RedisPassword string
}

func ConfigureOrDie() *Config {

	redisHost := GetEnvDefault("REDIS_HOST", "127.0.0.1:6379")

	if !strings.ContainsAny(redisHost, ":") {
		redisHost = redisHost + ":6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	config := &Config{
		Logger:        Logger,
		RedisHost:     redisHost,
		RedisPassword: redisPassword,
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			GetEnvDefault("POSTGRES_DB_HOST", "localhost"),
			GetEnvDefault("POSTGRES_DB_PORT", "5432"),
			GetEnvDefault("POSTGRES_DB_USER", "username"),
			GetEnvDefault("POSTGRES_DB_NAME", "email_reminder"),
			GetEnvDefault("POSTGRES_DB_PASSWORD", "password")),
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
