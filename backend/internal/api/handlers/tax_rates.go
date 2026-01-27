package handlers

import (
	"cutlass_analytics/internal/dto"
	"cutlass_analytics/internal/repositories"
	"cutlass_analytics/internal/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTaxRatesHandler(c *gin.Context, db *gorm.DB) {
	var req dto.TaxRatesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request parameters",
				Details: err.Error(),
			},
		})
		return
	}

	repo := repositories.NewTaxRateRepository(db)
	rates, err := repo.GetCurrentRates(types.Ocean(req.Ocean), req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch tax rates",
				Details: err.Error(),
			},
		})
		return
	}

	responses := make([]dto.CommodityTaxRateResponse, len(rates))
	var scrapedAt time.Time
	for i, rate := range rates {
		responses[i] = dto.CommodityTaxRateResponse{
			CommodityID:   rate.CommodityID,
			CommodityName: rate.Commodity.Name,
			Ocean:         string(rate.Ocean),
			TaxValue:      rate.TaxValue,
			ScrapedAt:     rate.ScrapedAt,
		}
		if rate.ScrapedAt.After(scrapedAt) {
			scrapedAt = rate.ScrapedAt
		}
	}

	response := dto.OceanTaxRatesResponse{
		Ocean:     req.Ocean,
		ScrapedAt: scrapedAt,
		TaxRates:  responses,
	}

	// Group by category if requested
	if req.GroupByCategory {
		byCategory := make(map[string][]dto.CommodityTaxRateResponse)
		for _, rate := range responses {
			// Get category from commodity
			for _, r := range rates {
				if r.CommodityID == rate.CommodityID {
					category := string(r.Commodity.Category)
					byCategory[category] = append(byCategory[category], rate)
					break
				}
			}
		}
		response.ByCategory = byCategory
	}

	c.JSON(http.StatusOK, response)
}

func GetTaxRateHistoryHandler(c *gin.Context, db *gorm.DB) {
	type CommodityIDPathParam struct {
		CommodityID uint `uri:"commodity_id" binding:"required,min=1"`
	}
	var param CommodityIDPathParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid commodity ID",
			},
		})
		return
	}

	var req dto.TaxRateHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request parameters",
			},
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	startDate, _ := req.ParsedStartDate()
	endDate, _ := req.ParsedEndDate()

	repo := repositories.NewTaxRateRepository(db)
	rates, err := repo.GetHistory(param.CommodityID, types.Ocean(req.Ocean), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch tax rate history",
			},
		})
		return
	}

	if len(rates) == 0 {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "NOT_FOUND",
				Message: "No tax rate history found",
			},
		})
		return
	}

	dataPoints := make([]dto.TaxRateHistoryPoint, len(rates))
	for i, rate := range rates {
		dataPoints[i] = dto.TaxRateHistoryPoint{
			Date:     rate.ScrapedAt,
			TaxValue: rate.TaxValue,
		}
	}

	response := dto.TaxRateHistoryResponse{
		Commodity: dto.CommodityBrief{
			ID:          rates[0].Commodity.ID,
			Name:        rates[0].Commodity.Name,
			DisplayName: rates[0].Commodity.DisplayName,
			Category:    string(rates[0].Commodity.Category),
		},
		Ocean:      req.Ocean,
		StartDate:  startDate,
		EndDate:    endDate,
		DataPoints: dataPoints,
	}

	c.JSON(http.StatusOK, response)
}

func CompareTaxRatesHandler(c *gin.Context, db *gorm.DB) {
	var req dto.TaxRateComparisonRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Commodity ID is required",
			},
		})
		return
	}

	repo := repositories.NewTaxRateRepository(db)
	rates, err := repo.CompareAcrossOceans(req.CommodityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch tax rate comparison",
			},
		})
		return
	}

	if len(rates) == 0 {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "NOT_FOUND",
				Message: "No tax rates found for this commodity",
			},
		})
		return
	}

	oceanRates := make([]dto.OceanTaxRateEntry, len(rates))
	var updatedAt time.Time
	for i, rate := range rates {
		oceanRates[i] = dto.OceanTaxRateEntry{
			Ocean:    string(rate.Ocean),
			TaxValue: rate.TaxValue,
		}
		if rate.ScrapedAt.After(updatedAt) {
			updatedAt = rate.ScrapedAt
		}
	}

	response := dto.TaxRateComparisonResponse{
		Commodity: dto.CommodityBrief{
			ID:          rates[0].Commodity.ID,
			Name:        rates[0].Commodity.Name,
			DisplayName: rates[0].Commodity.DisplayName,
			Category:    string(rates[0].Commodity.Category),
		},
		OceanRates: oceanRates,
		UpdatedAt:  updatedAt,
	}

	c.JSON(http.StatusOK, response)
}
