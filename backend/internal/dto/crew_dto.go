package dto

import "time"

// Request types
type CrewIDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type CrewGameIDParam struct {
	GameCrewID uint64 `uri:"game_crew_id" binding:"required"`
}

type CrewListRequest struct {
	OceanParam
	PaginationParams
	SortParams
	
	IsActive  *bool  `form:"is_active" binding:"omitempty"`
	FlagID    *uint  `form:"flag_id" binding:"omitempty,min=1"`
	CrewRank  string `form:"crew_rank" binding:"omitempty,oneof=Sailors 'Mostly Harmless' 'Scurvy Dogs' Scoundrels Blaggards 'Dread Pirates' 'Sea Lords' Imperials"`
	FameLevel string `form:"fame_level" binding:"omitempty"`
}

func (r *CrewListRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("name")
}

type CrewSearchRequest struct {
	OceanParam
	PaginationParams
	
	Query     string `form:"q" binding:"required,min=2,max=100"`
	IsActive  *bool  `form:"is_active" binding:"omitempty"`
}

func (r *CrewSearchRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type CrewStatsRequest struct {
	OceanParam
	PaginationParams
	SortParams
	
	MinBattles int    `form:"min_battles" binding:"omitempty,min=0"`
	MinWinRate float64 `form:"min_win_rate" binding:"omitempty,min=0,max=100"`
	CrewRank   string `form:"crew_rank" binding:"omitempty"`
	FlagID     *uint  `form:"flag_id" binding:"omitempty,min=1"`
}

func (r *CrewStatsRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("total_pvp_wins")
}

type CrewHistoryRequest struct {
	DateRangeParams
}

func (r *CrewHistoryRequest) SetDefaults() {
}

type CrewBattleRecordsRequest struct {
	DateRangeParams
	PaginationParams
}

func (r *CrewBattleRecordsRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type CrewFameHistoryRequest struct {
	DateRangeParams
	PaginationParams
}

func (r *CrewFameHistoryRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type CrewCompareRequest struct {
	Crew1ID uint `form:"crew1_id" binding:"required,min=1"`
	Crew2ID uint `form:"crew2_id" binding:"required,min=1,nefield=Crew1ID"`
}

type MultiCrewCompareRequest struct {
	CrewIDs []uint `form:"crew_ids" binding:"required,min=2,max=10,dive,min=1"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=wins win_rate battles rank"`
}

func (r *MultiCrewCompareRequest) SetDefaults() {
	if r.SortBy == "" {
		r.SortBy = "wins"
	}
}

type CrewTrendRequest struct {
	Period string `form:"period" binding:"omitempty,oneof=7d 30d 90d all"`
}

func (r *CrewTrendRequest) SetDefaults() {
	if r.Period == "" {
		r.Period = "30d"
	}
}

// Response types
type CrewBrief struct {
	ID         uint   `json:"id"`
	GameCrewID uint64 `json:"game_crew_id"`
	Name       string `json:"name"`
	Ocean      string `json:"ocean"`
}

type CrewResponse struct {
	ID          uint       `json:"id"`
	GameCrewID  uint64     `json:"game_crew_id"`
	Name        string     `json:"name"`
	Ocean       string     `json:"ocean"`
	IsActive    bool       `json:"is_active"`
	FirstSeenAt time.Time  `json:"first_seen_at"`
	LastSeenAt  time.Time  `json:"last_seen_at"`
	Flag        *FlagBrief `json:"flag,omitempty"`
	URLs        CrewURLs   `json:"urls"`
}

type CrewURLs struct {
	Info       string `json:"info"`
	BattleInfo string `json:"battle_info"`
}

type CrewDetailResponse struct {
	CrewResponse
	CurrentStats *CrewPVPStatsResponse `json:"current_stats,omitempty"`
	CurrentFame  *CrewFameResponse     `json:"current_fame,omitempty"`
}

type CrewListResponse struct {
	Crews      []CrewResponse `json:"crews"`
	Pagination Pagination     `json:"pagination"`
}

type CrewPVPStatsResponse struct {
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
	LastUpdated    time.Time `json:"last_updated"`
}

type CrewPVPStatsListResponse struct {
	Stats      []CrewPVPStatsResponse `json:"stats"`
	Pagination Pagination             `json:"pagination"`
}

type CrewBattleRecordResponse struct {
	ID             uint      `json:"id"`
	CrewID         uint      `json:"crew_id"`
	ScrapedAt      time.Time `json:"scraped_at"`
	CrewRank       string    `json:"crew_rank"`
	
	TotalPVPWins   int `json:"total_pvp_wins"`
	TotalPVPLosses int `json:"total_pvp_losses"`
	
	DailyPVPWins   int `json:"daily_pvp_wins"`
	DailyPVPLosses int `json:"daily_pvp_losses"`
	
	GreeterWins      int `json:"greeter_wins"`
	GreeterLosses    int `json:"greeter_losses"`
	BlockadeWins     int `json:"blockade_wins"`
	BlockadeLosses   int `json:"blockade_losses"`
	SeaMonsterWins   int `json:"sea_monster_wins"`
	SeaMonsterLosses int `json:"sea_monster_losses"`
	FlotillaWins     int `json:"flotilla_wins"`
	FlotillaLosses   int `json:"flotilla_losses"`
	
	WinRate      float64 `json:"win_rate"`
	TotalBattles int     `json:"total_battles"`
}

type CrewBattleRecordListResponse struct {
	CrewID  uint                       `json:"crew_id"`
	Records []CrewBattleRecordResponse `json:"records"`
}

type CrewFameResponse struct {
	CrewID    uint      `json:"crew_id"`
	ScrapedAt time.Time `json:"scraped_at"`
	FameLevel string    `json:"fame_level"`
	FameRank  *int      `json:"fame_rank,omitempty"`
}

type CrewFameHistoryResponse struct {
	CrewID  uint               `json:"crew_id"`
	History []CrewFameResponse `json:"history"`
}

type CrewReputationResponse struct {
	CrewID          uint      `json:"crew_id"`
	ScrapedAt       time.Time `json:"scraped_at"`
	ReputationType  string    `json:"reputation_type"`
	ReputationLevel string    `json:"reputation_level"`
	ReputationRank  *int      `json:"reputation_rank,omitempty"`
}

type CrewReputationSummaryResponse struct {
	CrewID      uint                     `json:"crew_id"`
	ScrapedAt   time.Time                `json:"scraped_at"`
	Reputations []CrewReputationResponse `json:"reputations"`
}

type CrewHistoryPointResponse struct {
	Date        time.Time `json:"date"`
	TotalWins   int       `json:"total_wins"`
	TotalLosses int       `json:"total_losses"`
	DailyWins   int       `json:"daily_wins"`
	DailyLosses int       `json:"daily_losses"`
	Rank        string    `json:"rank"`
	FameLevel   string    `json:"fame_level,omitempty"`
	WinRate     float64   `json:"win_rate"`
}

type CrewHistoricalStatsResponse struct {
	Crew        CrewBrief                  `json:"crew"`
	StartDate   time.Time                  `json:"start_date"`
	EndDate     time.Time                  `json:"end_date"`
	DataPoints  []CrewHistoryPointResponse `json:"data_points"`
	
	TotalWinsGained   int    `json:"total_wins_gained"`
	TotalLossesGained int    `json:"total_losses_gained"`
	StartingRank      string `json:"starting_rank"`
	EndingRank        string `json:"ending_rank"`
	BestRank          string `json:"best_rank"`
	WorstRank         string `json:"worst_rank"`
}

type DailyPVPChangeResponse struct {
	Date         time.Time `json:"date"`
	CrewID       uint      `json:"crew_id"`
	CrewName     string    `json:"crew_name"`
	WinsGained   int       `json:"wins_gained"`
	LossesGained int       `json:"losses_gained"`
	RankBefore   string    `json:"rank_before,omitempty"`
	RankAfter    string    `json:"rank_after,omitempty"`
	RankChanged  bool      `json:"rank_changed"`
}

type DailyActivityResponse struct {
	Ocean       string                   `json:"ocean"`
	Date        time.Time                `json:"date"`
	TotalWins   int                      `json:"total_wins"`
	TotalLosses int                      `json:"total_losses"`
	ActiveCrews int                      `json:"active_crews"`
	Changes     []DailyPVPChangeResponse `json:"changes"`
}

type CrewSearchResultResponse struct {
	CrewBrief
	FlagName  string `json:"flag_name,omitempty"`
	CrewRank  string `json:"crew_rank"`
	FameLevel string `json:"fame_level"`
	IsActive  bool   `json:"is_active"`
}

type CrewSearchResponse struct {
	Query      string                     `json:"query"`
	Results    []CrewSearchResultResponse `json:"results"`
	TotalCount int                        `json:"total_count"`
}
