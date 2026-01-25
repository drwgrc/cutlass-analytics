package database

import (
	"cutlass_analytics/internal/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

    err := db.AutoMigrate(
		&models.TaxRate{},
		&models.MarketOrder{},
    )
    if err != nil {
        return err
    }
	
	log.Println("Migrations completed successfully")
    return nil
}