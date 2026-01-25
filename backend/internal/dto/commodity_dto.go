package dto

import "time"

// Request types
type CommodityIDParam struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type CommodityNameParam struct {
	Name string `uri:"name" binding:"required,min=1"`
}

type CommodityListRequest struct {
	Category    string `form:"category" binding:"omitempty,oneof=basic herb mineral foraged refined ship_supply"`
	IsSpawnable *bool  `form:"is_spawnable" binding:"omitempty"`
	IsRare      *bool  `form:"is_rare" binding:"omitempty"`
	
	GroupByCategory bool `form:"group_by_category" binding:"omitempty"`
}

type CommoditySearchRequest struct {
	Query string `form:"q" binding:"required,min=2,max=100"`
}

type TaxRatesRequest struct {
	OceanParam
	
	Category string `form:"category" binding:"omitempty,oneof=basic herb mineral foraged refined ship_supply"`
	
	GroupByCategory bool `form:"group_by_category" binding:"omitempty"`
}

type TaxRateComparisonRequest struct {
	CommodityID uint `form:"commodity_id" binding:"required,min=1"`
}

type TaxRateHistoryRequest struct {
	OceanParam
	DateRangeParams
}

type MarketPricesRequest struct {
	OceanParam
	PaginationParams
	
	IslandID *uint `form:"island_id" binding:"omitempty,min=1"`
	
	CommodityID *uint `form:"commodity_id" binding:"omitempty,min=1"`
	
	HasBuyOffer  *bool `form:"has_buy_offer" binding:"omitempty"`
	HasSellOffer *bool `form:"has_sell_offer" binding:"omitempty"`
	
	SortBy string `form:"sort_by" binding:"omitempty,oneof=buy_price sell_price spread commodity island"`
}

func (r *MarketPricesRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.SortBy == "" {
		r.SortBy = "commodity"
	}
}

type IslandMarketPricesRequest struct {
	Category string `form:"category" binding:"omitempty,oneof=basic herb mineral foraged refined ship_supply"`
}

type CommodityMarketPricesRequest struct {
	OceanParam
	
	HasStock *bool `form:"has_stock" binding:"omitempty"`
	
	SortBy string `form:"sort_by" binding:"omitempty,oneof=buy_price sell_price island"`
}

func (r *CommodityMarketPricesRequest) SetDefaults() {
	if r.SortBy == "" {
		r.SortBy = "sell_price"
	}
}

type CommoditySpawnsRequest struct {
	OceanParam
	
	ConfirmedOnly bool `form:"confirmed_only" binding:"omitempty"`
}

type IslandSpawnsRequest struct {
	ConfirmedOnly bool `form:"confirmed_only" binding:"omitempty"`
}

type TradeRouteRequest struct {
	OceanParam
	PaginationParams
	
	CommodityID   *uint `form:"commodity_id" binding:"omitempty,min=1"`
	MinProfit     *int  `form:"min_profit" binding:"omitempty,min=0"`
	MinQuantity   *int  `form:"min_quantity" binding:"omitempty,min=1"`
	
	SortBy string `form:"sort_by" binding:"omitempty,oneof=profit max_profit quantity"`
}

func (r *TradeRouteRequest) SetDefaults() {
	r.PaginationParams.SetDefaults()
	if r.SortBy == "" {
		r.SortBy = "max_profit"
	}
}

type EconomySummaryRequest struct {
	OceanParam
	
	TopN int `form:"top_n" binding:"omitempty,min=1,max=25"`
}

func (r *EconomySummaryRequest) SetDefaults() {
	if r.TopN == 0 {
		r.TopN = 10
	}
}

// Response types
type CommodityBrief struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
}

type CommodityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
	IsSpawnable bool   `json:"is_spawnable"`
	IsRare      bool   `json:"is_rare"`
	Description string `json:"description,omitempty"`
	IconPath    string `json:"icon_path,omitempty"`
}

type CommodityDetailResponse struct {
	CommodityResponse
	SpawnIslands []IslandBrief          `json:"spawn_islands,omitempty"`
	TaxRates     []CommodityTaxRateResponse `json:"tax_rates,omitempty"`
}

type CommodityListResponse struct {
	Commodities []CommodityResponse `json:"commodities"`
	
	ByCategory map[string][]CommodityResponse `json:"by_category,omitempty"`
}

type CommodityTaxRateResponse struct {
	CommodityID   uint      `json:"commodity_id"`
	CommodityName string    `json:"commodity_name"`
	Ocean         string    `json:"ocean"`
	TaxValue      int       `json:"tax_value"`
	ScrapedAt     time.Time `json:"scraped_at"`
}

type OceanTaxRatesResponse struct {
	Ocean     string                     `json:"ocean"`
	ScrapedAt time.Time                  `json:"scraped_at"`
	TaxRates  []CommodityTaxRateResponse `json:"tax_rates"`
	
	ByCategory map[string][]CommodityTaxRateResponse `json:"by_category,omitempty"`
}

type TaxRateComparisonResponse struct {
	Commodity   CommodityBrief             `json:"commodity"`
	OceanRates  []OceanTaxRateEntry        `json:"ocean_rates"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

type OceanTaxRateEntry struct {
	Ocean    string `json:"ocean"`
	TaxValue int    `json:"tax_value"`
}

type TaxRateHistoryResponse struct {
	Commodity  CommodityBrief               `json:"commodity"`
	Ocean      string                       `json:"ocean"`
	StartDate  time.Time                    `json:"start_date"`
	EndDate    time.Time                    `json:"end_date"`
	DataPoints []TaxRateHistoryPoint        `json:"data_points"`
}

type TaxRateHistoryPoint struct {
	Date     time.Time `json:"date"`
	TaxValue int       `json:"tax_value"`
}

type MarketPriceResponse struct {
	IslandID      uint      `json:"island_id"`
	IslandName    string    `json:"island_name"`
	CommodityID   uint      `json:"commodity_id"`
	CommodityName string    `json:"commodity_name"`
	ScrapedAt     time.Time `json:"scraped_at"`
	
	BuyPrice    *int `json:"buy_price,omitempty"`
	BuyQuantity *int `json:"buy_quantity,omitempty"`
	
	SellPrice    *int `json:"sell_price,omitempty"`
	SellQuantity *int `json:"sell_quantity,omitempty"`
	
	Spread *int `json:"spread,omitempty"`
}

type IslandMarketPricesResponse struct {
	Island    IslandBrief           `json:"island"`
	ScrapedAt time.Time             `json:"scraped_at"`
	Prices    []MarketPriceResponse `json:"prices"`
}

type CommodityMarketPricesResponse struct {
	Commodity CommodityBrief        `json:"commodity"`
	Ocean     string                `json:"ocean"`
	ScrapedAt time.Time             `json:"scraped_at"`
	Prices    []MarketPriceResponse `json:"prices"`
	
	LowestSellPrice  *MarketPriceSummaryEntry `json:"lowest_sell_price,omitempty"`
	HighestBuyPrice  *MarketPriceSummaryEntry `json:"highest_buy_price,omitempty"`
}

type MarketPriceSummaryEntry struct {
	Island   IslandBrief `json:"island"`
	Price    int         `json:"price"`
	Quantity int         `json:"quantity"`
}

type CommoditySpawnResponse struct {
	Commodity CommodityBrief `json:"commodity"`
	Islands   []IslandSpawnEntry `json:"islands"`
}

type IslandSpawnEntry struct {
	Island      IslandBrief `json:"island"`
	Archipelago string      `json:"archipelago,omitempty"`
	IsConfirmed bool        `json:"is_confirmed"`
}

type IslandCommoditiesResponse struct {
	Island      IslandBrief      `json:"island"`
	Commodities []CommoditySpawnInfo `json:"commodities"`
}

type CommoditySpawnInfo struct {
	Commodity   CommodityBrief `json:"commodity"`
	IsConfirmed bool           `json:"is_confirmed"`
}

type TradeRouteResponse struct {
	Commodity   CommodityBrief `json:"commodity"`
	BuyIsland   IslandBrief    `json:"buy_island"`
	SellIsland  IslandBrief    `json:"sell_island"`
	BuyPrice    int            `json:"buy_price"`
	SellPrice   int            `json:"sell_price"`
	Profit      int            `json:"profit"` // per unit
	BuyQuantity int            `json:"buy_quantity"`
	SellQuantity int           `json:"sell_quantity"`
	MaxQuantity int            `json:"max_quantity"` // min of buy/sell quantity
	MaxProfit   int            `json:"max_profit"`   // profit * max_quantity
	ScrapedAt   time.Time      `json:"scraped_at"`
}

type TradeRouteListResponse struct {
	Ocean      string               `json:"ocean"`
	Routes     []TradeRouteResponse `json:"routes"`
	ScrapedAt  time.Time            `json:"scraped_at"`
	TotalRoutes int                 `json:"total_routes"`
}

type EconomySummaryResponse struct {
	Ocean              string    `json:"ocean"`
	ScrapedAt          time.Time `json:"scraped_at"`
	TotalCommodities   int       `json:"total_commodities"`
	TotalIslands       int       `json:"total_islands"`
	ColonizedIslands   int       `json:"colonized_islands"`
	
	MostTradedCommodities []CommodityTradeSummary `json:"most_traded_commodities"`
	
	MostActiveIslands []IslandActivitySummary `json:"most_active_islands"`
}

type CommodityTradeSummary struct {
	Commodity    CommodityBrief `json:"commodity"`
	TotalVolume  int            `json:"total_volume"`
	AvgPrice     float64        `json:"avg_price"`
	IslandsCount int            `json:"islands_count"`
}

type IslandActivitySummary struct {
	Island           IslandBrief `json:"island"`
	CommoditiesCount int         `json:"commodities_count"`
	TotalListings    int         `json:"total_listings"`
}
