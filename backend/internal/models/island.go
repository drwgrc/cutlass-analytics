package models

import (
	"cutlass_analytics/internal/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Island struct {
	gorm.Model
	GameIslandID uint64     `gorm:"uniqueIndex:idx_island_ocean;not null" json:"game_island_id"`
	Ocean        types.Ocean      `gorm:"uniqueIndex:idx_island_ocean;type:varchar(20);not null" json:"ocean"`
	Name         string     `gorm:"type:varchar(100);not null;index" json:"name"`
	
	ArchipelagoID *uint  `gorm:"index" json:"archipelago_id,omitempty"`
	
	Size        types.IslandSize `gorm:"type:varchar(20)" json:"size"`
	IsColonized bool       `gorm:"default:false;index" json:"is_colonized"`
	
	GovernorFlagID *uint  `gorm:"index" json:"governor_flag_id,omitempty"`
	GovernorName   string `gorm:"type:varchar(100)" json:"governor_name,omitempty"`
	
	Population int `gorm:"default:0" json:"population"`
	
	FirstSeenAt time.Time `gorm:"not null" json:"first_seen_at"`
	LastSeenAt  time.Time `gorm:"not null" json:"last_seen_at"`
	
	Archipelago       *Archipelago       `gorm:"foreignKey:ArchipelagoID" json:"archipelago,omitempty"`
	GovernorFlag      *Flag              `gorm:"foreignKey:GovernorFlagID" json:"governor_flag,omitempty"`
	Commodities       []IslandCommodity  `gorm:"foreignKey:IslandID" json:"commodities,omitempty"`
}

func (Island) TableName() string {
	return "islands"
}

func (i *Island) BeforeCreate(tx *gorm.DB) error {
	if i.FirstSeenAt.IsZero() {
		i.FirstSeenAt = time.Now()
	}
	if i.LastSeenAt.IsZero() {
		i.LastSeenAt = time.Now()
	}
	return nil
}

func (i *Island) GetYowebURL() string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/island/info.wm?islandid=%d", i.Ocean, i.GameIslandID)
}