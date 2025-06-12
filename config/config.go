package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig holds database config from .env
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// LoadConfig reads from .env and returns config struct
func LoadConfig() *DBConfig {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using system env variables.")
	}

	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

// ConnectDatabase opens GORM connection to PostgreSQL
func ConnectDatabase(cfg *DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	var db *gorm.DB
	var err error
	maxAttempts := 10

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		logrus.Infof("waiting for database (attempt %d/%d)...", attempts, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logrus.WithError(err).Fatalf("failed to connect to database after %d attempts.", maxAttempts)
	}

	return db
}
