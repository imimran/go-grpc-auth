package db

import (
	"fmt"
	"log"

	"github.com/imimran/go-grpc-auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)
}

func NewPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get db instance fail: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	log.Println("Connected to Postgres DB")
	return db, nil
}
