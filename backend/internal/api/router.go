package api

import (
	"net/http"
	"time"

	"cutlass_analytics/docs"
	"cutlass_analytics/internal/api/handlers"
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

    // API Documentation
    r.GET("/api/docs", serveSwaggerUI)
    r.GET("/api/docs/openapi.yaml", serveOpenAPISpec)

    // API routes
    api := r.Group("/api")
    {
        // Islands
        api.GET("/islands", func(c *gin.Context) { handlers.ListIslandsHandler(c, db) })
        api.GET("/islands/:id", func(c *gin.Context) { handlers.GetIslandHandler(c, db) })
        api.GET("/islands/game/:game_island_id", func(c *gin.Context) { handlers.GetIslandByGameIDHandler(c, db) })
        api.GET("/islands/:id/population", func(c *gin.Context) { handlers.GetIslandPopulationHandler(c, db) })
        api.GET("/islands/:id/governance", func(c *gin.Context) { handlers.GetIslandGovernanceHandler(c, db) })
        api.GET("/islands/:id/commodities", func(c *gin.Context) { handlers.GetIslandCommoditiesHandler(c, db) })

        // Crews
        api.GET("/crews", func(c *gin.Context) { handlers.ListCrewsHandler(c, db) })
        api.GET("/crews/:id", func(c *gin.Context) { handlers.GetCrewHandler(c, db) })
        api.GET("/crews/game/:game_crew_id", func(c *gin.Context) { handlers.GetCrewByGameIDHandler(c, db) })
        api.GET("/crews/:id/battles", func(c *gin.Context) { handlers.GetCrewBattlesHandler(c, db) })
        api.GET("/crews/:id/fame", func(c *gin.Context) { handlers.GetCrewFameHandler(c, db) })
        api.GET("/crews/:id/stats", func(c *gin.Context) { handlers.GetCrewStatsHandler(c, db) })

        // Flags
        api.GET("/flags", func(c *gin.Context) { handlers.ListFlagsHandler(c, db) })
        api.GET("/flags/:id", func(c *gin.Context) { handlers.GetFlagHandler(c, db) })
        api.GET("/flags/game/:game_flag_id", func(c *gin.Context) { handlers.GetFlagByGameIDHandler(c, db) })
        api.GET("/flags/:id/crews", func(c *gin.Context) { handlers.GetFlagCrewsHandler(c, db) })
        api.GET("/flags/:id/fame", func(c *gin.Context) { handlers.GetFlagFameHandler(c, db) })

        // Scrape Jobs
        api.GET("/scrape-jobs", func(c *gin.Context) { handlers.ListScrapeJobsHandler(c, db) })
        api.GET("/scrape-jobs/:id", func(c *gin.Context) { handlers.GetScrapeJobHandler(c, db) })
        api.GET("/scrape-jobs/status", func(c *gin.Context) { handlers.GetScrapeStatusHandler(c, db) })

        // Tax Rates
        api.GET("/tax-rates", func(c *gin.Context) { handlers.GetTaxRatesHandler(c, db) })
        api.GET("/tax-rates/:commodity_id/history", func(c *gin.Context) { handlers.GetTaxRateHistoryHandler(c, db) })
        api.GET("/tax-rates/compare", func(c *gin.Context) { handlers.CompareTaxRatesHandler(c, db) })
    }

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

func serveOpenAPISpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/x-yaml", docs.OpenAPISpec)
}

func serveSwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cutlass Analytics API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "/api/docs/openapi.yaml",
                dom_id: '#swagger-ui',
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
                layout: "BaseLayout",
                deepLinking: true,
                showExtensions: true,
                showCommonExtensions: true
            });
        };
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}