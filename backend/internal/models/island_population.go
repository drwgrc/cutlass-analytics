package models

import (
	"time"

	"gorm.io/gorm"
)

type IslandPopulation struct {
	gorm.Model
	IslandID   uint      `gorm:"uniqueIndex:idx_island_pop_date;not null" json:"island_id"`
	ScrapedAt  time.Time `gorm:"uniqueIndex:idx_island_pop_date;not null;index" json:"scraped_at"`
	Population int       `gorm:"not null" json:"population"`
	
	Island Island `gorm:"foreignKey:IslandID" json:"island,omitempty"`
}

func (IslandPopulation) TableName() string {
	return "island_populations"
}