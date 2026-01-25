package dto

import "time"

// Request types
type ScrapeJobIDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type ScrapeJobListRequest struct {
	PaginationParams
	SortParams
	
	Ocean   string `form:"ocean" binding:"omitempty,oneof=emerald meridian cerulean obsidian"`
	JobType string `form:"job_type" binding:"omitempty,oneof=crew_fame flag_fame crew_info battle_info daily_full"`
	Status  string `form:"status" binding:"omitempty,oneof=running completed failed"`
}

func (r *ScrapeJobListRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("started_at")
}

type TriggerScrapeRequest struct {
	Ocean   string `json:"ocean" binding:"required,oneof=emerald meridian cerulean obsidian"`
	JobType string `json:"job_type" binding:"omitempty,oneof=crew_fame flag_fame crew_info battle_info daily_full"`
}

func (r *TriggerScrapeRequest) SetDefaults() {
	if r.JobType == "" {
		r.JobType = "daily_full"
	}
}

type ScrapeHistoryRequest struct {
	OceanParam
	DateRangeParams
	PaginationParams
}

func (r *ScrapeHistoryRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type UpdateScrapeScheduleRequest struct {
	Ocean     string `json:"ocean" binding:"required,oneof=emerald meridian cerulean obsidian"`
	Frequency string `json:"frequency" binding:"required,oneof=hourly daily weekly"`
	IsEnabled bool   `json:"is_enabled"`
	
	Time string `json:"time" binding:"omitempty"`
}

// Response types
type ScrapeJobResponse struct {
	ID             uint       `json:"id"`
	Ocean          string     `json:"ocean"`
	JobType        string     `json:"job_type"`
	Status         string     `json:"status"`
	StartedAt      time.Time  `json:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	Duration       string     `json:"duration,omitempty"`
	DurationMs     int64      `json:"duration_ms,omitempty"`
	ItemsProcessed int        `json:"items_processed"`
	ItemsFailed    int        `json:"items_failed"`
	SuccessRate    float64    `json:"success_rate"`
	ErrorMessage   string     `json:"error_message,omitempty"`
}

type ScrapeJobListResponse struct {
	Jobs       []ScrapeJobResponse `json:"jobs"`
	Pagination Pagination          `json:"pagination"`
}

type ScrapeJobDetailResponse struct {
	ScrapeJobResponse
	
	CrewsProcessed  int `json:"crews_processed,omitempty"`
	CrewsFailed     int `json:"crews_failed,omitempty"`
	FlagsProcessed  int `json:"flags_processed,omitempty"`
	FlagsFailed     int `json:"flags_failed,omitempty"`
	RecordsCreated  int `json:"records_created,omitempty"`
	RecordsUpdated  int `json:"records_updated,omitempty"`
}

type ScrapeStatusResponse struct {
	IsRunning     bool               `json:"is_running"`
	CurrentJob    *ScrapeJobResponse `json:"current_job,omitempty"`
	LastCompleted *ScrapeJobResponse `json:"last_completed,omitempty"`
	NextScheduled *time.Time         `json:"next_scheduled,omitempty"`
}

type ScrapeScheduleResponse struct {
	Ocean         string     `json:"ocean"`
	Frequency     string     `json:"frequency"`
	NextRun       *time.Time `json:"next_run,omitempty"`
	LastRun       *time.Time `json:"last_run,omitempty"`
	LastStatus    string     `json:"last_status,omitempty"`
	IsEnabled     bool       `json:"is_enabled"`
}

type ScrapeScheduleListResponse struct {
	Schedules []ScrapeScheduleResponse `json:"schedules"`
}

type ScrapeTriggerResponse struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	JobID     uint              `json:"job_id"`
	Job       ScrapeJobResponse `json:"job"`
}

type ScrapeHistoryResponse struct {
	Ocean           string              `json:"ocean"`
	TotalJobs       int                 `json:"total_jobs"`
	SuccessfulJobs  int                 `json:"successful_jobs"`
	FailedJobs      int                 `json:"failed_jobs"`
	AvgDuration     string              `json:"avg_duration"`
	LastWeekJobs    []ScrapeJobResponse `json:"last_week_jobs"`
}

type ScrapeStatsResponse struct {
	TotalJobsToday      int       `json:"total_jobs_today"`
	TotalJobsThisWeek   int       `json:"total_jobs_this_week"`
	SuccessRateToday    float64   `json:"success_rate_today"`
	SuccessRateThisWeek float64   `json:"success_rate_this_week"`
	TotalItemsScraped   int       `json:"total_items_scraped"`
	LastSuccessfulScrape time.Time `json:"last_successful_scrape"`
	
	OceanStats []OceanScrapeStats `json:"ocean_stats"`
}

type OceanScrapeStats struct {
	Ocean            string    `json:"ocean"`
	LastScrapedAt    time.Time `json:"last_scraped_at"`
	LastStatus       string    `json:"last_status"`
	ItemsLastScraped int       `json:"items_last_scraped"`
}
