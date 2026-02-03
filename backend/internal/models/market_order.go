package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

// MarketOrder represents market order data from the buysell CSV export
// Each row in the CSV contains both buy and sell data for a specific shop
// Source: https://{ocean}.puzzlepirates.com/yoweb/econ/buysell.wm
type MarketOrder struct {
	gorm.Model
	Ocean         types.Ocean `gorm:"type:varchar(20);not null;index" json:"ocean"`
	IslandName    string      `gorm:"type:varchar(100);not null;index" json:"island_name"`
	CommodityName string      `gorm:"type:varchar(100);not null;index" json:"commodity_name"`
	ShopName      string      `gorm:"type:varchar(100);not null" json:"shop_name"`

	BuyPrice     int `gorm:"not null;default:0" json:"buy_price"`
	BuyQuantity  int `gorm:"not null;default:0" json:"buy_quantity"`
	SellPrice    int `gorm:"not null;default:0" json:"sell_price"`
	SellQuantity int `gorm:"not null;default:0" json:"sell_quantity"`

	ImportedAt time.Time `gorm:"not null;index" json:"imported_at"`
}

func (MarketOrder) TableName() string {
	return "market_orders"
}
