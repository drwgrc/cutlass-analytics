package dto

type DailyActivityRequest struct {
	OceanParam
	
	Date string `form:"date" binding:"omitempty"`
	
	IncludeChanges bool `form:"include_changes" binding:"omitempty"`
}

type DailyChangesRequest struct {
	OceanParam
	PaginationParams
	
	Date string `form:"date" binding:"omitempty"`
	
	MinWins   int `form:"min_wins" binding:"omitempty,min=0"`
	MinLosses int `form:"min_losses" binding:"omitempty,min=0"`
}

func (r *DailyChangesRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type OceanActivityRequest struct {
	OceanParam
	DateRangeParams
}

type ActivitySummaryRequest struct {
	OceanParam
	
	Period string `form:"period" binding:"omitempty,oneof=7d 30d 90d"`
}

func (r *ActivitySummaryRequest) SetDefaults() {
	if r.Period == "" {
		r.Period = "7d"
	}
}

type RankChangesRequest struct {
	OceanParam
	PaginationParams
	DateRangeParams
	
	PromotionsOnly bool `form:"promotions_only" binding:"omitempty"`
	DemotionsOnly  bool `form:"demotions_only" binding:"omitempty"`
}

func (r *RankChangesRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type CrewRankChangesRequest struct {
	DateRangeParams
}

type TrendingCrewsRequest struct {
	OceanParam
	PaginationParams
	
	Period string `form:"period" binding:"omitempty,oneof=24h 7d 30d"`
	
	Metric string `form:"metric" binding:"omitempty,oneof=wins battles win_rate"`
}

func (r *TrendingCrewsRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.Period == "" {
		r.Period = "7d"
	}
	if r.Metric == "" {
		r.Metric = "wins"
	}
}

type TrendingFlagsRequest struct {
	OceanParam
	PaginationParams
	
	Period string `form:"period" binding:"omitempty,oneof=24h 7d 30d"`
	
	Metric string `form:"metric" binding:"omitempty,oneof=wins battles win_rate crew_count"`
}

func (r *TrendingFlagsRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.Period == "" {
		r.Period = "7d"
	}
	if r.Metric == "" {
		r.Metric = "wins"
	}
}
