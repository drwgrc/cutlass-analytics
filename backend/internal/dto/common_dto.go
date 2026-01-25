package dto

import "time"

// Response types
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Meta struct {
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Pagination struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
	Services  map[string]string `json:"services,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type DeleteResponse struct {
	Deleted bool   `json:"deleted"`
	ID      uint   `json:"id"`
	Message string `json:"message,omitempty"`
}


// Request types
type PaginationParams struct {
	Page    int `form:"page" binding:"omitempty,min=1"`
	PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}

func (p *PaginationParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = 25
	}
}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *PaginationParams) Limit() int {
	return p.PerPage
}

type DateRangeParams struct {
	StartDate string `form:"start_date" binding:"omitempty"`
	EndDate   string `form:"end_date" binding:"omitempty"`
}

func (d *DateRangeParams) ParsedStartDate() (time.Time, error) {
	if d.StartDate == "" {
		return time.Now().AddDate(0, 0, -30), nil
	}
	return time.Parse("2006-01-02", d.StartDate)
}

func (d *DateRangeParams) ParsedEndDate() (time.Time, error) {
	if d.EndDate == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", d.EndDate)
}

func (d *DateRangeParams) Validate() error {
	start, err := d.ParsedStartDate()
	if err != nil {
		return err
	}
	end, err := d.ParsedEndDate()
	if err != nil {
		return err
	}
	if start.After(end) {
		return ErrStartDateAfterEndDate
	}
	return nil
}

type SortParams struct {
	SortBy    string `form:"sort_by" binding:"omitempty"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

func (s *SortParams) SetDefaults(defaultSortBy string) {
	if s.SortBy == "" {
		s.SortBy = defaultSortBy
	}
	if s.SortOrder == "" {
		s.SortOrder = "desc"
	}
}

func (s *SortParams) IsAscending() bool {
	return s.SortOrder == "asc"
}

type OceanParam struct {
	Ocean string `form:"ocean" binding:"required,oneof=emerald meridian cerulean obsidian"`
}

type OceanPathParam struct {
	Ocean string `uri:"ocean" binding:"required,oneof=emerald meridian cerulean obsidian"`
}

var (
	ErrStartDateAfterEndDate = &ValidationError{Field: "start_date", Message: "start_date must be before end_date"}
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
