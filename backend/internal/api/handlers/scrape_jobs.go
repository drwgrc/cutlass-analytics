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

func ListScrapeJobsHandler(c *gin.Context, db *gorm.DB) {
	var req dto.ScrapeJobListRequest
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

	repo := repositories.NewScrapeJobRepository(db)
	jobs, total, err := repo.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch scrape jobs",
				Details: err.Error(),
			},
		})
		return
	}

	responses := make([]dto.ScrapeJobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = toScrapeJobResponse(&job)
	}

	pagination := buildPagination(total, req.Page, req.PerPage)

	c.JSON(http.StatusOK, dto.ScrapeJobListResponse{
		Jobs:       responses,
		Pagination: pagination,
	})
}

func GetScrapeJobHandler(c *gin.Context, db *gorm.DB) {
	var param dto.ScrapeJobIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid scrape job ID",
			},
		})
		return
	}

	repo := repositories.NewScrapeJobRepository(db)
	job, err := repo.FindByID(param.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error: &dto.APIError{
					Code:    "NOT_FOUND",
					Message: "Scrape job not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch scrape job",
			},
		})
		return
	}

	response := dto.ScrapeJobDetailResponse{
		ScrapeJobResponse: toScrapeJobResponse(job),
	}

	c.JSON(http.StatusOK, response)
}

func GetScrapeStatusHandler(c *gin.Context, db *gorm.DB) {
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

	repo := repositories.NewScrapeJobRepository(db)
	currentJob, err := repo.GetCurrentStatus(types.Ocean(oceanParam.Ocean))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch scrape status",
			},
		})
		return
	}

	var currentJobResponse *dto.ScrapeJobResponse
	if currentJob != nil {
		response := toScrapeJobResponse(currentJob)
		currentJobResponse = &response
	}

	lastCompleted, err := repo.GetLastCompleted(types.Ocean(oceanParam.Ocean), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch last completed job",
			},
		})
		return
	}

	var lastCompletedResponse *dto.ScrapeJobResponse
	if lastCompleted != nil {
		response := toScrapeJobResponse(lastCompleted)
		lastCompletedResponse = &response
	}

	response := dto.ScrapeStatusResponse{
		IsRunning:     currentJob != nil,
		CurrentJob:    currentJobResponse,
		LastCompleted: lastCompletedResponse,
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions

func toScrapeJobResponse(job *models.ScrapeJob) dto.ScrapeJobResponse {
	duration := job.Duration()
	durationMs := duration.Milliseconds()

	var durationStr string
	if job.EndedAt != nil {
		durationStr = duration.String()
	}

	return dto.ScrapeJobResponse{
		ID:             job.ID,
		Ocean:          string(job.Ocean),
		JobType:        string(job.JobType),
		Status:         string(job.Status),
		StartedAt:      job.StartedAt,
		EndedAt:        job.EndedAt,
		Duration:       durationStr,
		DurationMs:     durationMs,
		ItemsProcessed: job.ItemsProcessed,
		ItemsFailed:    job.ItemsFailed,
		SuccessRate:    job.SuccessRate(),
		ErrorMessage:   job.ErrorMessage,
	}
}
