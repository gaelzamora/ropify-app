package db

import (
	"log"

	"github.com/gaelzamora/ropify-app/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	// Crear la extensión uuid-ossp primero
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatalf("Error creating uuid-ossp extension: %v", err)
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Garment{},
		&models.Outfit{},
	)
}
