package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type CrewReputationRecord struct {
	gorm.Model
	CrewID         uint           `gorm:"uniqueIndex:idx_crew_rep_date_type;not null" json:"crew_id"`
	ScrapedAt      time.Time      `gorm:"uniqueIndex:idx_crew_rep_date_type;not null;index" json:"scraped_at"`
	ReputationType types.ReputationType `gorm:"uniqueIndex:idx_crew_rep_date_type;type:varchar(20);not null" json:"reputation_type"`

	ReputationLevel string `gorm:"type:varchar(30)" json:"reputation_level"`
	ReputationRank  *int   `gorm:"index" json:"reputation_rank,omitempty"`

	Crew Crew `gorm:"foreignKey:CrewID" json:"crew,omitempty"`
}

func (CrewReputationRecord) TableName() string {
	return "crew_reputation_records"
}

func (r *CrewReputationRecord) IsRanked() bool {
	return r.ReputationRank != nil
}

func (r *CrewReputationRecord) GetRank() int {
	if r.ReputationRank == nil {
		return 0
	}
	return *r.ReputationRank
}
