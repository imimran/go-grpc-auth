package cmd

import (
	"log"

	"github.com/imimran/go-grpc-auth/user/domain"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}

	log.Println("✅ Auto migration complete")
	return nil
}
