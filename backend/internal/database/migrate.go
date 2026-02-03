package database

import (
	"cutlass_analytics/internal/models"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

    err := db.AutoMigrate(
		&models.Flag{},
		&models.Crew{},
		&models.CrewBattleRecord{},
		&models.CrewFameRecord{},
		&models.CrewReputationRecord{},
		&models.FlagFameRecord{},
		&models.CrewFlagHistory{},
		&models.ScrapeJob{},
		&models.Island{},
		&models.Archipelago{},
		&models.Commodity{},
		&models.CommodityTaxRate{},
		&models.MarketPrice{},
		&models.MarketOrder{},
		&models.IslandGovernanceHistory{},
		&models.IslandPopulation{},
		&models.IslandCommodity{},
    )
	
    if err != nil {
        return err
    }
	
	log.Println("Migrations completed successfully")
    return nil
}

func CreateIndexes(db *gorm.DB) error {
	// Index for finding latest battle record per crew
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crew_battle_latest 
		ON crew_battle_records (crew_id, scraped_at DESC)
	`).Error; err != nil {
		return err
	}

	// Index for finding latest fame record per crew
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crew_fame_latest 
		ON crew_fame_records (crew_id, scraped_at DESC)
	`).Error; err != nil {
		return err
	}

	// Index for finding latest fame record per flag
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_flag_fame_latest 
		ON flag_fame_records (flag_id, scraped_at DESC)
	`).Error; err != nil {
		return err
	}

	// Index for date range queries on battle records
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crew_battle_date_range 
		ON crew_battle_records (scraped_at, crew_id)
	`).Error; err != nil {
		return err
	}

	// Index for active crews by ocean
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crews_ocean_active 
		ON crews (ocean, is_active) WHERE is_active = true
	`).Error; err != nil {
		return err
	}

	// Index for active flags by ocean
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_flags_ocean_active 
		ON flags (ocean, is_active) WHERE is_active = true
	`).Error; err != nil {
		return err
	}

	// Index for crew flag history (current memberships)
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crew_flag_history_active 
		ON crew_flag_history (crew_id) WHERE left_at IS NULL
	`).Error; err != nil {
		return err
	}

	// Index for reputation records by type
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crew_rep_type_date 
		ON crew_reputation_records (reputation_type, scraped_at DESC)
	`).Error; err != nil {
		return err
	}

	return nil
}

func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.ScrapeJob{},
		&models.CrewFlagHistory{},
		&models.FlagFameRecord{},
		&models.CrewReputationRecord{},
		&models.CrewFameRecord{},
		&models.CrewBattleRecord{},
		&models.Crew{},
		&models.Flag{},
	)
}

func ResetDatabase(db *gorm.DB) error {
	if err := DropAllTables(db); err != nil {
		return err
	}
	if err := AutoMigrate(db); err != nil {
		return err
	}
	return CreateIndexes(db)
}
