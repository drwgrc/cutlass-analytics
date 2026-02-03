package poller

import (
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CSVPoller handles fetching and importing market order data from the Puzzle Pirates buysell endpoint
type CSVPoller struct {
	db      *gorm.DB
	client  *http.Client
	oceans  []types.Ocean
}

// NewCSVPoller creates a new CSV poller instance
func NewCSVPoller(db *gorm.DB, oceans []types.Ocean) *CSVPoller {
	return &CSVPoller{
		db:     db,
		client: &http.Client{Timeout: 30 * time.Second},
		oceans: oceans,
	}
}

// getBuySellURL returns the buysell CSV URL for a given ocean
func getBuySellURL(ocean types.Ocean) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/econ/buysell.wm", ocean)
}

// Run fetches and imports market orders from all configured oceans
func (p *CSVPoller) Run() error {
	log.Println("CSV poller: Starting market order import...")

	var allOrders []models.MarketOrder
	now := time.Now()

	for _, ocean := range p.oceans {
		orders, err := p.fetchAndParse(ocean, now)
		if err != nil {
			log.Printf("CSV poller: Error fetching %s ocean: %v", ocean, err)
			continue
		}
		log.Printf("CSV poller: Fetched %d orders from %s ocean", len(orders), ocean)
		allOrders = append(allOrders, orders...)
	}

	if len(allOrders) == 0 {
		log.Println("CSV poller: No orders fetched from any ocean")
		return nil
	}

	// Import orders atomically
	if err := p.importOrders(allOrders); err != nil {
		return fmt.Errorf("failed to import orders: %w", err)
	}

	log.Printf("CSV poller: Successfully imported %d market orders", len(allOrders))
	return nil
}

// fetchAndParse fetches the CSV from a specific ocean and parses it
func (p *CSVPoller) fetchAndParse(ocean types.Ocean, importTime time.Time) ([]models.MarketOrder, error) {
	url := getBuySellURL(ocean)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CSV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return p.parseCSV(resp.Body, ocean, importTime)
}

// parseCSV reads and parses CSV data from a reader
// Expected format: island,commodity,store,buy_price,buy_quantity,sell_price,sell_quantity
func (p *CSVPoller) parseCSV(r io.Reader, ocean types.Ocean, importTime time.Time) ([]models.MarketOrder, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, nil // Empty file or only headers
	}

	orders := make([]models.MarketOrder, 0, len(records)-1)

	// Skip header row
	for i, record := range records[1:] {
		order, err := p.parseRecord(record, ocean, importTime)
		if err != nil {
			log.Printf("CSV poller: Error parsing row %d: %v", i+2, err)
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// parseRecord converts a CSV record to a MarketOrder
// Format: island,commodity,store,buy_price,buy_quantity,sell_price,sell_quantity
func (p *CSVPoller) parseRecord(record []string, ocean types.Ocean, importTime time.Time) (models.MarketOrder, error) {
	if len(record) < 7 {
		return models.MarketOrder{}, fmt.Errorf("invalid record length: expected 7 fields, got %d", len(record))
	}

	islandName := strings.TrimSpace(record[0])
	if islandName == "" {
		return models.MarketOrder{}, fmt.Errorf("island name is empty")
	}

	commodityName := strings.TrimSpace(record[1])
	if commodityName == "" {
		return models.MarketOrder{}, fmt.Errorf("commodity name is empty")
	}

	shopName := strings.TrimSpace(record[2])
	if shopName == "" {
		return models.MarketOrder{}, fmt.Errorf("shop name is empty")
	}

	buyPrice, err := strconv.Atoi(strings.TrimSpace(record[3]))
	if err != nil {
		return models.MarketOrder{}, fmt.Errorf("invalid buy_price: %w", err)
	}

	buyQuantity, err := strconv.Atoi(strings.TrimSpace(record[4]))
	if err != nil {
		return models.MarketOrder{}, fmt.Errorf("invalid buy_quantity: %w", err)
	}

	sellPrice, err := strconv.Atoi(strings.TrimSpace(record[5]))
	if err != nil {
		return models.MarketOrder{}, fmt.Errorf("invalid sell_price: %w", err)
	}

	sellQuantity, err := strconv.Atoi(strings.TrimSpace(record[6]))
	if err != nil {
		return models.MarketOrder{}, fmt.Errorf("invalid sell_quantity: %w", err)
	}

	return models.MarketOrder{
		Ocean:         ocean,
		IslandName:    islandName,
		CommodityName: commodityName,
		ShopName:      shopName,
		BuyPrice:      buyPrice,
		BuyQuantity:   buyQuantity,
		SellPrice:     sellPrice,
		SellQuantity:  sellQuantity,
		ImportedAt:    importTime,
	}, nil
}

// importOrders replaces all existing market orders with new ones atomically
func (p *CSVPoller) importOrders(orders []models.MarketOrder) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		// Delete all existing market orders
		if err := tx.Exec("DELETE FROM market_orders").Error; err != nil {
			return fmt.Errorf("failed to clear existing orders: %w", err)
		}

		// Batch insert new orders for better performance
		batchSize := 100
		for i := 0; i < len(orders); i += batchSize {
			end := i + batchSize
			if end > len(orders) {
				end = len(orders)
			}
			batch := orders[i:end]

			if err := tx.Create(&batch).Error; err != nil {
				return fmt.Errorf("failed to insert orders batch: %w", err)
			}
		}

		return nil
	})
}
