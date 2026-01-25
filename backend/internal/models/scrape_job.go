package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type ScrapeJobStatus string

const (
	ScrapeJobStatusRunning   ScrapeJobStatus = "running"
	ScrapeJobStatusCompleted ScrapeJobStatus = "completed"
	ScrapeJobStatusFailed    ScrapeJobStatus = "failed"
)

type ScrapeJobType string

const (
	ScrapeJobTypeCrewFame   ScrapeJobType = "crew_fame"
	ScrapeJobTypeFlagFame   ScrapeJobType = "flag_fame"
	ScrapeJobTypeCrewInfo   ScrapeJobType = "crew_info"
	ScrapeJobTypeBattleInfo ScrapeJobType = "battle_info"
	ScrapeJobTypeDailyFull  ScrapeJobType = "daily_full"
)

type ScrapeJob struct {
	gorm.Model
	Ocean     types.Ocean           `gorm:"type:varchar(20);not null;index" json:"ocean"`
	JobType   ScrapeJobType   `gorm:"type:varchar(50);not null" json:"job_type"`
	StartedAt time.Time       `gorm:"not null" json:"started_at"`
	EndedAt   *time.Time      `json:"ended_at,omitempty"`
	Status    ScrapeJobStatus `gorm:"type:varchar(20);default:'running'" json:"status"`

	ItemsProcessed int    `gorm:"default:0" json:"items_processed"`
	ItemsFailed    int    `gorm:"default:0" json:"items_failed"`
	ErrorMessage   string `gorm:"type:text" json:"error_message,omitempty"`
}

func (ScrapeJob) TableName() string {
	return "scrape_jobs"
}

func (s *ScrapeJob) BeforeCreate(tx *gorm.DB) error {
	if s.StartedAt.IsZero() {
		s.StartedAt = time.Now()
	}
	if s.Status == "" {
		s.Status = ScrapeJobStatusRunning
	}
	return nil
}

func (s *ScrapeJob) MarkCompleted(db *gorm.DB) error {
	now := time.Now()
	s.EndedAt = &now
	s.Status = ScrapeJobStatusCompleted
	return db.Save(s).Error
}

func (s *ScrapeJob) MarkFailed(db *gorm.DB, err error) error {
	now := time.Now()
	s.EndedAt = &now
	s.Status = ScrapeJobStatusFailed
	if err != nil {
		s.ErrorMessage = err.Error()
	}
	return db.Save(s).Error
}

func (s *ScrapeJob) IncrementProcessed() {
	s.ItemsProcessed++
}

func (s *ScrapeJob) IncrementFailed() {
	s.ItemsFailed++
}

func (s *ScrapeJob) Duration() time.Duration {
	endTime := time.Now()
	if s.EndedAt != nil {
		endTime = *s.EndedAt
	}
	return endTime.Sub(s.StartedAt)
}

func (s *ScrapeJob) IsRunning() bool {
	return s.Status == ScrapeJobStatusRunning
}

func (s *ScrapeJob) SuccessRate() float64 {
	total := s.ItemsProcessed + s.ItemsFailed
	if total == 0 {
		return 0
	}
	return float64(s.ItemsProcessed) / float64(total) * 100
}
