package models

import (
	"cutlass_analytics/internal/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Flag struct {
	gorm.Model
	GameFlagID  uint64    `gorm:"uniqueIndex:idx_flag_ocean;not null" json:"game_flag_id"`
	Ocean       types.Ocean     `gorm:"uniqueIndex:idx_flag_ocean;type:varchar(20);not null" json:"ocean"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	FirstSeenAt time.Time `gorm:"not null" json:"first_seen_at"`
	LastSeenAt  time.Time `gorm:"not null" json:"last_seen_at"`

	Crews           []Crew           `gorm:"foreignKey:FlagID" json:"crews,omitempty"`
	FlagFameRecords []FlagFameRecord `gorm:"foreignKey:FlagID" json:"fame_records,omitempty"`
}

func (Flag) TableName() string {
	return "flags"
}

func (f *Flag) BeforeCreate(tx *gorm.DB) error {
	if f.FirstSeenAt.IsZero() {
		f.FirstSeenAt = time.Now()
	}
	if f.LastSeenAt.IsZero() {
		f.LastSeenAt = time.Now()
	}
	return nil
}

func (f *Flag) GetYowebURL() string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/flag/info.wm?flagid=%d", f.Ocean, f.GameFlagID)
}

func (f *Flag) ActiveCrewCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&Crew{}).Where("flag_id = ? AND is_active = ?", f.ID, true).Count(&count)
	return count
}
