package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// type Handlers struct {
//     rankingRepo      *repository.RankingRepository
//     powerRankingRepo *repository.PowerRankingRepository
//     marketRepo       *repository.MarketRepository
// }

func NewRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
    }))

    // Initialize handlers with repositories
    // h := &Handlers{
    //     rankingRepo:      repository.NewRankingRepository(db),
    //     powerRankingRepo: repository.NewPowerRankingRepository(db),
    //     marketRepo:       repository.NewMarketRepository(db),
    // }

    // Health check
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API routes
    // api := r.Group("/api")
    {
        // Rankings
        // api.GET("/rankings/:entityType", h.GetRankings)
        // api.GET("/rankings/:entityType/:id", h.GetRankingDetail)
        // api.GET("/rankings/:entityType/:id/history", h.GetRankingHistory)

        // Power Rankings
        // api.GET("/power-rankings/:entityType", h.GetPowerRankings)

        // Market
        // api.GET("/commodities", h.GetCommodities)
        // api.GET("/islands", h.GetIslands)
        // api.GET("/market/orders", h.GetMarketOrders)
    }

    return r
}