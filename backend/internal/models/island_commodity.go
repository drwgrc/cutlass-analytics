package models

import "gorm.io/gorm"

type IslandCommodity struct {
	gorm.Model
	IslandID    uint `gorm:"uniqueIndex:idx_island_commodity;not null" json:"island_id"`
	CommodityID uint `gorm:"uniqueIndex:idx_island_commodity;not null" json:"commodity_id"`
	
	IsConfirmed bool `gorm:"default:true" json:"is_confirmed"`
	
	Island    Island    `gorm:"foreignKey:IslandID" json:"island,omitempty"`
	Commodity Commodity `gorm:"foreignKey:CommodityID" json:"commodity,omitempty"`
}

func (IslandCommodity) TableName() string {
	return "island_commodities"
}
