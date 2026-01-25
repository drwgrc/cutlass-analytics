package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type CrewFameRecord struct {
	gorm.Model
	CrewID    uint      `gorm:"uniqueIndex:idx_crew_fame_date;not null" json:"crew_id"`
	ScrapedAt time.Time `gorm:"uniqueIndex:idx_crew_fame_date;not null;index" json:"scraped_at"`

	FameLevel types.FameLevel `gorm:"type:varchar(30)" json:"fame_level"`
	FameRank  *int      `gorm:"index" json:"fame_rank,omitempty"`

	Crew Crew `gorm:"foreignKey:CrewID" json:"crew,omitempty"`
}

func (CrewFameRecord) TableName() string {
	return "crew_fame_records"
}

func (r *CrewFameRecord) IsRanked() bool {
	return r.FameRank != nil
}

func (r *CrewFameRecord) GetRank() int {
	if r.FameRank == nil {
		return 0
	}
	return *r.FameRank
}
