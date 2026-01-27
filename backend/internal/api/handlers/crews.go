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

func ListCrewsHandler(c *gin.Context, db *gorm.DB) {
	var req dto.CrewListRequest
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

	repo := repositories.NewCrewRepository(db)
	crews, total, err := repo.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch crews",
				Details: err.Error(),
			},
		})
		return
	}

	responses := make([]dto.CrewResponse, len(crews))
	for i, crew := range crews {
		responses[i] = toCrewResponse(&crew)
	}

	pagination := buildPagination(total, req.Page, req.PerPage)

	c.JSON(http.StatusOK, dto.CrewListResponse{
		Crews:      responses,
		Pagination: pagination,
	})
}

func GetCrewHandler(c *gin.Context, db *gorm.DB) {
	var param dto.CrewIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid crew ID",
			},
		})
		return
	}

	repo := repositories.NewCrewRepository(db)
	crew, err := repo.FindByID(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Crew not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch crew",
			},
		})
		return
	}

	response := toCrewResponse(crew)
	c.JSON(http.StatusOK, response)
}

func GetCrewByGameIDHandler(c *gin.Context, db *gorm.DB) {
	var param dto.CrewGameIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid game crew ID",
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

	repo := repositories.NewCrewRepository(db)
	crew, err := repo.FindByGameID(param.GameCrewID, types.Ocean(oceanParam.Ocean))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Crew not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch crew",
			},
		})
		return
	}

	response := toCrewResponse(crew)
	c.JSON(http.StatusOK, response)
}

func GetCrewBattlesHandler(c *gin.Context, db *gorm.DB) {
	var param dto.CrewIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid crew ID",
			},
		})
		return
	}

	var req dto.CrewBattleRecordsRequest
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

	repo := repositories.NewCrewRepository(db)
	records, err := repo.GetBattleRecords(param.ID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch battle records",
			},
		})
		return
	}

	battleRecords := make([]dto.CrewBattleRecordResponse, len(records))
	for i, record := range records {
		battleRecords[i] = dto.CrewBattleRecordResponse{
			ID:             record.ID,
			CrewID:         record.CrewID,
			ScrapedAt:      record.ScrapedAt,
			CrewRank:       string(record.CrewRank),
			TotalPVPWins:   record.TotalPVPWins,
			TotalPVPLosses: record.TotalPVPLosses,
			DailyPVPWins:   record.DailyPVPWins,
			DailyPVPLosses: record.DailyPVPLosses,
			WinRate:        record.WinRate(),
			TotalBattles:   record.TotalBattles(),
		}
	}

	response := dto.CrewBattleRecordListResponse{
		CrewID:  param.ID,
		Records: battleRecords,
	}

	c.JSON(http.StatusOK, response)
}

func GetCrewFameHandler(c *gin.Context, db *gorm.DB) {
	var param dto.CrewIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid crew ID",
			},
		})
		return
	}

	repo := repositories.NewCrewRepository(db)
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

	fameHistory := make([]dto.CrewFameResponse, len(records))
	for i, record := range records {
		fameHistory[i] = dto.CrewFameResponse{
			CrewID:    record.CrewID,
			ScrapedAt: record.ScrapedAt,
			FameLevel: string(record.FameLevel),
			FameRank:  record.FameRank,
		}
	}

	response := dto.CrewFameHistoryResponse{
		CrewID:  param.ID,
		History: fameHistory,
	}

	c.JSON(http.StatusOK, response)
}

func GetCrewStatsHandler(c *gin.Context, db *gorm.DB) {
	var param dto.CrewIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid crew ID",
			},
		})
		return
	}

	repo := repositories.NewCrewRepository(db)
	crew, err := repo.FindByID(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Crew not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch crew",
			},
		})
		return
	}

	battleRecord, err := repo.GetCurrentStats(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "No battle records found for this crew",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch battle stats",
			},
		})
		return
	}

	var flagID *uint
	var flagName string
	if crew.Flag != nil {
		flagID = crew.FlagID
		flagName = crew.Flag.Name
	}

	response := dto.CrewPVPStatsResponse{
		CrewID:         crew.ID,
		GameCrewID:     crew.GameCrewID,
		Name:           crew.Name,
		Ocean:          string(crew.Ocean),
		FlagID:         flagID,
		FlagName:       flagName,
		CrewRank:       string(battleRecord.CrewRank),
		TotalPVPWins:   battleRecord.TotalPVPWins,
		TotalPVPLosses: battleRecord.TotalPVPLosses,
		WinRate:        battleRecord.WinRate(),
		TotalBattles:   battleRecord.TotalBattles(),
		LastUpdated:    battleRecord.ScrapedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions

func toCrewResponse(crew *models.Crew) dto.CrewResponse {
	response := dto.CrewResponse{
		ID:          crew.ID,
		GameCrewID:  crew.GameCrewID,
		Name:        crew.Name,
		Ocean:       string(crew.Ocean),
		IsActive:    crew.IsActive,
		FirstSeenAt: crew.FirstSeenAt,
		LastSeenAt:  crew.LastSeenAt,
		URLs: dto.CrewURLs{
			Info:       crew.GetYowebURL(),
			BattleInfo: crew.GetBattleInfoURL(),
		},
	}

	if crew.Flag != nil {
		response.Flag = &dto.FlagBrief{
			ID:         crew.Flag.ID,
			GameFlagID: crew.Flag.GameFlagID,
			Name:       crew.Flag.Name,
			Ocean:      string(crew.Flag.Ocean),
		}
	}

	return response
}
