package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type CommodityTaxRate struct {
	gorm.Model
	CommodityID uint      `gorm:"uniqueIndex:idx_commodity_ocean_date;not null" json:"commodity_id"`
	Ocean       types.Ocean     `gorm:"uniqueIndex:idx_commodity_ocean_date;type:varchar(20);not null" json:"ocean"`
	ScrapedAt   time.Time `gorm:"uniqueIndex:idx_commodity_ocean_date;not null;index" json:"scraped_at"`
	
	TaxValue int `gorm:"not null" json:"tax_value"`
	
	Commodity Commodity `gorm:"foreignKey:CommodityID" json:"commodity,omitempty"`
}

func (CommodityTaxRate) TableName() string {
	return "commodity_tax_rates"
}
