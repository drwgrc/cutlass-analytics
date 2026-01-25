package models

import (
	"time"

	"gorm.io/gorm"
)

type MarketOrder struct {
	ID uint `gorm:"primaryKey" json:"id"`
	IslandName string `gorm:"size:100;not null" json:"island_name"`
	CommodityName string `gorm:"size:100;not null" json:"commodity_name"`
	StoreName string `gorm:"size:100;not null" json:"store_name"`
	BuyPrice float64 `gorm:"not null" json:"buy_price"`
	BuyQuantity int `gorm:"not null" json:"buy_quantity"`
	SellPrice float64 `gorm:"not null" json:"sell_price"`
	SellQuantity int `gorm:"not null" json:"sell_quantity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}