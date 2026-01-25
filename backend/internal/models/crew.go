package models

import (
	"cutlass_analytics/internal/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Crew struct {
	gorm.Model
	GameCrewID  uint64    `gorm:"uniqueIndex:idx_crew_ocean;not null" json:"game_crew_id"`
	Ocean       types.Ocean     `gorm:"uniqueIndex:idx_crew_ocean;type:varchar(20);not null" json:"ocean"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	FlagID      *uint     `gorm:"index" json:"flag_id,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	FirstSeenAt time.Time `gorm:"not null" json:"first_seen_at"`
	LastSeenAt  time.Time `gorm:"not null" json:"last_seen_at"`

	Flag              *Flag                  `gorm:"foreignKey:FlagID" json:"flag,omitempty"`
	BattleRecords     []CrewBattleRecord     `gorm:"foreignKey:CrewID" json:"battle_records,omitempty"`
	FameRecords       []CrewFameRecord       `gorm:"foreignKey:CrewID" json:"fame_records,omitempty"`
	ReputationRecords []CrewReputationRecord `gorm:"foreignKey:CrewID" json:"reputation_records,omitempty"`
}

func (Crew) TableName() string {
	return "crews"
}

func (c *Crew) BeforeCreate(tx *gorm.DB) error {
	if c.FirstSeenAt.IsZero() {
		c.FirstSeenAt = time.Now()
	}
	if c.LastSeenAt.IsZero() {
		c.LastSeenAt = time.Now()
	}
	return nil
}

func (c *Crew) GetYowebURL() string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/crew/info.wm?crewid=%d", c.Ocean, c.GameCrewID)
}

func (c *Crew) GetBattleInfoURL() string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/crew/battleinfo.wm?crewid=%d&classic=$classic", c.Ocean, c.GameCrewID)
}

func (c *Crew) HasFlag() bool {
	return c.FlagID != nil
}

func (c *Crew) GetLatestBattleRecord(db *gorm.DB) (*CrewBattleRecord, error) {
	var record CrewBattleRecord
	err := db.Where("crew_id = ?", c.ID).
		Order("scraped_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Crew) GetLatestFameRecord(db *gorm.DB) (*CrewFameRecord, error) {
	var record CrewFameRecord
	err := db.Where("crew_id = ?", c.ID).
		Order("scraped_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
