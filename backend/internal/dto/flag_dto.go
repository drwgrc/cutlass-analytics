package dto

import "time"

// Request types
type FlagIDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type FlagGameIDParam struct {
	GameFlagID uint64 `uri:"game_flag_id" binding:"required"`
}

type FlagListRequest struct {
	OceanParam
	PaginationParams
	SortParams
	
	IsActive  *bool  `form:"is_active" binding:"omitempty"`
	FameLevel string `form:"fame_level" binding:"omitempty"`
	MinCrews  int    `form:"min_crews" binding:"omitempty,min=0"`
}

func (r *FlagListRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("name")
}

type FlagSearchRequest struct {
	OceanParam
	PaginationParams
	
	Query    string `form:"q" binding:"required,min=2,max=100"`
	IsActive *bool  `form:"is_active" binding:"omitempty"`
}

func (r *FlagSearchRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type FlagStatsRequest struct {
	OceanParam
	PaginationParams
	SortParams
	
	MinBattles int     `form:"min_battles" binding:"omitempty,min=0"`
	MinWinRate float64 `form:"min_win_rate" binding:"omitempty,min=0,max=100"`
	MinCrews   int     `form:"min_crews" binding:"omitempty,min=0"`
	FameLevel  string  `form:"fame_level" binding:"omitempty"`
}

func (r *FlagStatsRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("total_pvp_wins")
}

type FlagCrewsRequest struct {
	PaginationParams
	SortParams
	
	IsActive *bool `form:"is_active" binding:"omitempty"`
}

func (r *FlagCrewsRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("total_pvp_wins")
}

type FlagHistoryRequest struct {
	DateRangeParams
}

type FlagFameHistoryRequest struct {
	DateRangeParams
	PaginationParams
}

func (r *FlagFameHistoryRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type FlagMembershipHistoryRequest struct {
	DateRangeParams
	PaginationParams
	
	IncludeJoins  bool `form:"include_joins" binding:"omitempty"`
	IncludeLeaves bool `form:"include_leaves" binding:"omitempty"`
}

func (r *FlagMembershipHistoryRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.IncludeJoins = true
	r.IncludeLeaves = true
}

type FlagCompareRequest struct {
	Flag1ID uint `form:"flag1_id" binding:"required,min=1"`
	Flag2ID uint `form:"flag2_id" binding:"required,min=1,nefield=Flag1ID"`
}

type MultiFlagCompareRequest struct {
	FlagIDs []uint `form:"flag_ids" binding:"required,min=2,max=10,dive,min=1"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=wins win_rate crews fame"`
}

func (r *MultiFlagCompareRequest) SetDefaults() {
	if r.SortBy == "" {
		r.SortBy = "wins"
	}
}

type FlagTrendRequest struct {
	Period string `form:"period" binding:"omitempty,oneof=7d 30d 90d all"`
}

func (r *FlagTrendRequest) SetDefaults() {
	if r.Period == "" {
		r.Period = "30d"
	}
}

// Response types
type FlagBrief struct {
	ID         uint   `json:"id"`
	GameFlagID uint64 `json:"game_flag_id"`
	Name       string `json:"name"`
	Ocean      string `json:"ocean"`
}

type FlagResponse struct {
	ID          uint      `json:"id"`
	GameFlagID  uint64    `json:"game_flag_id"`
	Name        string    `json:"name"`
	Ocean       string    `json:"ocean"`
	IsActive    bool      `json:"is_active"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	URL         string    `json:"url"`
}

type FlagDetailResponse struct {
	FlagResponse
	CrewCount    int                    `json:"crew_count"`
	CurrentFame  *FlagFameResponse      `json:"current_fame,omitempty"`
	CurrentStats *FlagPVPStatsResponse  `json:"current_stats,omitempty"`
	Crews        []CrewPVPStatsResponse `json:"crews,omitempty"`
}

type FlagListResponse struct {
	Flags      []FlagResponse `json:"flags"`
	Pagination Pagination     `json:"pagination"`
}

type FlagPVPStatsResponse struct {
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
	LastUpdated    time.Time `json:"last_updated"`
}

type FlagPVPStatsListResponse struct {
	Stats      []FlagPVPStatsResponse `json:"stats"`
	Pagination Pagination             `json:"pagination"`
}

type FlagFameResponse struct {
	FlagID    uint      `json:"flag_id"`
	ScrapedAt time.Time `json:"scraped_at"`
	FameLevel string    `json:"fame_level"`
	FameRank  *int      `json:"fame_rank,omitempty"`
}

type FlagFameHistoryResponse struct {
	FlagID  uint               `json:"flag_id"`
	History []FlagFameResponse `json:"history"`
}

type FlagMembershipLogResponse struct {
	Flag     *FlagBrief `json:"flag,omitempty"`
	JoinedAt time.Time  `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at,omitempty"`
	Days     int        `json:"days"`
	IsCurrent bool      `json:"is_current"`
}

type CrewFlagHistoryResponse struct {
	Crew        CrewBrief                   `json:"crew"`
	CurrentFlag *FlagBrief                  `json:"current_flag,omitempty"`
	History     []FlagMembershipLogResponse `json:"history"`
}

type FlagMembershipResponse struct {
	Flag          FlagBrief   `json:"flag"`
	CurrentCrews  []CrewBrief `json:"current_crews"`
	TotalCrews    int         `json:"total_crews"`
	RecentJoins   []CrewFlagChangeResponse `json:"recent_joins,omitempty"`
	RecentLeaves  []CrewFlagChangeResponse `json:"recent_leaves,omitempty"`
}

type CrewFlagChangeResponse struct {
	Crew      CrewBrief `json:"crew"`
	Timestamp time.Time `json:"timestamp"`
}

type FlagHistoryPointResponse struct {
	Date        time.Time `json:"date"`
	CrewCount   int       `json:"crew_count"`
	TotalWins   int       `json:"total_wins"`
	TotalLosses int       `json:"total_losses"`
	FameLevel   string    `json:"fame_level"`
	FameRank    *int      `json:"fame_rank,omitempty"`
	WinRate     float64   `json:"win_rate"`
}

type FlagHistoricalStatsResponse struct {
	Flag       FlagBrief                  `json:"flag"`
	StartDate  time.Time                  `json:"start_date"`
	EndDate    time.Time                  `json:"end_date"`
	DataPoints []FlagHistoryPointResponse `json:"data_points"`
	
	CrewsJoined int `json:"crews_joined"`
	CrewsLeft   int `json:"crews_left"`
	NetCrewChange int `json:"net_crew_change"`
}

type FlagSearchResultResponse struct {
	FlagBrief
	CrewCount int    `json:"crew_count"`
	FameLevel string `json:"fame_level"`
	IsActive  bool   `json:"is_active"`
}

type FlagSearchResponse struct {
	Query      string                     `json:"query"`
	Results    []FlagSearchResultResponse `json:"results"`
	TotalCount int                        `json:"total_count"`
}
