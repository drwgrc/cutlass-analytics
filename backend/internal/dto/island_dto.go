package dto

import "time"


type IslandIDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type IslandGameIDParam struct {
	GameIslandID uint64 `uri:"game_island_id" binding:"required"`
}

type ArchipelagoIDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type IslandListRequest struct {
	OceanParam
	PaginationParams
	SortParams
	
	ArchipelagoID *uint  `form:"archipelago_id" binding:"omitempty,min=1"`
	Size          string `form:"size" binding:"omitempty,oneof=outpost medium large"`
	IsColonized   *bool  `form:"is_colonized" binding:"omitempty"`
	HasCommodity  string `form:"has_commodity" binding:"omitempty"`
	GovernorFlagID *uint `form:"governor_flag_id" binding:"omitempty,min=1"`
}

func (r *IslandListRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	r.SortParams.SetDefaults("name")
}

type IslandSearchRequest struct {
	OceanParam
	PaginationParams
	
	Query       string `form:"q" binding:"required,min=2,max=100"`
	IsColonized *bool  `form:"is_colonized" binding:"omitempty"`
}

func (r *IslandSearchRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type IslandPopulationHistoryRequest struct {
	DateRangeParams
}

type OceanPopulationRequest struct {
	OceanParam
	
	TopN int `form:"top_n" binding:"omitempty,min=1,max=50"`
}

func (r *OceanPopulationRequest) SetDefaults() {
	if r.TopN == 0 {
		r.TopN = 10
	}
}

type IslandBuildingListRequest struct {
	PaginationParams
	
	BuildingType string `form:"building_type" binding:"omitempty"`
	IsActive     *bool  `form:"is_active" binding:"omitempty"`
}

func (r *IslandBuildingListRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}

type RentComparisonRequest struct {
	OceanParam
	
	BuildingType string `form:"building_type" binding:"required"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=stall_rent shoppe_rent island_name"`
}

func (r *RentComparisonRequest) SetDefaults() {
	if r.SortBy == "" {
		r.SortBy = "stall_rent"
	}
}

type ArchipelagoListRequest struct {
	OceanParam
	
	IncludeIslands bool `form:"include_islands" binding:"omitempty"`
}

type IslandGovernanceHistoryRequest struct {
	DateRangeParams
	PaginationParams
}

func (r *IslandGovernanceHistoryRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
}


// Response types
type IslandBrief struct {
	ID           uint   `json:"id"`
	GameIslandID uint64 `json:"game_island_id"`
	Name         string `json:"name"`
	Ocean        string `json:"ocean"`
	IsColonized  bool   `json:"is_colonized"`
}

type ArchipelagoBrief struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

type IslandResponse struct {
	ID           uint              `json:"id"`
	GameIslandID uint64            `json:"game_island_id"`
	Name         string            `json:"name"`
	Ocean        string            `json:"ocean"`
	Size         string            `json:"size"` // outpost, medium, large
	IsColonized  bool              `json:"is_colonized"`
	Population   int               `json:"population"`
	Archipelago  *ArchipelagoBrief `json:"archipelago,omitempty"`
	Governor     *IslandGovernor   `json:"governor,omitempty"`
	FirstSeenAt  time.Time         `json:"first_seen_at"`
	LastSeenAt   time.Time         `json:"last_seen_at"`
	URL          string            `json:"url"`
}

type IslandGovernor struct {
	FlagID       *uint  `json:"flag_id,omitempty"`
	FlagName     string `json:"flag_name,omitempty"`
	GovernorName string `json:"governor_name,omitempty"`
}

type IslandDetailResponse struct {
	IslandResponse
	Commodities []CommodityBrief        `json:"commodities,omitempty"`
	Buildings   []IslandBuildingResponse `json:"buildings,omitempty"`
	TaxSettings *IslandTaxSettingResponse `json:"tax_settings,omitempty"`
	RentPrices  []ShoppeRentPriceResponse `json:"rent_prices,omitempty"`
}

type IslandListResponse struct {
	Islands    []IslandResponse `json:"islands"`
	Pagination Pagination       `json:"pagination"`
}

type IslandBuildingResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	BuildingType string `json:"building_type"`
	OwnerName    string `json:"owner_name,omitempty"`
	ShoppeClass  string `json:"shoppe_class,omitempty"`
	IsActive     bool   `json:"is_active"`
}

type IslandBuildingListResponse struct {
	IslandID   uint                     `json:"island_id"`
	IslandName string                   `json:"island_name"`
	Buildings  []IslandBuildingResponse `json:"buildings"`
	
	Summary BuildingSummary `json:"summary"`
}

type BuildingSummary struct {
	TotalBuildings    int `json:"total_buildings"`
	InfrastructureCount int `json:"infrastructure_count"`
	ShoppeCount       int `json:"shoppe_count"`
	BazaarCount       int `json:"bazaar_count"`
	HousingCount      int `json:"housing_count"`
}

type IslandTaxSettingResponse struct {
	IslandID     uint      `json:"island_id"`
	ScrapedAt    time.Time `json:"scraped_at"`
	ShoppeTax    *float64  `json:"shoppe_tax,omitempty"`
	StallTax     *float64  `json:"stall_tax,omitempty"`
	HousingTax   *float64  `json:"housing_tax,omitempty"`
	CommodityTax *float64  `json:"commodity_tax,omitempty"`
	DockingFee   *int      `json:"docking_fee,omitempty"`
}

type ShoppeRentPriceResponse struct {
	BuildingType string    `json:"building_type"`
	StallRent    *int      `json:"stall_rent,omitempty"`
	ShoppeRent   *int      `json:"shoppe_rent,omitempty"`
	ScrapedAt    time.Time `json:"scraped_at"`
}

type IslandRentComparisonResponse struct {
	BuildingType string                       `json:"building_type"`
	Ocean        string                       `json:"ocean"`
	Islands      []IslandRentComparisonEntry  `json:"islands"`
	UpdatedAt    time.Time                    `json:"updated_at"`
}

type IslandRentComparisonEntry struct {
	Island     IslandBrief `json:"island"`
	StallRent  *int        `json:"stall_rent,omitempty"`
	ShoppeRent *int        `json:"shoppe_rent,omitempty"`
}

type ArchipelagoResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	Ocean       string           `json:"ocean"`
	Color       string           `json:"color,omitempty"`
	IslandCount int              `json:"island_count"`
	Islands     []IslandBrief    `json:"islands,omitempty"`
}

type ArchipelagoListResponse struct {
	Archipelagos []ArchipelagoResponse `json:"archipelagos"`
}

type IslandPopulationResponse struct {
	IslandID   uint      `json:"island_id"`
	ScrapedAt  time.Time `json:"scraped_at"`
	Population int       `json:"population"`
}

type IslandPopulationHistoryResponse struct {
	Island     IslandBrief                `json:"island"`
	StartDate  time.Time                  `json:"start_date"`
	EndDate    time.Time                  `json:"end_date"`
	DataPoints []IslandPopulationResponse `json:"data_points"`
	
	MinPopulation int `json:"min_population"`
	MaxPopulation int `json:"max_population"`
	AvgPopulation int `json:"avg_population"`
}

type OceanPopulationResponse struct {
	Ocean           string                        `json:"ocean"`
	TotalPopulation int                           `json:"total_population"`
	IslandCount     int                           `json:"island_count"`
	ScrapedAt       time.Time                     `json:"scraped_at"`
	TopIslands      []IslandPopulationRankEntry   `json:"top_islands"`
}

type IslandPopulationRankEntry struct {
	Rank       int         `json:"rank"`
	Island     IslandBrief `json:"island"`
	Population int         `json:"population"`
	Percentage float64     `json:"percentage"`
}

type IslandGovernanceHistoryResponse struct {
	Island  IslandBrief                 `json:"island"`
	History []GovernanceChangeResponse  `json:"history"`
}

type GovernanceChangeResponse struct {
	FlagID       *uint      `json:"flag_id,omitempty"`
	FlagName     string     `json:"flag_name,omitempty"`
	GovernorName string     `json:"governor_name,omitempty"`
	StartedAt    time.Time  `json:"started_at"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	ChangeType   string     `json:"change_type,omitempty"`
	IsCurrent    bool       `json:"is_current"`
}

type IslandSearchResultResponse struct {
	IslandBrief
	ArchipelagoName string   `json:"archipelago_name,omitempty"`
	Size            string   `json:"size"`
	Population      int      `json:"population"`
	Commodities     []string `json:"commodities,omitempty"`
}

type IslandSearchResponse struct {
	Query      string                       `json:"query"`
	Results    []IslandSearchResultResponse `json:"results"`
	TotalCount int                          `json:"total_count"`
}