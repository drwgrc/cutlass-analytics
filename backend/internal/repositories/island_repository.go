package repositories

import (
	"cutlass_analytics/internal/dto"
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type IslandRepository struct {
	db *gorm.DB
}

func NewIslandRepository(db *gorm.DB) *IslandRepository {
	return &IslandRepository{db: db}
}

func (r *IslandRepository) FindByID(id uint) (*models.Island, error) {
	var island models.Island
	err := r.db.Preload("Archipelago").Preload("GovernorFlag").
		First(&island, id).Error
	if err != nil {
		return nil, err
	}
	return &island, nil
}

func (r *IslandRepository) FindByGameID(gameID uint64, ocean types.Ocean) (*models.Island, error) {
	var island models.Island
	err := r.db.Preload("Archipelago").Preload("GovernorFlag").
		Where("game_island_id = ? AND ocean = ?", gameID, ocean).
		First(&island).Error
	if err != nil {
		return nil, err
	}
	return &island, nil
}

func (r *IslandRepository) List(req dto.IslandListRequest) ([]models.Island, int64, error) {
	query := r.db.Model(&models.Island{})

	// Apply filters
	if req.Ocean != "" {
		query = query.Where("ocean = ?", req.Ocean)
	}
	if req.ArchipelagoID != nil {
		query = query.Where("archipelago_id = ?", *req.ArchipelagoID)
	}
	if req.Size != "" {
		query = query.Where("size = ?", req.Size)
	}
	if req.IsColonized != nil {
		query = query.Where("is_colonized = ?", *req.IsColonized)
	}
	if req.GovernorFlagID != nil {
		query = query.Where("governor_flag_id = ?", *req.GovernorFlagID)
	}
	if req.HasCommodity != "" {
		query = query.Joins("JOIN island_commodities ON island_commodities.island_id = islands.id").
			Joins("JOIN commodities ON commodities.id = island_commodities.commodity_id").
			Where("commodities.name = ? OR commodities.display_name = ?", req.HasCommodity, req.HasCommodity)
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
	query = query.Preload("Archipelago").Preload("GovernorFlag")

	var islands []models.Island
	if err := query.Find(&islands).Error; err != nil {
		return nil, 0, err
	}

	return islands, total, nil
}

func (r *IslandRepository) GetLatestPopulation(islandID uint) (*models.IslandPopulation, error) {
	var population models.IslandPopulation
	err := r.db.Where("island_id = ?", islandID).
		Order("scraped_at DESC").
		First(&population).Error
	if err != nil {
		return nil, err
	}
	return &population, nil
}

func (r *IslandRepository) GetPopulationHistory(islandID uint, startDate, endDate time.Time) ([]models.IslandPopulation, error) {
	var populations []models.IslandPopulation
	err := r.db.Where("island_id = ? AND scraped_at >= ? AND scraped_at <= ?", islandID, startDate, endDate).
		Order("scraped_at ASC").
		Find(&populations).Error
	if err != nil {
		return nil, err
	}
	return populations, nil
}

func (r *IslandRepository) GetGovernanceHistory(islandID uint) ([]models.IslandGovernanceHistory, error) {
	var history []models.IslandGovernanceHistory
	err := r.db.Where("island_id = ?", islandID).
		Preload("Flag").
		Order("started_at DESC").
		Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *IslandRepository) GetCommodities(islandID uint) ([]models.IslandCommodity, error) {
	var commodities []models.IslandCommodity
	err := r.db.Where("island_id = ?", islandID).
		Preload("Commodity").
		Find(&commodities).Error
	if err != nil {
		return nil, err
	}
	return commodities, nil
}
