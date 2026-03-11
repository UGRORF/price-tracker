package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port    int
	Timeout int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{}

	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	cfg.Database = *dbConfig

	serverConfig, err := loadServerConfig()
	if err != nil {
		return nil, err
	}
	cfg.Server = *serverConfig

	return cfg, nil
}

func loadDatabaseConfig() (*DatabaseConfig, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "price_tracker"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func loadServerConfig() (*ServerConfig, error) {
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	timeout, err := strconv.Atoi(getEnv("SERVER_TIMEOUT", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_TIMEOUT: %w", err)
	}

	return &ServerConfig{
		Port:    port,
		Timeout: timeout,
	}, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
