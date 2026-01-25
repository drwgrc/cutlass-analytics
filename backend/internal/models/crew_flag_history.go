package models

import (
	"time"

	"gorm.io/gorm"
)

type CrewFlagHistory struct {
	gorm.Model
	CrewID   uint       `gorm:"not null;index" json:"crew_id"`
	FlagID   *uint      `gorm:"index" json:"flag_id,omitempty"`
	JoinedAt time.Time  `gorm:"not null" json:"joined_at"`
	LeftAt   *time.Time `json:"left_at,omitempty"`

	Crew Crew  `gorm:"foreignKey:CrewID" json:"crew,omitempty"`
	Flag *Flag `gorm:"foreignKey:FlagID" json:"flag,omitempty"`
}

func (CrewFlagHistory) TableName() string {
	return "crew_flag_history"
}

func (h *CrewFlagHistory) IsActive() bool {
	return h.LeftAt == nil
}

func (h *CrewFlagHistory) Duration() time.Duration {
	endTime := time.Now()
	if h.LeftAt != nil {
		endTime = *h.LeftAt
	}
	return endTime.Sub(h.JoinedAt)
}

func (h *CrewFlagHistory) DurationDays() int {
	return int(h.Duration().Hours() / 24)
}

func (h *CrewFlagHistory) IsIndependent() bool {
	return h.FlagID == nil
}
