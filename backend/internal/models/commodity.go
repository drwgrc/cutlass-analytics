package models

import (
	"cutlass_analytics/internal/types"

	"gorm.io/gorm"
)

type Commodity struct {
	gorm.Model
	Name        string            `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
	DisplayName string            `gorm:"type:varchar(100);not null" json:"display_name"`
	Category    types.CommodityCategory `gorm:"type:varchar(30);not null;index" json:"category"`
	
	IsSpawnable bool `gorm:"default:false" json:"is_spawnable"`
	
	IsRare bool `gorm:"default:false" json:"is_rare"`
	
	Description string `gorm:"type:text" json:"description,omitempty"`
	
	IconPath string `gorm:"type:varchar(255)" json:"icon_path,omitempty"`
	
	IslandCommodities []IslandCommodity `gorm:"foreignKey:CommodityID" json:"island_commodities,omitempty"`
	TaxRates          []CommodityTaxRate `gorm:"foreignKey:CommodityID" json:"tax_rates,omitempty"`
}

func (Commodity) TableName() string {
	return "commodities"
}