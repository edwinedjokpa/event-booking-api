package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Config struct {
	DatabaseURL string
	RedisAddr   string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   os.Getenv("REDIS_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return nil, os.ErrNotExist
	}

	if cfg.RedisAddr == "" {
		return nil, os.ErrNotExist
	}

	if cfg.JWTSecret == "" {
		return nil, os.ErrNotExist
	}

	return cfg, nil
}
