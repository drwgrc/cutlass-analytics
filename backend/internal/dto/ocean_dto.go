package dto

import "time"

type OceanResponse struct {
	Name      string `json:"name"`
	DisplayName string `json:"display_name"`
	IsActive  bool   `json:"is_active"`
	BaseURL   string `json:"base_url"`
}

type OceanListResponse struct {
	Oceans []OceanResponse `json:"oceans"`
}

type OceanStatsResponse struct {
	Ocean           string    `json:"ocean"`
	DisplayName     string    `json:"display_name"`
	TotalCrews      int       `json:"total_crews"`
	ActiveCrews     int       `json:"active_crews"`
	TotalFlags      int       `json:"total_flags"`
	ActiveFlags     int       `json:"active_flags"`
	TotalPVPBattles int       `json:"total_pvp_battles"`
	TotalPVPWins    int       `json:"total_pvp_wins"`
	LastScrapedAt   time.Time `json:"last_scraped_at"`
	
}

type OceanComparisonResponse struct {
	Oceans    []OceanStatsResponse `json:"oceans"`
	UpdatedAt time.Time            `json:"updated_at"`
}
