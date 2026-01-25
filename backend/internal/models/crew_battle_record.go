package models

import (
	"cutlass_analytics/internal/types"
	"time"

	"gorm.io/gorm"
)

type CrewBattleRecord struct {
	gorm.Model
	CrewID    uint      `gorm:"uniqueIndex:idx_crew_battle_date;not null" json:"crew_id"`
	ScrapedAt time.Time `gorm:"uniqueIndex:idx_crew_battle_date;not null;index" json:"scraped_at"`

	CrewRank types.CrewRank `gorm:"type:varchar(30)" json:"crew_rank"`

	TotalPVPWins   int `gorm:"default:0" json:"total_pvp_wins"`
	TotalPVPLosses int `gorm:"default:0" json:"total_pvp_losses"`

	DailyPVPWins   int `gorm:"default:0" json:"daily_pvp_wins"`
	DailyPVPLosses int `gorm:"default:0" json:"daily_pvp_losses"`

	DataHash string `gorm:"type:varchar(64)" json:"data_hash,omitempty"`

	Crew Crew `gorm:"foreignKey:CrewID" json:"crew,omitempty"`
}

func (CrewBattleRecord) TableName() string {
	return "crew_battle_records"
}

func (r *CrewBattleRecord) WinRate() float64 {
	total := r.TotalPVPWins + r.TotalPVPLosses
	if total == 0 {
		return 0
	}
	return float64(r.TotalPVPWins) / float64(total) * 100
}

func (r *CrewBattleRecord) TotalBattles() int {
	return r.TotalPVPWins + r.TotalPVPLosses
}

func (r *CrewBattleRecord) DailyWinRate() float64 {
	total := r.DailyPVPWins + r.DailyPVPLosses
	if total == 0 {
		return 0
	}
	return float64(r.DailyPVPWins) / float64(total) * 100
}

func (r *CrewBattleRecord) DailyTotalBattles() int {
	return r.DailyPVPWins + r.DailyPVPLosses
}

func (r *CrewBattleRecord) HasActivity() bool {
	return r.DailyPVPWins > 0 || r.DailyPVPLosses > 0
}

func (r *CrewBattleRecord) CalculateDeltas(previous *CrewBattleRecord) {
	if previous == nil {
		r.DailyPVPWins = r.TotalPVPWins
		r.DailyPVPLosses = r.TotalPVPLosses
		return
	}
	r.DailyPVPWins = r.TotalPVPWins - previous.TotalPVPWins
	r.DailyPVPLosses = r.TotalPVPLosses - previous.TotalPVPLosses
}
