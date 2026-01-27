package scraper

import (
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"gorm.io/gorm"
)

// Scraper handles web scraping operations
type Scraper struct {
	db        *gorm.DB
	collector *colly.Collector
	job       *models.ScrapeJob
	ocean     types.Ocean
}

// fetchHTML fetches HTML content from a URL
func (s *Scraper) fetchHTML(url string) (string, error) {
	var htmlContent string
	var fetchErr error

	// Create a temporary collector for this request to avoid handler conflicts
	tempCollector := s.collector.Clone()
	tempCollector.OnResponse(func(r *colly.Response) {
		htmlContent = string(r.Body)
	})
	tempCollector.OnError(func(r *colly.Response, err error) {
		fetchErr = err
	})

	if err := tempCollector.Visit(url); err != nil {
		return "", err
	}

	if fetchErr != nil {
		return "", fetchErr
	}

	return htmlContent, nil
}

// NewScraper creates a new scraper instance with rate limiting
func NewScraper(db *gorm.DB, ocean types.Ocean, jobType models.ScrapeJobType) (*Scraper, error) {
	// Create scrape job
	job := &models.ScrapeJob{
		Ocean:   ocean,
		JobType: jobType,
		Status:  models.ScrapeJobStatusRunning,
	}

	if err := db.Create(job).Error; err != nil {
		return nil, fmt.Errorf("failed to create scrape job: %w", err)
	}

	// Create collector with rate limiting
	collector := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Set rate limiting: 1 request per second
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*.puzzlepirates.com",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	// Set user agent
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

	// Handle errors
	collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with error: %v\n", r.Request.URL, err)
	})

	return &Scraper{
		db:        db,
		collector: collector,
		job:       job,
		ocean:     ocean,
	}, nil
}

// Run executes the scraper based on job type
func (s *Scraper) Run() error {
	defer func() {
		// Reload job to get latest counters
		s.db.First(s.job, s.job.ID)
		if s.job.Status == models.ScrapeJobStatusRunning {
			if err := s.job.MarkCompleted(s.db); err != nil {
				log.Printf("Failed to mark job as completed: %v", err)
			}
		}
	}()

	switch s.job.JobType {
	case models.ScrapeJobTypeDailyFull:
		// Run all scrapers
		if err := s.ScrapeIslands(); err != nil {
			log.Printf("Error scraping islands: %v", err)
			s.job.IncrementFailed()
		}
		if err := s.ScrapeTaxRates(); err != nil {
			log.Printf("Error scraping tax rates: %v", err)
			s.job.IncrementFailed()
		}
		if err := s.ScrapeCrews(); err != nil {
			log.Printf("Error scraping crews: %v", err)
			s.job.IncrementFailed()
		}
		if err := s.ScrapeFlags(); err != nil {
			log.Printf("Error scraping flags: %v", err)
			s.job.IncrementFailed()
		}
	case models.ScrapeJobTypeCrewInfo:
		if err := s.ScrapeCrews(); err != nil {
			return s.job.MarkFailed(s.db, err)
		}
	case models.ScrapeJobTypeCrewFame:
		if err := s.ScrapeCrewFame(); err != nil {
			return s.job.MarkFailed(s.db, err)
		}
	case models.ScrapeJobTypeFlagFame:
		if err := s.ScrapeFlagFame(); err != nil {
			return s.job.MarkFailed(s.db, err)
		}
	case models.ScrapeJobTypeBattleInfo:
		if err := s.ScrapeBattleInfo(); err != nil {
			return s.job.MarkFailed(s.db, err)
		}
	}

	return nil
}

// ScrapeIslands scrapes all island data by looping through island IDs 0-120
func (s *Scraper) ScrapeIslands() error {
	scrapedAt := time.Now()
	processedCount := 0

	// Loop through island IDs from 0 to 120
	for islandID := uint64(0); islandID <= 120; islandID++ {
		url := GetIslandInfoURL(s.ocean, islandID)
		htmlContent, err := s.fetchHTML(url)
		if err != nil {
			log.Printf("Failed to fetch island %d: %v", islandID, err)
			s.job.IncrementFailed()
			continue
		}

		// Check if island is uncolonized
		if strings.Contains(htmlContent, "Shiver me timbers: The island is uncolonized.") {
			// Skip uncolonized islands
			continue
		}

		// Parse island info
		islandData, err := ParseIslandInfo(htmlContent, islandID, s.ocean)
		if err != nil {
			log.Printf("Failed to parse island %d: %v", islandID, err)
			s.job.IncrementFailed()
			continue
		}

		// Process and save island
		if err := s.processIsland(*islandData, scrapedAt); err != nil {
			log.Printf("Error processing island %d: %v", islandID, err)
			s.job.IncrementFailed()
			continue
		}

		processedCount++
		s.job.IncrementProcessed()
		s.db.Save(s.job)
	}

	log.Printf("Successfully processed %d islands for ocean %s", processedCount, s.ocean)
	return nil
}

// processIsland processes a single island and saves all related data
func (s *Scraper) processIsland(data IslandData, scrapedAt time.Time) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Find or create archipelago
		var archipelago models.Archipelago
		if data.Archipelago != "" {
			err := tx.Where("ocean = ? AND name = ?", s.ocean, data.Archipelago).
				FirstOrCreate(&archipelago, models.Archipelago{
					Ocean: s.ocean,
					Name:  data.Archipelago,
				}).Error
			if err != nil {
				return fmt.Errorf("failed to get/create archipelago: %w", err)
			}
		}

		// Find or create island
		var island models.Island
		err := tx.Where("game_island_id = ? AND ocean = ?", data.GameIslandID, s.ocean).
			FirstOrCreate(&island, models.Island{
				GameIslandID: data.GameIslandID,
				Ocean:        s.ocean,
				Name:         data.Name,
				Size:         data.Size,
				IsColonized:  data.IsColonized,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to get/create island: %w", err)
		}

		// Update island fields
		island.Name = data.Name
		island.Size = data.Size
		island.IsColonized = data.IsColonized
		if archipelago.ID > 0 {
			island.ArchipelagoID = &archipelago.ID
		}
		island.LastSeenAt = scrapedAt

		// Handle governor flag
		if data.GovernorFlag != "" {
			var flagID uint64
			if _, err := fmt.Sscanf(data.GovernorFlag, "%d", &flagID); err == nil {
				var flag models.Flag
				if err := tx.Where("game_flag_id = ? AND ocean = ?", flagID, s.ocean).First(&flag).Error; err == nil {
					island.GovernorFlagID = &flag.ID
				}
			}
		}
		if data.GovernorName != "" {
			island.GovernorName = data.GovernorName
		}

		if err := tx.Save(&island).Error; err != nil {
			return fmt.Errorf("failed to save island: %w", err)
		}

		// Create population record
		if data.Population > 0 {
			pop := models.IslandPopulation{
				IslandID:  island.ID,
				ScrapedAt: scrapedAt,
				Population: data.Population,
			}
			// Use FirstOrCreate to avoid duplicates
			if err := tx.Where("island_id = ? AND scraped_at = ?", island.ID, scrapedAt).
				FirstOrCreate(&pop).Error; err != nil {
				return fmt.Errorf("failed to create population record: %w", err)
			}
		}

		// Update governance history if governor changed
		var lastGov models.IslandGovernanceHistory
		err = tx.Where("island_id = ? AND ended_at IS NULL", island.ID).
			Order("started_at DESC").First(&lastGov).Error

		governorChanged := false
		if err == gorm.ErrRecordNotFound {
			governorChanged = true
		} else if err == nil {
			currentFlagID := island.GovernorFlagID
			if (lastGov.FlagID == nil && currentFlagID != nil) ||
				(lastGov.FlagID != nil && currentFlagID == nil) ||
				(lastGov.FlagID != nil && currentFlagID != nil && *lastGov.FlagID != *currentFlagID) ||
				lastGov.GovernorName != island.GovernorName {
				governorChanged = true
			}
		}

		if governorChanged {
			// End previous governance
			if err == nil {
				now := time.Now()
				lastGov.EndedAt = &now
				tx.Save(&lastGov)
			}

			// Create new governance record
			gov := models.IslandGovernanceHistory{
				IslandID:     island.ID,
				FlagID:        island.GovernorFlagID,
				GovernorName:  island.GovernorName,
				StartedAt:     scrapedAt,
				ChangeType:    "scrape",
			}
			if err := tx.Create(&gov).Error; err != nil {
				return fmt.Errorf("failed to create governance history: %w", err)
			}
		}

		// Process commodities
		for _, commName := range data.Commodities {
			var commodity models.Commodity
			if err := tx.Where("name = ?", commName).First(&commodity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Create commodity if it doesn't exist
					commodity = models.Commodity{
						Name:        commName,
						DisplayName: commName,
						Category:    types.CommodityCategoryBasic, // Default category
					}
					if err := tx.Create(&commodity).Error; err != nil {
						return fmt.Errorf("failed to create commodity: %w", err)
					}
				} else {
					return fmt.Errorf("failed to find commodity: %w", err)
				}
			}

			var islandComm models.IslandCommodity
			if err := tx.Where("island_id = ? AND commodity_id = ?", island.ID, commodity.ID).
				FirstOrCreate(&islandComm, models.IslandCommodity{
					IslandID:    island.ID,
					CommodityID: commodity.ID,
					IsConfirmed: true,
				}).Error; err != nil {
				return fmt.Errorf("failed to create island commodity: %w", err)
			}
		}

		return nil
	})
}

// ScrapeTaxRates scrapes tax rates for all commodities
func (s *Scraper) ScrapeTaxRates() error {
	url := GetTaxRatesURL(s.ocean)
	htmlContent, err := s.fetchHTML(url)
	if err != nil {
		return fmt.Errorf("failed to fetch tax rates: %w", err)
	}

	// Parse tax rates
	rates, err := ParseTaxRates(htmlContent, s.ocean)
	if err != nil {
		return fmt.Errorf("failed to parse tax rates: %w", err)
	}

	if len(rates) == 0 {
		log.Printf("WARNING: Parsed 0 tax rates for ocean %s from URL: %s", s.ocean, url)
		return fmt.Errorf("no tax rates found in parsed HTML for ocean %s", s.ocean)
	}

	log.Printf("Successfully parsed %d tax rates for ocean %s", len(rates), s.ocean)

	scrapedAt := time.Now()

	for _, rateData := range rates {
		if err := s.processTaxRate(rateData, scrapedAt); err != nil {
			log.Printf("Error processing tax rate for %s: %v", rateData.CommodityName, err)
			s.job.IncrementFailed()
			continue
		}
		s.job.IncrementProcessed()
		s.db.Save(s.job)
	}

	return nil
}

// processTaxRate processes a single tax rate
func (s *Scraper) processTaxRate(data TaxRateData, scrapedAt time.Time) error {
	// Find or create commodity
	var commodity models.Commodity
	if err := s.db.Where("name = ?", data.CommodityName).First(&commodity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create commodity if it doesn't exist
			commodity = models.Commodity{
				Name:        data.CommodityName,
				DisplayName: data.CommodityName,
				Category:    types.CommodityCategoryBasic, // Default category
			}
			if err := s.db.Create(&commodity).Error; err != nil {
				return fmt.Errorf("failed to create commodity: %w", err)
			}
		} else {
			return fmt.Errorf("failed to find commodity: %w", err)
		}
	}

	// Create tax rate record
	taxRate := models.CommodityTaxRate{
		CommodityID: commodity.ID,
		Ocean:       s.ocean,
		ScrapedAt:   scrapedAt,
		TaxValue:    data.TaxValue,
	}

	// Use FirstOrCreate to avoid duplicates
	if err := s.db.Where("commodity_id = ? AND ocean = ? AND scraped_at = ?",
		commodity.ID, s.ocean, scrapedAt).FirstOrCreate(&taxRate).Error; err != nil {
		return fmt.Errorf("failed to create tax rate: %w", err)
	}

	return nil
}

// ScrapeCrews scrapes all crew data including fame and battle info
func (s *Scraper) ScrapeCrews() error {
	// First, get crew list from fame list
	url := GetCrewFameListURL(s.ocean)
	htmlContent, err := s.fetchHTML(url)
	if err != nil {
		return fmt.Errorf("failed to fetch crew fame list: %w", err)
	}

	// Parse crew fame list
	crews, err := ParseCrewFameList(htmlContent, s.ocean)
	if err != nil {
		return fmt.Errorf("failed to parse crew fame list: %w", err)
	}

	if len(crews) == 0 {
		log.Printf("WARNING: Parsed 0 crews for ocean %s from URL: %s", s.ocean, url)
		return fmt.Errorf("no crews found in parsed HTML for ocean %s", s.ocean)
	}

	log.Printf("Successfully parsed %d crews for ocean %s", len(crews), s.ocean)

	scrapedAt := time.Now()

	for _, crewData := range crews {
		if err := s.processCrew(crewData, scrapedAt); err != nil {
			log.Printf("Error processing crew %d: %v", crewData.CrewID, err)
			s.job.IncrementFailed()
			continue
		}
		s.job.IncrementProcessed()
		s.db.Save(s.job)
	}

	return nil
}

// processCrew processes a single crew and all related data
func (s *Scraper) processCrew(fameData CrewFameData, scrapedAt time.Time) error {
	// Fetch crew info page
	url := GetCrewInfoURL(s.ocean, fameData.CrewID)
	crewInfoHTML, err := s.fetchHTML(url)
	if err != nil {
		return fmt.Errorf("failed to fetch crew info: %w", err)
	}

	crewData, err := ParseCrewInfo(crewInfoHTML, fameData.CrewID, s.ocean)
	if err != nil {
		return fmt.Errorf("failed to parse crew info: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Find or create crew
		var crew models.Crew
		err := tx.Where("game_crew_id = ? AND ocean = ?", fameData.CrewID, s.ocean).
			FirstOrCreate(&crew, models.Crew{
				GameCrewID: fameData.CrewID,
				Ocean:      s.ocean,
				Name:       crewData.Name,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to get/create crew: %w", err)
		}

		// Update crew fields
		crew.Name = crewData.Name
		crew.LastSeenAt = scrapedAt

		// Handle flag
		if crewData.FlagID != nil {
			var flag models.Flag
			if err := tx.Where("game_flag_id = ? AND ocean = ?", *crewData.FlagID, s.ocean).
				First(&flag).Error; err == nil {
				oldFlagID := crew.FlagID
				crew.FlagID = &flag.ID

				// Update flag history if flag changed
				if oldFlagID == nil || *oldFlagID != flag.ID {
					// End previous flag history
					var lastFlagHist models.CrewFlagHistory
					if err := tx.Where("crew_id = ? AND left_at IS NULL", crew.ID).
						First(&lastFlagHist).Error; err == nil {
						now := time.Now()
						lastFlagHist.LeftAt = &now
						tx.Save(&lastFlagHist)
					}

					// Create new flag history
					flagHist := models.CrewFlagHistory{
						CrewID:  crew.ID,
						FlagID:  &flag.ID,
						JoinedAt: scrapedAt,
					}
					if err := tx.Create(&flagHist).Error; err != nil {
						return fmt.Errorf("failed to create flag history: %w", err)
					}
				}
			}
		} else {
			// Crew is independent
			oldFlagID := crew.FlagID
			crew.FlagID = nil

			if oldFlagID != nil {
				// End previous flag history
				var lastFlagHist models.CrewFlagHistory
				if err := tx.Where("crew_id = ? AND left_at IS NULL", crew.ID).
					First(&lastFlagHist).Error; err == nil {
					now := time.Now()
					lastFlagHist.LeftAt = &now
					tx.Save(&lastFlagHist)
				}
			}
		}

		if err := tx.Save(&crew).Error; err != nil {
			return fmt.Errorf("failed to save crew: %w", err)
		}

		// Create fame record
		fameRecord := models.CrewFameRecord{
			CrewID:    crew.ID,
			ScrapedAt: scrapedAt,
			FameLevel: fameData.FameLevel,
			FameRank:  fameData.Rank,
		}
		if err := tx.Where("crew_id = ? AND scraped_at = ?", crew.ID, scrapedAt).
			FirstOrCreate(&fameRecord).Error; err != nil {
			return fmt.Errorf("failed to create fame record: %w", err)
		}

		// Fetch and process battle info
		battleURL := GetCrewBattleInfoURL(s.ocean, fameData.CrewID)
		battleHTML, err := s.fetchHTML(battleURL)
		if err == nil {
			battleData, err := ParseCrewBattleInfo(battleHTML, fameData.CrewID)
			if err == nil {
				// Get previous battle record for delta calculation
				var prevRecord models.CrewBattleRecord
				err := tx.Where("crew_id = ?", crew.ID).
					Order("scraped_at DESC").First(&prevRecord).Error

				battleRecord := models.CrewBattleRecord{
					CrewID:         crew.ID,
					ScrapedAt:      scrapedAt,
					CrewRank:       battleData.CrewRank,
					TotalPVPWins:   battleData.TotalPVPWins,
					TotalPVPLosses: battleData.TotalPVPLosses,
				}

				if err == nil {
					battleRecord.CalculateDeltas(&prevRecord)
				} else {
					battleRecord.CalculateDeltas(nil)
				}

				if err := tx.Where("crew_id = ? AND scraped_at = ?", crew.ID, scrapedAt).
					FirstOrCreate(&battleRecord).Error; err != nil {
					return fmt.Errorf("failed to create battle record: %w", err)
				}
			}
		}

		return nil
	})
}

// ScrapeCrewFame scrapes only crew fame data
func (s *Scraper) ScrapeCrewFame() error {
	return s.ScrapeCrews() // Reuse ScrapeCrews which includes fame
}

// ScrapeBattleInfo scrapes only battle info for existing crews
func (s *Scraper) ScrapeBattleInfo() error {
	var crews []models.Crew
	if err := s.db.Where("ocean = ? AND is_active = ?", s.ocean, true).Find(&crews).Error; err != nil {
		return fmt.Errorf("failed to fetch crews: %w", err)
	}

	scrapedAt := time.Now()

	for _, crew := range crews {
		battleURL := GetCrewBattleInfoURL(s.ocean, crew.GameCrewID)
		battleHTML, err := s.fetchHTML(battleURL)
		if err != nil {
			log.Printf("Failed to fetch battle info for crew %d: %v", crew.GameCrewID, err)
			s.job.IncrementFailed()
			continue
		}

		battleData, err := ParseCrewBattleInfo(battleHTML, crew.GameCrewID)
		if err != nil {
			log.Printf("Failed to parse battle info for crew %d: %v", crew.GameCrewID, err)
			s.job.IncrementFailed()
			continue
		}

		// Get previous battle record for delta calculation
		var prevRecord models.CrewBattleRecord
		err = s.db.Where("crew_id = ?", crew.ID).
			Order("scraped_at DESC").First(&prevRecord).Error

		battleRecord := models.CrewBattleRecord{
			CrewID:         crew.ID,
			ScrapedAt:      scrapedAt,
			CrewRank:       battleData.CrewRank,
			TotalPVPWins:   battleData.TotalPVPWins,
			TotalPVPLosses: battleData.TotalPVPLosses,
		}

		if err == nil {
			battleRecord.CalculateDeltas(&prevRecord)
		} else {
			battleRecord.CalculateDeltas(nil)
		}

		if err := s.db.Where("crew_id = ? AND scraped_at = ?", crew.ID, scrapedAt).
			FirstOrCreate(&battleRecord).Error; err != nil {
			log.Printf("Failed to create battle record for crew %d: %v", crew.GameCrewID, err)
			s.job.IncrementFailed()
			continue
		}

		s.job.IncrementProcessed()
		s.db.Save(s.job)
	}

	return nil
}

// ScrapeFlags scrapes all flag data
func (s *Scraper) ScrapeFlags() error {
	// First, get flag list from fame list
	url := GetFlagFameListURL(s.ocean)
	htmlContent, err := s.fetchHTML(url)
	if err != nil {
		return fmt.Errorf("failed to fetch flag fame list: %w", err)
	}

	// Parse flag fame list
	flags, err := ParseFlagFameList(htmlContent, s.ocean)
	if err != nil {
		return fmt.Errorf("failed to parse flag fame list: %w", err)
	}

	if len(flags) == 0 {
		log.Printf("WARNING: Parsed 0 flags for ocean %s from URL: %s", s.ocean, url)
		return fmt.Errorf("no flags found in parsed HTML for ocean %s", s.ocean)
	}

	log.Printf("Successfully parsed %d flags for ocean %s", len(flags), s.ocean)

	scrapedAt := time.Now()

	for _, flagData := range flags {
		if err := s.processFlag(flagData, scrapedAt); err != nil {
			log.Printf("Error processing flag %d: %v", flagData.FlagID, err)
			s.job.IncrementFailed()
			continue
		}
		s.job.IncrementProcessed()
		s.db.Save(s.job)
	}

	return nil
}

// processFlag processes a single flag and all related data
func (s *Scraper) processFlag(fameData FlagFameData, scrapedAt time.Time) error {
	// Fetch flag info page
	url := GetFlagInfoURL(s.ocean, fameData.FlagID)
	flagInfoHTML, err := s.fetchHTML(url)
	if err != nil {
		return fmt.Errorf("failed to fetch flag info: %w", err)
	}

	flagData, err := ParseFlagInfo(flagInfoHTML, fameData.FlagID, s.ocean)
	if err != nil {
		return fmt.Errorf("failed to parse flag info: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Find or create flag
		var flag models.Flag
		err := tx.Where("game_flag_id = ? AND ocean = ?", fameData.FlagID, s.ocean).
			FirstOrCreate(&flag, models.Flag{
				GameFlagID: fameData.FlagID,
				Ocean:      s.ocean,
				Name:       flagData.Name,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to get/create flag: %w", err)
		}

		// Update flag fields
		flag.Name = flagData.Name
		flag.LastSeenAt = scrapedAt

		if err := tx.Save(&flag).Error; err != nil {
			return fmt.Errorf("failed to save flag: %w", err)
		}

		// Create fame record
		fameRecord := models.FlagFameRecord{
			FlagID:    flag.ID,
			ScrapedAt: scrapedAt,
			FameLevel: fameData.FameLevel,
			FameRank:  fameData.Rank,
		}
		if err := tx.Where("flag_id = ? AND scraped_at = ?", flag.ID, scrapedAt).
			FirstOrCreate(&fameRecord).Error; err != nil {
			return fmt.Errorf("failed to create fame record: %w", err)
		}

		return nil
	})
}

// ScrapeFlagFame scrapes only flag fame data
func (s *Scraper) ScrapeFlagFame() error {
	return s.ScrapeFlags() // Reuse ScrapeFlags which includes fame
}
