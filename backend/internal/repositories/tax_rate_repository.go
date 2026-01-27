package repositories

import (
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type TaxRateRepository struct {
	db *gorm.DB
}

func NewTaxRateRepository(db *gorm.DB) *TaxRateRepository {
	return &TaxRateRepository{db: db}
}

func (r *TaxRateRepository) GetCurrentRates(ocean types.Ocean, category string) ([]models.CommodityTaxRate, error) {
	// Get the latest tax rate for each commodity in the ocean
	query := r.db.Model(&models.CommodityTaxRate{}).
		Where("ocean = ?", ocean).
		Where("scraped_at = (SELECT MAX(scraped_at) FROM commodity_tax_rates WHERE commodity_id = commodity_tax_rates.commodity_id AND ocean = ?)", ocean).
		Preload("Commodity")

	if category != "" {
		query = query.Joins("JOIN commodities ON commodities.id = commodity_tax_rates.commodity_id").
			Where("commodities.category = ?", category)
	}

	var rates []models.CommodityTaxRate
	err := query.Find(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *TaxRateRepository) GetHistory(commodityID uint, ocean types.Ocean, startDate, endDate time.Time) ([]models.CommodityTaxRate, error) {
	var rates []models.CommodityTaxRate
	err := r.db.Where("commodity_id = ? AND ocean = ? AND scraped_at >= ? AND scraped_at <= ?", commodityID, ocean, startDate, endDate).
		Preload("Commodity").
		Order("scraped_at ASC").
		Find(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *TaxRateRepository) CompareAcrossOceans(commodityID uint) ([]models.CommodityTaxRate, error) {
	// Get the latest tax rate for this commodity across all oceans
	var rates []models.CommodityTaxRate
	
	oceans := []types.Ocean{types.OceanEmerald, types.OceanMeridian, types.OceanCerulean, types.OceanObsidian}
	
	for _, ocean := range oceans {
		var rate models.CommodityTaxRate
		err := r.db.Where("commodity_id = ? AND ocean = ?", commodityID, ocean).
			Preload("Commodity").
			Order("scraped_at DESC").
			First(&rate).Error
		if err == nil {
			rates = append(rates, rate)
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	
	return rates, nil
}
