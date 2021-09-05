package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type environment struct {
	DSN     string
	Brokers []string
}

func readEnvironment() (environment, error) {
	// Ignore error because .env file may not exist, in this case real environment variables will be used
	_ = godotenv.Load()

	dsn, ok := os.LookupEnv("DATABASE_CONNECTION_STRING")
	if !ok || len(dsn) == 0 {
		return environment{}, fmt.Errorf("DATABASE_CONNECTION_STRING environment variable is required")
	}

	brokers, ok := os.LookupEnv("KAFKA_BROKERS")
	if !ok || len(brokers) == 0 {
		return environment{}, fmt.Errorf("KAFKA_BROKERS environment variable is required")
	}

	env := environment{
		DSN:     dsn,
		Brokers: strings.Split(brokers, ","),
	}

	return env, nil
}
