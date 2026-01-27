package repositories

import (
	"cutlass_analytics/internal/dto"
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type CrewRepository struct {
	db *gorm.DB
}

func NewCrewRepository(db *gorm.DB) *CrewRepository {
	return &CrewRepository{db: db}
}

func (r *CrewRepository) FindByID(id uint) (*models.Crew, error) {
	var crew models.Crew
	err := r.db.Preload("Flag").
		First(&crew, id).Error
	if err != nil {
		return nil, err
	}
	return &crew, nil
}

func (r *CrewRepository) FindByGameID(gameID uint64, ocean types.Ocean) (*models.Crew, error) {
	var crew models.Crew
	err := r.db.Preload("Flag").
		Where("game_crew_id = ? AND ocean = ?", gameID, ocean).
		First(&crew).Error
	if err != nil {
		return nil, err
	}
	return &crew, nil
}

func (r *CrewRepository) List(req dto.CrewListRequest) ([]models.Crew, int64, error) {
	query := r.db.Model(&models.Crew{})

	// Apply filters
	if req.Ocean != "" {
		query = query.Where("ocean = ?", req.Ocean)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}
	if req.FlagID != nil {
		query = query.Where("flag_id = ?", *req.FlagID)
	}
	if req.CrewRank != "" {
		// Need to join with latest battle record to filter by rank
		query = query.Joins("LEFT JOIN crew_battle_records ON crew_battle_records.crew_id = crews.id").
			Where("crew_battle_records.crew_rank = ?", req.CrewRank).
			Where("crew_battle_records.scraped_at = (SELECT MAX(scraped_at) FROM crew_battle_records WHERE crew_id = crews.id)")
	}
	if req.FameLevel != "" {
		// Need to join with latest fame record to filter by fame level
		query = query.Joins("LEFT JOIN crew_fame_records ON crew_fame_records.crew_id = crews.id").
			Where("crew_fame_records.fame_level = ?", req.FameLevel).
			Where("crew_fame_records.scraped_at = (SELECT MAX(scraped_at) FROM crew_fame_records WHERE crew_id = crews.id)")
	}

	// Count total before pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	query = applySorting(query, req.SortBy, req.SortOrder)

	// Apply pagination
	query = query.Offset(req.Offset()).Limit(req.Limit())

	// Preload relationships
	query = query.Preload("Flag")

	var crews []models.Crew
	if err := query.Find(&crews).Error; err != nil {
		return nil, 0, err
	}

	return crews, total, nil
}

func (r *CrewRepository) GetLatestBattleRecord(crewID uint) (*models.CrewBattleRecord, error) {
	var record models.CrewBattleRecord
	err := r.db.Where("crew_id = ?", crewID).
		Order("scraped_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *CrewRepository) GetBattleRecords(crewID uint, startDate, endDate time.Time) ([]models.CrewBattleRecord, error) {
	var records []models.CrewBattleRecord
	err := r.db.Where("crew_id = ? AND scraped_at >= ? AND scraped_at <= ?", crewID, startDate, endDate).
		Order("scraped_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *CrewRepository) GetLatestFameRecord(crewID uint) (*models.CrewFameRecord, error) {
	var record models.CrewFameRecord
	err := r.db.Where("crew_id = ?", crewID).
		Order("scraped_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *CrewRepository) GetFameHistory(crewID uint) ([]models.CrewFameRecord, error) {
	var records []models.CrewFameRecord
	err := r.db.Where("crew_id = ?", crewID).
		Order("scraped_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *CrewRepository) GetCurrentStats(crewID uint) (*models.CrewBattleRecord, error) {
	return r.GetLatestBattleRecord(crewID)
}
