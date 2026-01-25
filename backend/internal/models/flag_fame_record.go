package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type FlagFameRecord struct {
	gorm.Model
	FlagID    uint      `gorm:"uniqueIndex:idx_flag_fame_date;not null" json:"flag_id"`
	ScrapedAt time.Time `gorm:"uniqueIndex:idx_flag_fame_date;not null;index" json:"scraped_at"`

	FameLevel types.FameLevel `gorm:"type:varchar(30)" json:"fame_level"`
	FameRank  *int      `gorm:"index" json:"fame_rank,omitempty"`

	Flag Flag `gorm:"foreignKey:FlagID" json:"flag,omitempty"`
}

func (FlagFameRecord) TableName() string {
	return "flag_fame_records"
}

func (r *FlagFameRecord) IsRanked() bool {
	return r.FameRank != nil
}

func (r *FlagFameRecord) GetRank() int {
	if r.FameRank == nil {
		return 0
	}
	return *r.FameRank
}
