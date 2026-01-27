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

func ListFlagsHandler(c *gin.Context, db *gorm.DB) {
	var req dto.FlagListRequest
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

	repo := repositories.NewFlagRepository(db)
	flags, total, err := repo.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch flags",
				Details: err.Error(),
			},
		})
		return
	}

	responses := make([]dto.FlagResponse, len(flags))
	for i, flag := range flags {
		responses[i] = toFlagResponse(&flag)
	}

	pagination := buildPagination(total, req.Page, req.PerPage)

	c.JSON(http.StatusOK, dto.FlagListResponse{
		Flags:      responses,
		Pagination: pagination,
	})
}

func GetFlagHandler(c *gin.Context, db *gorm.DB) {
	var param dto.FlagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid flag ID",
			},
		})
		return
	}

	repo := repositories.NewFlagRepository(db)
	flag, err := repo.FindByID(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Flag not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch flag",
			},
		})
		return
	}

	response := toFlagResponse(flag)
	c.JSON(http.StatusOK, response)
}

func GetFlagByGameIDHandler(c *gin.Context, db *gorm.DB) {
	var param dto.FlagGameIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid game flag ID",
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

	repo := repositories.NewFlagRepository(db)
	flag, err := repo.FindByGameID(param.GameFlagID, types.Ocean(oceanParam.Ocean))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Flag not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch flag",
			},
		})
		return
	}

	response := toFlagResponse(flag)
	c.JSON(http.StatusOK, response)
}

func GetFlagCrewsHandler(c *gin.Context, db *gorm.DB) {
	var param dto.FlagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid flag ID",
			},
		})
		return
	}

	var req dto.FlagCrewsRequest
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
	req.SetDefaults()

	repo := repositories.NewFlagRepository(db)
	crews, err := repo.GetCrews(param.ID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch flag crews",
			},
		})
		return
	}

	// Get latest battle records for stats
	crewRepo := repositories.NewCrewRepository(db)
	crewStats := make([]dto.CrewPVPStatsResponse, 0, len(crews))
	for _, crew := range crews {
		battleRecord, err := crewRepo.GetLatestBattleRecord(crew.ID)
		if err != nil {
			continue // Skip crews without battle records
		}

		crewStats = append(crewStats, dto.CrewPVPStatsResponse{
			CrewID:         crew.ID,
			GameCrewID:     crew.GameCrewID,
			Name:           crew.Name,
			Ocean:          string(crew.Ocean),
			FlagID:         crew.FlagID,
			CrewRank:       string(battleRecord.CrewRank),
			TotalPVPWins:   battleRecord.TotalPVPWins,
			TotalPVPLosses: battleRecord.TotalPVPLosses,
			WinRate:        battleRecord.WinRate(),
			TotalBattles:   battleRecord.TotalBattles(),
			LastUpdated:    battleRecord.ScrapedAt,
		})
	}

	// Apply sorting if needed
	// For now, we'll return as-is. Could add sorting logic here if needed.

	// Get flag details
	flag, err := repo.FindByID(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch flag",
			},
		})
		return
	}

	response := dto.FlagDetailResponse{
		FlagResponse: toFlagResponse(flag),
		CrewCount:    len(crews),
		Crews:        crewStats,
	}

	c.JSON(http.StatusOK, response)
}

func GetFlagFameHandler(c *gin.Context, db *gorm.DB) {
	var param dto.FlagIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid flag ID",
			},
		})
		return
	}

	repo := repositories.NewFlagRepository(db)
	records, err := repo.GetFameHistory(param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch fame history",
			},
		})
		return
	}

	fameHistory := make([]dto.FlagFameResponse, len(records))
	for i, record := range records {
		fameHistory[i] = dto.FlagFameResponse{
			FlagID:    record.FlagID,
			ScrapedAt: record.ScrapedAt,
			FameLevel: string(record.FameLevel),
			FameRank:  record.FameRank,
		}
	}

	response := dto.FlagFameHistoryResponse{
		FlagID:  param.ID,
		History: fameHistory,
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions

func toFlagResponse(flag *models.Flag) dto.FlagResponse {
	return dto.FlagResponse{
		ID:          flag.ID,
		GameFlagID:  flag.GameFlagID,
		Name:        flag.Name,
		Ocean:       string(flag.Ocean),
		IsActive:    flag.IsActive,
		FirstSeenAt: flag.FirstSeenAt,
		LastSeenAt:  flag.LastSeenAt,
		URL:         flag.GetYowebURL(),
	}
}
