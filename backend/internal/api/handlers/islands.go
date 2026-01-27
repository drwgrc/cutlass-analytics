package handlers

import (
	"cutlass_analytics/internal/dto"
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/repositories"
	"cutlass_analytics/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListIslandsHandler(c *gin.Context, db *gorm.DB) {
	var req dto.IslandListRequest
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
	req.SetDefaults()

	repo := repositories.NewIslandRepository(db)
	islands, total, err := repo.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch islands",
				Details: err.Error(),
			},
		})
		return
	}

	responses := make([]dto.IslandResponse, len(islands))
	for i, island := range islands {
		responses[i] = toIslandResponse(&island)
	}

	pagination := buildPagination(total, req.Page, req.PerPage)

	c.JSON(http.StatusOK, dto.IslandListResponse{
		Islands:    responses,
		Pagination: pagination,
	})
}

func GetIslandHandler(c *gin.Context, db *gorm.DB) {
	var param dto.IslandIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid island ID",
			},
		})
		return
	}

	repo := repositories.NewIslandRepository(db)
	island, err := repo.FindByID(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Island not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch island",
			},
		})
		return
	}

	response := toIslandResponse(island)
	c.JSON(http.StatusOK, response)
}

func GetIslandByGameIDHandler(c *gin.Context, db *gorm.DB) {
	var param dto.IslandGameIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid game island ID",
			},
		})
		return
	}

	var oceanParam dto.OceanParam
	if err := c.ShouldBindQuery(&oceanParam); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Ocean parameter is required",
			},
		})
		return
	}

	repo := repositories.NewIslandRepository(db)
	island, err := repo.FindByGameID(param.GameIslandID, types.Ocean(oceanParam.Ocean))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Island not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch island",
			},
		})
		return
	}

	response := toIslandResponse(island)
	c.JSON(http.StatusOK, response)
}

func GetIslandPopulationHandler(c *gin.Context, db *gorm.DB) {
	var param dto.IslandIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid island ID",
			},
		})
		return
	}

	var req dto.IslandPopulationHistoryRequest
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

	repo := repositories.NewIslandRepository(db)
	populations, err := repo.GetPopulationHistory(param.ID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch population history",
			},
		})
		return
	}

	// Get island info
	island, err := repo.FindByID(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch island",
			},
		})
		return
	}

	dataPoints := make([]dto.IslandPopulationResponse, len(populations))
	for i, pop := range populations {
		dataPoints[i] = dto.IslandPopulationResponse{
			IslandID:  pop.IslandID,
			ScrapedAt: pop.ScrapedAt,
			Population: pop.Population,
		}
	}

	// Calculate stats
	minPop, maxPop, avgPop := calculatePopulationStats(populations)

	response := dto.IslandPopulationHistoryResponse{
		Island: dto.IslandBrief{
			ID:           island.ID,
			GameIslandID: island.GameIslandID,
			Name:         island.Name,
			Ocean:        string(island.Ocean),
			IsColonized:  island.IsColonized,
		},
		StartDate:  startDate,
		EndDate:    endDate,
		DataPoints: dataPoints,
		MinPopulation: minPop,
		MaxPopulation: maxPop,
		AvgPopulation: avgPop,
	}

	c.JSON(http.StatusOK, response)
}

func GetIslandGovernanceHandler(c *gin.Context, db *gorm.DB) {
	var param dto.IslandIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid island ID",
			},
		})
		return
	}

	repo := repositories.NewIslandRepository(db)
	history, err := repo.GetGovernanceHistory(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch governance history",
			},
		})
		return
	}

	// Get island info
	island, err := repo.FindByID(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch island",
			},
		})
		return
	}

	governanceHistory := make([]dto.GovernanceChangeResponse, len(history))
	for i, h := range history {
		var flagName string
		if h.Flag != nil {
			flagName = h.Flag.Name
		}
		governanceHistory[i] = dto.GovernanceChangeResponse{
			FlagID:       h.FlagID,
			FlagName:     flagName,
			GovernorName: h.GovernorName,
			StartedAt:    h.StartedAt,
			EndedAt:       h.EndedAt,
			ChangeType:   h.ChangeType,
			IsCurrent:    h.IsCurrent(),
		}
	}

	response := dto.IslandGovernanceHistoryResponse{
		Island: dto.IslandBrief{
			ID:           island.ID,
			GameIslandID: island.GameIslandID,
			Name:         island.Name,
			Ocean:        string(island.Ocean),
			IsColonized:  island.IsColonized,
		},
		History: governanceHistory,
	}

	c.JSON(http.StatusOK, response)
}

func GetIslandCommoditiesHandler(c *gin.Context, db *gorm.DB) {
	var param dto.IslandIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid island ID",
			},
		})
		return
	}

	repo := repositories.NewIslandRepository(db)
	commodities, err := repo.GetCommodities(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch island commodities",
			},
		})
		return
	}

	// Get island info
	island, err := repo.FindByID(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch island",
			},
		})
		return
	}

	commodityInfos := make([]dto.CommoditySpawnInfo, len(commodities))
	for i, ic := range commodities {
		commodityInfos[i] = dto.CommoditySpawnInfo{
			Commodity: dto.CommodityBrief{
				ID:          ic.Commodity.ID,
				Name:        ic.Commodity.Name,
				DisplayName: ic.Commodity.DisplayName,
				Category:    string(ic.Commodity.Category),
			},
			IsConfirmed: ic.IsConfirmed,
		}
	}

	response := dto.IslandCommoditiesResponse{
		Island: dto.IslandBrief{
			ID:           island.ID,
			GameIslandID: island.GameIslandID,
			Name:         island.Name,
			Ocean:        string(island.Ocean),
			IsColonized:  island.IsColonized,
		},
		Commodities: commodityInfos,
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions

func toIslandResponse(island *models.Island) dto.IslandResponse {
	response := dto.IslandResponse{
		ID:           island.ID,
		GameIslandID: island.GameIslandID,
		Name:         island.Name,
		Ocean:        string(island.Ocean),
		Size:         string(island.Size),
		IsColonized:  island.IsColonized,
		FirstSeenAt:  island.FirstSeenAt,
		LastSeenAt:   island.LastSeenAt,
		URL:          island.GetYowebURL(),
	}

	// Get latest population
	// Note: This could be optimized by joining with island_populations
	// For now, we'll use the population field on the island model
	response.Population = island.Population

	// Set archipelago if present
	if island.Archipelago != nil {
		response.Archipelago = &dto.ArchipelagoBrief{
			ID:    island.Archipelago.ID,
			Name:  island.Archipelago.Name,
			Color: island.Archipelago.Color,
		}
	}

	// Set governor if present
	if island.GovernorFlag != nil {
		response.Governor = &dto.IslandGovernor{
			FlagID:       island.GovernorFlagID,
			FlagName:     island.GovernorFlag.Name,
			GovernorName: island.GovernorName,
		}
	}

	return response
}

func buildPagination(total int64, page, perPage int) dto.Pagination {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	if totalPages == 0 {
		totalPages = 1
	}

	return dto.Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func calculatePopulationStats(populations []models.IslandPopulation) (min, max, avg int) {
	if len(populations) == 0 {
		return 0, 0, 0
	}

	min = populations[0].Population
	max = populations[0].Population
	sum := 0

	for _, pop := range populations {
		if pop.Population < min {
			min = pop.Population
		}
		if pop.Population > max {
			max = pop.Population
		}
		sum += pop.Population
	}

	avg = sum / len(populations)
	return min, max, avg
}
