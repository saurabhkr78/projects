package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func (c *Config) Validate() error {
	if c.Port == "" {
		return errors.New("PORT is not set")
	}
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is not set")
	}
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is not set")
	}
	return nil

}

func Load() (*Config, error) {
	_ = godotenv.Load() // Load environment variables from .env file

	cfg := &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
