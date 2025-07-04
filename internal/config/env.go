package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost  string
	Port        string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	SSLMode     string
	DatabaseURL string
}

var Envs = initConfig()

func initConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING: .env file not found. Using system variables.")
	}

	config := &Config{
		PublicHost:  getEnv("PUBLIC_HOST", "http://localhost"),
		Port:        getEnv("PORT", "8080"),
		DBUser:      getEnv("DB_USER", "ecom_user"),
		DBPassword:  getEnv("DB_PASSWORD", "ecom_password"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "ecom_backend"),
		SSLMode:     getEnv("DB_SSL_MODE", "disable"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if config.DatabaseURL == "" {
		config.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
			config.SSLMode,
		)
	}

	validateConfig(config)

	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func validateConfig(cfg *Config) {
	require := map[string]string{
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
		"DB_HOST":     cfg.DBHost,
		"DB_PORT":     cfg.DBPort,
		"DB_NAME":     cfg.DBName,
	}

	for name, value := range require {
		if value == "" {
			log.Fatalf("Required configuration not found: %s", name)
		}
	}

	if cfg.SSLMode != "disable" && cfg.SSLMode != "require" &&
		cfg.SSLMode != "verify-ca" && cfg.SSLMode != "verify-full" {
		log.Printf("Aviso: DB_SSL_MODE '%s' pode não ser válido. Valores aceitos: disable, require, verify-ca, verify-full", cfg.SSLMode)
	}
}

func (c *Config) GetDBConfig() map[string]string {
	return map[string]string{
		"host":     c.DBHost,
		"port":     c.DBPort,
		"user":     c.DBUser,
		"password": c.DBPassword,
		"dbname":   c.DBName,
		"sslmode":  c.SSLMode,
	}
}