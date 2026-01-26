package api

import (
	"time"

	"cutlass_analytics/internal/dto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
    }))

    r.GET("/api/health", func(c *gin.Context) {
        healthCheckHandler(c, db)
    })

    return r
}

func healthCheckHandler(c *gin.Context, db *gorm.DB) {
	// Get underlying sql.DB to test connectivity
	sqlDB, err := db.DB()
	if err != nil {
		response := dto.HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Services: map[string]string{
				"database": "unhealthy",
			},
		}
		c.JSON(503, response)
		return
	}

	// Ping database to verify connectivity
	if err := sqlDB.Ping(); err != nil {
		response := dto.HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Services: map[string]string{
				"database": "unhealthy",
			},
		}
		c.JSON(503, response)
		return
	}

	response := dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"database": "healthy",
		},
	}
	c.JSON(200, response)
}