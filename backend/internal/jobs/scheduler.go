package jobs

import (
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/scraper"
	"cutlass_analytics/internal/types"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// Scheduler manages cron jobs for scraping operations
type Scheduler struct {
	cron    *cron.Cron
	db      *gorm.DB
	stop    chan struct{}
	wg      sync.WaitGroup
	running bool
	mu      sync.Mutex
}

// NewScheduler creates a new scheduler instance
func NewScheduler(db *gorm.DB) *Scheduler {
	// Create cron with PST timezone
	pstLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Printf("Warning: Failed to load PST timezone, using UTC: %v", err)
		pstLocation = time.UTC
	}

	c := cron.New(cron.WithLocation(pstLocation))

	return &Scheduler{
		cron: c,
		db:   db,
		stop: make(chan struct{}),
	}
}

// Start initializes and starts the scheduler
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	// Schedule daily scraper at 3:30 AM PST
	// Cron format: minute hour day month day-of-week
	_, err := s.cron.AddFunc("30 3 * * *", s.runDailyScrapers)
	if err != nil {
		return err
	}

	s.cron.Start()
	s.running = true
	log.Println("Scheduler started - Daily scraper scheduled for 3:30 AM PST")

	return nil
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	ctx := s.cron.Stop()
	select {
	case <-ctx.Done():
		log.Println("Scheduler stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Println("Warning: Scheduler stop timeout")
	}

	close(s.stop)
	s.wg.Wait()
	s.running = false
}

// RunOnce runs all scrapers immediately in a goroutine
// This is useful for running jobs on server startup
func (s *Scheduler) RunOnce() {
	go s.runDailyScrapers()
	log.Println("Triggered immediate scraper job run on server startup")
}

// runDailyScrapers runs all scrapers for all oceans
func (s *Scheduler) runDailyScrapers() {
	log.Println("Starting daily scraper job for all oceans")

	oceans := []types.Ocean{
		types.OceanEmerald,
		types.OceanMeridian,
		types.OceanCerulean,
		types.OceanObsidian,
	}

	var wg sync.WaitGroup
	for _, ocean := range oceans {
		wg.Add(1)
		go func(o types.Ocean) {
			defer wg.Done()
			if err := s.runScraperForOcean(o); err != nil {
				log.Printf("Error running scraper for ocean %s: %v", o, err)
			}
		}(ocean)
	}

	wg.Wait()
	log.Println("Daily scraper job completed for all oceans")
}

// runScraperForOcean runs the full scraper for a specific ocean
func (s *Scheduler) runScraperForOcean(ocean types.Ocean) error {
	log.Printf("Starting scraper for ocean: %s", ocean)

	scraperInstance, err := scraper.NewScraper(s.db, ocean, models.ScrapeJobTypeDailyFull)
	if err != nil {
		return err
	}

	if err := scraperInstance.Run(); err != nil {
		log.Printf("Scraper failed for ocean %s: %v", ocean, err)
		return err
	}

	log.Printf("Scraper completed successfully for ocean: %s", ocean)
	return nil
}
