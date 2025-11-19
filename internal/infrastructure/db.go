package infrastructure

import (
	"log"

	"github.com/imimran/go-grpc-auth/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate User table
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	return db
}
