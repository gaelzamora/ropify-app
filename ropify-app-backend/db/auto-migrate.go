package db

import (
	"github.com/gaelzamora/ropify-app/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Garment{},
		&models.Outfit{},
	)
}
