package config

import (
	"os"
)

type DBConfig struct {
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Config struct {
	DB DBConfig
}

func LoadConfig() *Config {
	return &Config{
		DB: DBConfig{
			User:     getEnv("DB_USER", "default_user"),
			Password: getEnv("DB_PASSWORD", "default_password"),
			DBName:   getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
