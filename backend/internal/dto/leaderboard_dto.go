package dto

import "time"

// Request types
type CrewLeaderboardRequest struct {
	OceanParam
	PaginationParams
	
	Type string `form:"type" binding:"omitempty,oneof=wins win_rate battles rank"`
	
	MinBattles int `form:"min_battles" binding:"omitempty,min=0"`
}

func (r *CrewLeaderboardRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.Type == "" {
		r.Type = "wins"
	}
	if r.Type == "win_rate" && r.MinBattles == 0 {
		r.MinBattles = 10
	}
}

type FlagLeaderboardRequest struct {
	OceanParam
	PaginationParams
	
	Type string `form:"type" binding:"omitempty,oneof=wins win_rate battles crews fame"`
	
	MinBattles int `form:"min_battles" binding:"omitempty,min=0"`
	MinCrews   int `form:"min_crews" binding:"omitempty,min=0"`
}

func (r *FlagLeaderboardRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.Type == "" {
		r.Type = "wins"
	}
	if r.Type == "win_rate" && r.MinBattles == 0 {
		r.MinBattles = 10
	}
}

type FameLeaderboardRequest struct {
	OceanParam
	PaginationParams
	
	EntityType string `form:"entity_type" binding:"required,oneof=crew flag"`
}

func (r *FameLeaderboardRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type ReputationLeaderboardRequest struct {
	OceanParam
	PaginationParams
	
	EntityType string `form:"entity_type" binding:"required,oneof=crew flag"`
	
	ReputationType string `form:"reputation_type" binding:"required,oneof=Conqueror Explorer Patron Magnate"`
}

func (r *ReputationLeaderboardRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type DailyLeaderboardRequest struct {
	OceanParam
	PaginationParams
	
	Date string `form:"date" binding:"omitempty"`
	
	Type string `form:"type" binding:"omitempty,oneof=wins battles win_rate"`
}

func (r *DailyLeaderboardRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.Type == "" {
		r.Type = "wins"
	}
}

type LeaderboardSummaryRequest struct {
	OceanParam
	
	Limit int `form:"limit" binding:"omitempty,min=1,max=25"`
}

func (r *LeaderboardSummaryRequest) SetDefaults() {
	if r.Limit == 0 {
		r.Limit = 5
	}
}

// Response types
type CrewLeaderboardEntryResponse struct {
	Rank           int     `json:"rank"`
	CrewID         uint    `json:"crew_id"`
	GameCrewID     uint64  `json:"game_crew_id"`
	Name           string  `json:"name"`
	Ocean          string  `json:"ocean"`
	FlagID         *uint   `json:"flag_id,omitempty"`
	FlagName       string  `json:"flag_name,omitempty"`
	CrewRank       string  `json:"crew_rank"`
	TotalPVPWins   int     `json:"total_pvp_wins"`
	TotalPVPLosses int     `json:"total_pvp_losses"`
	WinRate        float64 `json:"win_rate"`
	TotalBattles   int     `json:"total_battles"`
}

type CrewLeaderboardResponse struct {
	Ocean       string                         `json:"ocean"`
	Type        string                         `json:"type"`
	Title       string                         `json:"title"`
	Description string                         `json:"description,omitempty"`
	UpdatedAt   time.Time                      `json:"updated_at"`
	Entries     []CrewLeaderboardEntryResponse `json:"entries"`
	Pagination  Pagination                     `json:"pagination"`
}

type FlagLeaderboardEntryResponse struct {
	Rank           int     `json:"rank"`
	FlagID         uint    `json:"flag_id"`
	GameFlagID     uint64  `json:"game_flag_id"`
	Name           string  `json:"name"`
	Ocean          string  `json:"ocean"`
	FameLevel      string  `json:"fame_level"`
	CrewCount      int     `json:"crew_count"`
	TotalPVPWins   int     `json:"total_pvp_wins"`
	TotalPVPLosses int     `json:"total_pvp_losses"`
	WinRate        float64 `json:"win_rate"`
	TotalBattles   int     `json:"total_battles"`
}

type FlagLeaderboardResponse struct {
	Ocean       string                         `json:"ocean"`
	Type        string                         `json:"type"`
	Title       string                         `json:"title"`
	Description string                         `json:"description,omitempty"`
	UpdatedAt   time.Time                      `json:"updated_at"`
	Entries     []FlagLeaderboardEntryResponse `json:"entries"`
	Pagination  Pagination                     `json:"pagination"`
}

type CrewFameLeaderboardEntryResponse struct {
	Rank      int    `json:"rank"`
	CrewID    uint   `json:"crew_id"`
	Name      string `json:"name"`
	Ocean     string `json:"ocean"`
	FlagName  string `json:"flag_name,omitempty"`
	FameLevel string `json:"fame_level"`
	CrewRank  string `json:"crew_rank"`
}

type CrewFameLeaderboardResponse struct {
	Ocean      string                             `json:"ocean"`
	UpdatedAt  time.Time                          `json:"updated_at"`
	Entries    []CrewFameLeaderboardEntryResponse `json:"entries"`
	Pagination Pagination                         `json:"pagination"`
}

type FlagFameLeaderboardEntryResponse struct {
	Rank      int    `json:"rank"`
	FlagID    uint   `json:"flag_id"`
	Name      string `json:"name"`
	Ocean     string `json:"ocean"`
	FameLevel string `json:"fame_level"`
	CrewCount int    `json:"crew_count"`
}

type FlagFameLeaderboardResponse struct {
	Ocean      string                             `json:"ocean"`
	UpdatedAt  time.Time                          `json:"updated_at"`
	Entries    []FlagFameLeaderboardEntryResponse `json:"entries"`
	Pagination Pagination                         `json:"pagination"`
}

type ReputationLeaderboardEntryResponse struct {
	Rank            int    `json:"rank"`
	EntityID        uint   `json:"entity_id"`
	EntityType      string `json:"entity_type"`
	Name            string `json:"name"`
	Ocean           string `json:"ocean"`
	ReputationLevel string `json:"reputation_level"`
}

type ReputationLeaderboardResponse struct {
	Ocean          string                               `json:"ocean"`
	ReputationType string                               `json:"reputation_type"`
	EntityType     string                               `json:"entity_type"`
	UpdatedAt      time.Time                            `json:"updated_at"`
	Entries        []ReputationLeaderboardEntryResponse `json:"entries"`
	Pagination     Pagination                           `json:"pagination"`
}

type LeaderboardSummaryResponse struct {
	Ocean     string    `json:"ocean"`
	UpdatedAt time.Time `json:"updated_at"`
	
	TopCrewsByWins    []CrewLeaderboardEntryResponse `json:"top_crews_by_wins"`
	TopCrewsByWinRate []CrewLeaderboardEntryResponse `json:"top_crews_by_win_rate"`
	TopCrewsByRank    []CrewLeaderboardEntryResponse `json:"top_crews_by_rank"`
	TopFlagsByWins    []FlagLeaderboardEntryResponse `json:"top_flags_by_wins"`
	TopFlagsByWinRate []FlagLeaderboardEntryResponse `json:"top_flags_by_win_rate"`
}

type DailyCrewLeaderboardEntryResponse struct {
	Rank       int       `json:"rank"`
	CrewID     uint      `json:"crew_id"`
	Name       string    `json:"name"`
	FlagName   string    `json:"flag_name,omitempty"`
	DailyWins  int       `json:"daily_wins"`
	DailyLosses int      `json:"daily_losses"`
	DailyBattles int     `json:"daily_battles"`
	DailyWinRate float64 `json:"daily_win_rate"`
}

type DailyLeaderboardResponse struct {
	Ocean     string                              `json:"ocean"`
	Date      time.Time                           `json:"date"`
	Type      string                              `json:"type"`
	Entries   []DailyCrewLeaderboardEntryResponse `json:"entries"`
}

type LeaderboardTypesResponse struct {
	CrewLeaderboards []LeaderboardTypeInfo `json:"crew_leaderboards"`
	FlagLeaderboards []LeaderboardTypeInfo `json:"flag_leaderboards"`
}

type LeaderboardTypeInfo struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
}
