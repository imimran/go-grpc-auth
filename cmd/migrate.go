package cmd

import (
	"log"

	userDomain "github.com/imimran/go-grpc-auth/user/domain"
	addressDomain "github.com/imimran/go-grpc-auth/address/domain" // Import your address domain
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// 1. Optional: Enable PostGIS extension if not already enabled
	// This is required for the 'geography' type in your Address model
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error; err != nil {
		log.Printf("Warning: Could not enable PostGIS extension: %v", err)
	}

	// 2. Run migrations for all models
	err := db.AutoMigrate(
		&userDomain.User{},
		&addressDomain.Address{}, // Add this line
	)

	if err != nil {
		return err
	}

	log.Println("Auto migration complete for User and Address tables")
	return nil
}
