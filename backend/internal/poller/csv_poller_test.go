package poller

import (
	"cutlass_analytics/internal/models"
	"cutlass_analytics/internal/types"
	"strings"
	"testing"
	"time"
)

func TestParseRecord_ValidRecord(t *testing.T) {
	p := &CSVPoller{}
	now := time.Now()
	ocean := types.OceanEmerald

	tests := []struct {
		name     string
		record   []string
		expected models.MarketOrder
	}{
		{
			name:   "basic buy order",
			record: []string{"Maia-Insel", "Sugar cane", "Ferklstall", "4", "100", "0", "0"},
			expected: models.MarketOrder{
				Ocean:         ocean,
				IslandName:    "Maia-Insel",
				CommodityName: "Sugar cane",
				ShopName:      "Ferklstall",
				BuyPrice:      4,
				BuyQuantity:   100,
				SellPrice:     0,
				SellQuantity:  0,
			},
		},
		{
			name:   "buy and sell order",
			record: []string{"Chachapoya-Insel", "Iron", "Deinsklave's Schmiede-Laden", "12", "0", "40", "450"},
			expected: models.MarketOrder{
				Ocean:         ocean,
				IslandName:    "Chachapoya-Insel",
				CommodityName: "Iron",
				ShopName:      "Deinsklave's Schmiede-Laden",
				BuyPrice:      12,
				BuyQuantity:   0,
				SellPrice:     40,
				SellQuantity:  450,
			},
		},
		{
			name:   "sell only order",
			record: []string{"Maia-Insel", "Emeralds", "Ferklnachschub", "0", "0", "470", "4"},
			expected: models.MarketOrder{
				Ocean:         ocean,
				IslandName:    "Maia-Insel",
				CommodityName: "Emeralds",
				ShopName:      "Ferklnachschub",
				BuyPrice:      0,
				BuyQuantity:   0,
				SellPrice:     470,
				SellQuantity:  4,
			},
		},
		{
			name:   "whitespace trimming",
			record: []string{" Maia-Insel ", " Wood ", " Test Shop ", " 10 ", " 50 ", " 30 ", " 0 "},
			expected: models.MarketOrder{
				Ocean:         ocean,
				IslandName:    "Maia-Insel",
				CommodityName: "Wood",
				ShopName:      "Test Shop",
				BuyPrice:      10,
				BuyQuantity:   50,
				SellPrice:     30,
				SellQuantity:  0,
			},
		},
		{
			name:   "large quantities",
			record: []string{"TestIsland", "TestCommodity", "TestShop", "15", "1001", "19", "72"},
			expected: models.MarketOrder{
				Ocean:         ocean,
				IslandName:    "TestIsland",
				CommodityName: "TestCommodity",
				ShopName:      "TestShop",
				BuyPrice:      15,
				BuyQuantity:   1001,
				SellPrice:     19,
				SellQuantity:  72,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := p.parseRecord(tt.record, ocean, now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if order.Ocean != tt.expected.Ocean {
				t.Errorf("Ocean = %s, want %s", order.Ocean, tt.expected.Ocean)
			}
			if order.IslandName != tt.expected.IslandName {
				t.Errorf("IslandName = %s, want %s", order.IslandName, tt.expected.IslandName)
			}
			if order.CommodityName != tt.expected.CommodityName {
				t.Errorf("CommodityName = %s, want %s", order.CommodityName, tt.expected.CommodityName)
			}
			if order.ShopName != tt.expected.ShopName {
				t.Errorf("ShopName = %s, want %s", order.ShopName, tt.expected.ShopName)
			}
			if order.BuyPrice != tt.expected.BuyPrice {
				t.Errorf("BuyPrice = %d, want %d", order.BuyPrice, tt.expected.BuyPrice)
			}
			if order.BuyQuantity != tt.expected.BuyQuantity {
				t.Errorf("BuyQuantity = %d, want %d", order.BuyQuantity, tt.expected.BuyQuantity)
			}
			if order.SellPrice != tt.expected.SellPrice {
				t.Errorf("SellPrice = %d, want %d", order.SellPrice, tt.expected.SellPrice)
			}
			if order.SellQuantity != tt.expected.SellQuantity {
				t.Errorf("SellQuantity = %d, want %d", order.SellQuantity, tt.expected.SellQuantity)
			}
			if order.ImportedAt != now {
				t.Errorf("ImportedAt not set correctly")
			}
		})
	}
}

func TestParseRecord_InvalidRecord(t *testing.T) {
	p := &CSVPoller{}
	now := time.Now()
	ocean := types.OceanEmerald

	tests := []struct {
		name   string
		record []string
	}{
		{
			name:   "too few fields",
			record: []string{"Island", "Commodity", "Shop", "100", "50", "0"},
		},
		{
			name:   "empty island name",
			record: []string{"", "Commodity", "Shop", "100", "50", "0", "0"},
		},
		{
			name:   "empty commodity name",
			record: []string{"Island", "", "Shop", "100", "50", "0", "0"},
		},
		{
			name:   "empty shop name",
			record: []string{"Island", "Commodity", "", "100", "50", "0", "0"},
		},
		{
			name:   "invalid buy_price",
			record: []string{"Island", "Commodity", "Shop", "abc", "50", "0", "0"},
		},
		{
			name:   "invalid buy_quantity",
			record: []string{"Island", "Commodity", "Shop", "100", "abc", "0", "0"},
		},
		{
			name:   "invalid sell_price",
			record: []string{"Island", "Commodity", "Shop", "100", "50", "abc", "0"},
		},
		{
			name:   "invalid sell_quantity",
			record: []string{"Island", "Commodity", "Shop", "100", "50", "0", "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := p.parseRecord(tt.record, ocean, now)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestParseCSV_ValidData(t *testing.T) {
	p := &CSVPoller{}
	now := time.Now()
	ocean := types.OceanEmerald

	csvData := `island,commodity,store,buy_price,buy_quantity,sell_price,sell_quantity
Maia-Insel,Sugar cane,Ferklstall,4,100,0,0
Maia-Insel,Iron,Karlimero's Schmiede-Laden,11,1001,20,0
Chachapoya-Insel,Hemp,Powerolli's Ausstatter-Laden,2,0,10,30`

	reader := strings.NewReader(csvData)
	orders, err := p.parseCSV(reader, ocean, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got %d", len(orders))
	}

	// Verify first order
	if orders[0].IslandName != "Maia-Insel" || orders[0].CommodityName != "Sugar cane" {
		t.Errorf("first order incorrect: %+v", orders[0])
	}
	if orders[0].BuyPrice != 4 || orders[0].BuyQuantity != 100 {
		t.Errorf("first order buy data incorrect: %+v", orders[0])
	}
	if orders[0].Ocean != ocean {
		t.Errorf("first order ocean incorrect: %s, want %s", orders[0].Ocean, ocean)
	}

	// Verify second order
	if orders[1].ShopName != "Karlimero's Schmiede-Laden" {
		t.Errorf("second order shop name incorrect: %s", orders[1].ShopName)
	}
	if orders[1].BuyPrice != 11 || orders[1].SellPrice != 20 {
		t.Errorf("second order prices incorrect: %+v", orders[1])
	}

	// Verify third order (sell only)
	if orders[2].BuyQuantity != 0 || orders[2].SellQuantity != 30 {
		t.Errorf("third order quantities incorrect: %+v", orders[2])
	}
}

func TestParseCSV_EmptyData(t *testing.T) {
	p := &CSVPoller{}
	now := time.Now()
	ocean := types.OceanEmerald

	// Data with only header
	csvData := `island,commodity,store,buy_price,buy_quantity,sell_price,sell_quantity`

	reader := strings.NewReader(csvData)
	orders, err := p.parseCSV(reader, ocean, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(orders))
	}
}

func TestGetBuySellURL(t *testing.T) {
	tests := []struct {
		ocean    types.Ocean
		expected string
	}{
		{types.OceanEmerald, "https://emerald.puzzlepirates.com/yoweb/econ/buysell.wm"},
		{types.OceanMeridian, "https://meridian.puzzlepirates.com/yoweb/econ/buysell.wm"},
		{types.OceanCerulean, "https://cerulean.puzzlepirates.com/yoweb/econ/buysell.wm"},
	}

	for _, tt := range tests {
		t.Run(string(tt.ocean), func(t *testing.T) {
			url := getBuySellURL(tt.ocean)
			if url != tt.expected {
				t.Errorf("getBuySellURL(%s) = %s, want %s", tt.ocean, url, tt.expected)
			}
		})
	}
}

func TestNewCSVPoller(t *testing.T) {
	oceans := []types.Ocean{types.OceanEmerald, types.OceanMeridian}
	p := NewCSVPoller(nil, oceans)

	if p == nil {
		t.Fatal("expected non-nil poller")
	}
	if len(p.oceans) != 2 {
		t.Errorf("expected 2 oceans, got %d", len(p.oceans))
	}
	if p.client == nil {
		t.Error("expected non-nil HTTP client")
	}
}
