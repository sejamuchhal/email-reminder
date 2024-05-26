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
	RabbitMQURL      string
	ReminderQueue    string
}

func ConfigureOrDie() *Config {

	redisHost := GetEnvDefault("REDIS_HOST", "127.0.0.1:6379")

	reminderQueuequeue := GetEnvDefault("REMINDER_QUEUE", "reminder_queue")

	if !strings.ContainsAny(redisHost, ":") {
		redisHost = redisHost + ":6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	mailersendAPIKey := os.Getenv("MAILERSEND_API_KEY")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

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
		RabbitMQURL:      rabbitMQURL,
		ReminderQueue:    reminderQueuequeue,
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
