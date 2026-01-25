package dto

import "time"

type CrewComparisonResponse struct {
	Crew1 CrewComparisonData `json:"crew1"`
	Crew2 CrewComparisonData `json:"crew2"`
	
	WinsLeader     string `json:"wins_leader"`
	WinRateLeader  string `json:"win_rate_leader"`
	BattlesLeader  string `json:"battles_leader"`
	RankLeader     string `json:"rank_leader"`
	FameLeader     string `json:"fame_leader"`
	
	ComparedAt time.Time `json:"compared_at"`
}

type CrewComparisonData struct {
	Crew           CrewBrief `json:"crew"`
	FlagName       string    `json:"flag_name,omitempty"`
	CrewRank       string    `json:"crew_rank"`
	FameLevel      string    `json:"fame_level"`
	TotalPVPWins   int       `json:"total_pvp_wins"`
	TotalPVPLosses int       `json:"total_pvp_losses"`
	WinRate        float64   `json:"win_rate"`
	TotalBattles   int       `json:"total_battles"`
}

type FlagComparisonResponse struct {
	Flag1 FlagComparisonData `json:"flag1"`
	Flag2 FlagComparisonData `json:"flag2"`
	
	WinsLeader     string `json:"wins_leader"`
	WinRateLeader  string `json:"win_rate_leader"`
	CrewCountLeader string `json:"crew_count_leader"`
	FameLeader     string `json:"fame_leader"`
	
	ComparedAt time.Time `json:"compared_at"`
}

type FlagComparisonData struct {
	Flag           FlagBrief `json:"flag"`
	FameLevel      string    `json:"fame_level"`
	CrewCount      int       `json:"crew_count"`
	TotalPVPWins   int       `json:"total_pvp_wins"`
	TotalPVPLosses int       `json:"total_pvp_losses"`
	WinRate        float64   `json:"win_rate"`
	TotalBattles   int       `json:"total_battles"`
}

type MultiCrewComparisonResponse struct {
	Crews      []CrewComparisonData `json:"crews"`
	SortedBy   string               `json:"sorted_by"`
	ComparedAt time.Time            `json:"compared_at"`
}

type MultiFlagComparisonResponse struct {
	Flags      []FlagComparisonData `json:"flags"`
	SortedBy   string               `json:"sorted_by"`
	ComparedAt time.Time            `json:"compared_at"`
}

type CrewTrendResponse struct {
	Crew       CrewBrief       `json:"crew"`
	Period     string          `json:"period"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    time.Time       `json:"end_date"`
	
	WinsTrend       []TrendPoint `json:"wins_trend"`
	WinRateTrend    []TrendPoint `json:"win_rate_trend"`
	BattlesTrend    []TrendPoint `json:"battles_trend"`
	
	TotalWinsGained   int     `json:"total_wins_gained"`
	TotalLossesGained int     `json:"total_losses_gained"`
	WinRateChange     float64 `json:"win_rate_change"`
	RankChanges       int     `json:"rank_changes"`
}

type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type FlagTrendResponse struct {
	Flag       FlagBrief       `json:"flag"`
	Period     string          `json:"period"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    time.Time       `json:"end_date"`
	
	WinsTrend      []TrendPoint `json:"wins_trend"`
	WinRateTrend   []TrendPoint `json:"win_rate_trend"`
	CrewCountTrend []TrendPoint `json:"crew_count_trend"`
	
	TotalWinsGained   int     `json:"total_wins_gained"`
	TotalLossesGained int     `json:"total_losses_gained"`
	WinRateChange     float64 `json:"win_rate_change"`
	NetCrewChange     int     `json:"net_crew_change"`
}

type OceanActivityResponse struct {
	Ocean          string           `json:"ocean"`
	Period         string           `json:"period"`
	StartDate      time.Time        `json:"start_date"`
	EndDate        time.Time        `json:"end_date"`
	
	DailyActivity  []DailyActivitySummary `json:"daily_activity"`
	
	TotalBattles   int     `json:"total_battles"`
	TotalWins      int     `json:"total_wins"`
	AvgDailyBattles float64 `json:"avg_daily_battles"`
	MostActiveDay  string  `json:"most_active_day"`
	LeastActiveDay string  `json:"least_active_day"`
}

type DailyActivitySummary struct {
	Date         time.Time `json:"date"`
	TotalBattles int       `json:"total_battles"`
	TotalWins    int       `json:"total_wins"`
	TotalLosses  int       `json:"total_losses"`
	ActiveCrews  int       `json:"active_crews"`
}

type RankChangeResponse struct {
	Crew       CrewBrief         `json:"crew"`
	Period     string            `json:"period"`
	Changes    []RankChangeEvent `json:"changes"`
	TotalChanges int             `json:"total_changes"`
}

type RankChangeEvent struct {
	Date     time.Time `json:"date"`
	FromRank string    `json:"from_rank"`
	ToRank   string    `json:"to_rank"`
	IsPromotion bool   `json:"is_promotion"`
}

type RecentRankChangesResponse struct {
	Ocean   string                    `json:"ocean"`
	Period  string                    `json:"period"`
	Changes []CrewRankChangeResponse  `json:"changes"`
}

type CrewRankChangeResponse struct {
	Crew        CrewBrief `json:"crew"`
	FlagName    string    `json:"flag_name,omitempty"`
	Date        time.Time `json:"date"`
	FromRank    string    `json:"from_rank"`
	ToRank      string    `json:"to_rank"`
	IsPromotion bool      `json:"is_promotion"`
}
