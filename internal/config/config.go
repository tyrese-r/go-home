package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	ServerAddress string
	DBPath        string
}

// New returns a Config with values from environment variables or defaults
func New() *Config {
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data.db"
	}

	return &Config{
		ServerAddress: serverAddr,
		DBPath:        dbPath,
	}
}
