package models

import (
	"cutlass_analytics/internal/types"

	"gorm.io/gorm"
)

type Archipelago struct {
	gorm.Model
	Ocean       types.Ocean  `gorm:"uniqueIndex:idx_archipelago_ocean;type:varchar(20);not null" json:"ocean"`
	Name        string `gorm:"uniqueIndex:idx_archipelago_ocean;type:varchar(100);not null" json:"name"`
	DisplayName string `gorm:"type:varchar(100)" json:"display_name"`
	
	Color string `gorm:"type:varchar(30)" json:"color,omitempty"`
	
	Islands []Island `gorm:"foreignKey:ArchipelagoID" json:"islands,omitempty"`
}

func (Archipelago) TableName() string {
	return "archipelagos"
}