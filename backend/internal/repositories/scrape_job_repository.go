package repositories

import (
	"cutlass_analytics/internal/dto"
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"

	"gorm.io/gorm"
)

type ScrapeJobRepository struct {
	db *gorm.DB
}

func NewScrapeJobRepository(db *gorm.DB) *ScrapeJobRepository {
	return &ScrapeJobRepository{db: db}
}

func (r *ScrapeJobRepository) FindByID(id uint) (*models.ScrapeJob, error) {
	var job models.ScrapeJob
	err := r.db.First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *ScrapeJobRepository) List(req dto.ScrapeJobListRequest) ([]models.ScrapeJob, int64, error) {
	query := r.db.Model(&models.ScrapeJob{})

	// Apply filters
	if req.Ocean != "" {
		query = query.Where("ocean = ?", req.Ocean)
	}
	if req.JobType != "" {
		query = query.Where("job_type = ?", req.JobType)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Count total before pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	query = applySorting(query, req.SortBy, req.SortOrder)

	// Apply pagination
	query = query.Offset(req.Offset()).Limit(req.Limit())

	var jobs []models.ScrapeJob
	if err := query.Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (r *ScrapeJobRepository) GetCurrentStatus(ocean types.Ocean) (*models.ScrapeJob, error) {
	var job models.ScrapeJob
	err := r.db.Where("ocean = ? AND status = ?", ocean, models.ScrapeJobStatusRunning).
		Order("started_at DESC").
		First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *ScrapeJobRepository) GetLastCompleted(ocean types.Ocean, jobType models.ScrapeJobType) (*models.ScrapeJob, error) {
	var job models.ScrapeJob
	query := r.db.Where("ocean = ? AND status = ?", ocean, models.ScrapeJobStatusCompleted)
	if jobType != "" {
		query = query.Where("job_type = ?", jobType)
	}
	err := query.Order("ended_at DESC").
		First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *ScrapeJobRepository) GetStats(ocean types.Ocean) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total jobs
	var totalJobs int64
	r.db.Model(&models.ScrapeJob{}).Where("ocean = ?", ocean).Count(&totalJobs)
	stats["total_jobs"] = totalJobs

	// Successful jobs
	var successfulJobs int64
	r.db.Model(&models.ScrapeJob{}).
		Where("ocean = ? AND status = ?", ocean, models.ScrapeJobStatusCompleted).
		Count(&successfulJobs)
	stats["successful_jobs"] = successfulJobs

	// Failed jobs
	var failedJobs int64
	r.db.Model(&models.ScrapeJob{}).
		Where("ocean = ? AND status = ?", ocean, models.ScrapeJobStatusFailed).
		Count(&failedJobs)
	stats["failed_jobs"] = failedJobs

	// Success rate
	if totalJobs > 0 {
		stats["success_rate"] = float64(successfulJobs) / float64(totalJobs) * 100
	} else {
		stats["success_rate"] = 0.0
	}

	return stats, nil
}
