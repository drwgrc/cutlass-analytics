package models

import (
	"time"

	"gorm.io/gorm"
)

type MarketPrice struct {
	gorm.Model
	IslandID    uint      `gorm:"uniqueIndex:idx_market_price;not null;index" json:"island_id"`
	CommodityID uint      `gorm:"uniqueIndex:idx_market_price;not null;index" json:"commodity_id"`
	ScrapedAt   time.Time `gorm:"uniqueIndex:idx_market_price;not null;index" json:"scraped_at"`
	
	BuyPrice    *int `json:"buy_price,omitempty"`
	BuyQuantity *int `json:"buy_quantity,omitempty"`
	
	SellPrice    *int `json:"sell_price,omitempty"`
	SellQuantity *int `json:"sell_quantity,omitempty"`
	
	Island    Island    `gorm:"foreignKey:IslandID" json:"island,omitempty"`
	Commodity Commodity `gorm:"foreignKey:CommodityID" json:"commodity,omitempty"`
}

func (MarketPrice) TableName() string {
	return "market_prices"
}

func (p *MarketPrice) Spread() *int {
	if p.BuyPrice == nil || p.SellPrice == nil {
		return nil
	}
	spread := *p.SellPrice - *p.BuyPrice
	return &spread
}
