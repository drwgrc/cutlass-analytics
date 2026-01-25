package models

import (
	"time"

	"gorm.io/gorm"
)

type IslandGovernanceHistory struct {
	gorm.Model
	IslandID     uint       `gorm:"not null;index" json:"island_id"`
	FlagID       *uint      `gorm:"index" json:"flag_id,omitempty"`
	GovernorName string     `gorm:"type:varchar(100)" json:"governor_name,omitempty"`
	StartedAt    time.Time  `gorm:"not null" json:"started_at"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	
	ChangeType string `gorm:"type:varchar(50)" json:"change_type,omitempty"`
	
	Island Island `gorm:"foreignKey:IslandID" json:"island,omitempty"`
	Flag   *Flag  `gorm:"foreignKey:FlagID" json:"flag,omitempty"`
}

func (IslandGovernanceHistory) TableName() string {
	return "island_governance_history"
}

func (h *IslandGovernanceHistory) IsCurrent() bool {
	return h.EndedAt == nil
}