package repositories

import (
	"cutlass_analytics/internal/dto"
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"

	"gorm.io/gorm"
)

type FlagRepository struct {
	db *gorm.DB
}

func NewFlagRepository(db *gorm.DB) *FlagRepository {
	return &FlagRepository{db: db}
}

func (r *FlagRepository) FindByID(id uint) (*models.Flag, error) {
	var flag models.Flag
	err := r.db.First(&flag, id).Error
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

func (r *FlagRepository) FindByGameID(gameID uint64, ocean types.Ocean) (*models.Flag, error) {
	var flag models.Flag
	err := r.db.Where("game_flag_id = ? AND ocean = ?", gameID, ocean).
		First(&flag).Error
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

func (r *FlagRepository) List(req dto.FlagListRequest) ([]models.Flag, int64, error) {
	query := r.db.Model(&models.Flag{})

	// Apply filters
	if req.Ocean != "" {
		query = query.Where("ocean = ?", req.Ocean)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}
	if req.FameLevel != "" {
		// Need to join with latest fame record to filter by fame level
		query = query.Joins("LEFT JOIN flag_fame_records ON flag_fame_records.flag_id = flags.id").
			Where("flag_fame_records.fame_level = ?", req.FameLevel).
			Where("flag_fame_records.scraped_at = (SELECT MAX(scraped_at) FROM flag_fame_records WHERE flag_id = flags.id)")
	}
	if req.MinCrews > 0 {
		// Subquery to count active crews per flag
		query = query.Where("(SELECT COUNT(*) FROM crews WHERE crews.flag_id = flags.id AND crews.is_active = true) >= ?", req.MinCrews)
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

	var flags []models.Flag
	if err := query.Find(&flags).Error; err != nil {
		return nil, 0, err
	}

	return flags, total, nil
}

func (r *FlagRepository) GetCrews(flagID uint, isActive *bool) ([]models.Crew, error) {
	query := r.db.Where("flag_id = ?", flagID)
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	var crews []models.Crew
	err := query.Find(&crews).Error
	if err != nil {
		return nil, err
	}
	return crews, nil
}

func (r *FlagRepository) GetActiveCrewCount(flagID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Crew{}).
		Where("flag_id = ? AND is_active = ?", flagID, true).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FlagRepository) GetLatestFameRecord(flagID uint) (*models.FlagFameRecord, error) {
	var record models.FlagFameRecord
	err := r.db.Where("flag_id = ?", flagID).
		Order("scraped_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *FlagRepository) GetFameHistory(flagID uint) ([]models.FlagFameRecord, error) {
	var records []models.FlagFameRecord
	err := r.db.Where("flag_id = ?", flagID).
		Order("scraped_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
