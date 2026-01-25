package models

import (
	"time"

	"gorm.io/gorm"
)

type TaxRate struct {
	ID uint `gorm:"primaryKey" json:"id"`
	CommodityName string `gorm:"size:100;not null" json:"commodity_name"`
	TaxRate float64 `gorm:"not null" json:"tax_rate"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}