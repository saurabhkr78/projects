package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
}

func Load() (*Config, error) {
	// Try to load .env for local development.
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	cfg := &Config{
		Port:         os.Getenv("PORT"),
		MongoURI:     os.Getenv("MONGO_URI"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
	}

	if cfg.MongoURI == "" {
		return nil, errors.New("MONGO_URI is required")
	}

	if cfg.DatabaseName == "" {
		return nil, errors.New("DATABASE_NAME is required")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
