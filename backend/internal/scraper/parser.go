package scraper

import (
	"cutlass_analytics/internal/types"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// IslandData represents parsed island information
type IslandData struct {
	GameIslandID  uint64
	Name          string
	Archipelago   string
	Size          types.IslandSize
	IsColonized   bool
	Population    int
	GovernorFlag  string
	GovernorName  string
	Commodities   []string
}

// TaxRateData represents parsed tax rate information
type TaxRateData struct {
	CommodityName string
	TaxValue      int
}

// CrewFameData represents parsed crew fame list entry
type CrewFameData struct {
	CrewID    uint64
	Name      string
	FameLevel types.FameLevel
	Rank      *int
}

// CrewData represents parsed crew information
type CrewData struct {
	GameCrewID uint64
	Name       string
	FlagID     *uint64
	FlagName   string
}

// CrewBattleData represents parsed crew battle information
type CrewBattleData struct {
	CrewRank      types.CrewRank
	TotalPVPWins  int
	TotalPVPLosses int
}

// FlagFameData represents parsed flag fame list entry
type FlagFameData struct {
	FlagID    uint64
	Name      string
	FameLevel types.FameLevel
	Rank      *int
}

// FlagData represents parsed flag information
type FlagData struct {
	GameFlagID uint64
	Name       string
}


// ParseIslandList parses the island list page (showAll=true)
func ParseIslandList(html string, ocean types.Ocean) ([]IslandData, error) {
	var islands []IslandData

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return islands, fmt.Errorf("ParseIslandList: failed to parse HTML: %w", err)
	}

	// Structure: main <center> in body, each island in its own <center> tag
	doc.Find("body > center").Each(func(_ int, mainCenter *goquery.Selection) {
		// Find all island <center> tags within the main center
		mainCenter.Find("center").Each(func(_ int, islandCenter *goquery.Selection) {
			var island IslandData
			island.IsColonized = true // All islands on this page are colonized

			// Extract island name from <font> tag
			islandCenter.Find("font").Each(func(_ int, font *goquery.Selection) {
				if island.Name == "" {
					island.Name = strings.TrimSpace(font.Text())
				}
			})

			// Extract full text content for regex matching
			fullText := islandCenter.Text()

			// Extract population from "Population: [number]"
			popRe := regexp.MustCompile(`(?i)population[:\s]+(\d+)`)
			popMatches := popRe.FindStringSubmatch(fullText)
			if len(popMatches) > 1 {
				if pop, err := strconv.Atoi(strings.ReplaceAll(popMatches[1], ",", "")); err == nil {
					island.Population = pop
				}
			}

			// Extract archipelago from "Located in the X archipelago."
			archRe := regexp.MustCompile(`(?i)located\s+in\s+the\s+([A-Za-z\s]+?)\s+archipelago`)
			archMatches := archRe.FindStringSubmatch(fullText)
			if len(archMatches) > 1 {
				island.Archipelago = strings.TrimSpace(archMatches[1])
			}

			// Extract commodities from "Exports: [Commodities[]]" pattern
			// Commodities are separated by space and comma
			exportsRe := regexp.MustCompile(`(?i)exports[:\s]+(.+?)(?:\n|$)`)
			exportsMatches := exportsRe.FindStringSubmatch(fullText)
			if len(exportsMatches) > 1 {
				exportsText := strings.TrimSpace(exportsMatches[1])
				// Split by comma and space, then clean up each commodity
				commodityParts := strings.Split(exportsText, ", ")
				for _, part := range commodityParts {
					comm := strings.TrimSpace(part)
					if comm != "" {
						island.Commodities = append(island.Commodities, comm)
					}
				}
			}

			// Try to extract island ID from any link with islandid parameter (if name is a link)
			islandCenter.Find("a[href*='islandid=']").Each(func(_ int, link *goquery.Selection) {
				href, exists := link.Attr("href")
				if exists {
					re := regexp.MustCompile(`islandid=(\d+)`)
					matches := re.FindStringSubmatch(href)
					if len(matches) > 1 {
						if id, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
							island.GameIslandID = id
						}
					}
				}
			})

			// Add island if we have at least a name
			if island.Name != "" {
				islands = append(islands, island)
			}
		})
	})

	return islands, nil
}

// ParseIslandInfo parses an individual island info page
func ParseIslandInfo(html string, islandID uint64, ocean types.Ocean) (*IslandData, error) {
	// #region agent log
	logEntry := map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "A",
		"location":     "parser.go:147",
		"message":      "ParseIslandInfo entry",
		"data": map[string]interface{}{
			"islandID":    islandID,
			"htmlLength":  len(html),
			"hasNewline": strings.Contains(html, "\n"),
			"hasTab":     strings.Contains(html, "\t"),
			"hasCR":      strings.Contains(html, "\r"),
		},
		"timestamp": time.Now().UnixMilli(),
	}
	if logBytes, err := json.Marshal(logEntry); err == nil {
		if f, err := os.OpenFile("/Users/andrewgarcia/projects/cutlass-analytics/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			f.Write(logBytes)
			f.WriteString("\n")
			f.Close()
		}
	}
	// #endregion

	// #region agent log
	controlCharCount := 0
	for _, r := range html {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			controlCharCount++
		}
	}
	logEntry2 := map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "B",
		"location":     "parser.go:147",
		"message":      "Control character check",
		"data": map[string]interface{}{
			"islandID":          islandID,
			"controlCharCount":   controlCharCount,
			"first100Chars":     html[:min(100, len(html))],
		},
		"timestamp": time.Now().UnixMilli(),
	}
	if logBytes, err := json.Marshal(logEntry2); err == nil {
		if f, err := os.OpenFile("/Users/andrewgarcia/projects/cutlass-analytics/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			f.Write(logBytes)
			f.WriteString("\n")
			f.Close()
		}
	}
	// #endregion

	island := &IslandData{
		GameIslandID: islandID,
		IsColonized:  true, // If we get here, the island is colonized
	}

	// Parse HTML using goquery (designed for parsing HTML strings directly)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		// #region agent log
		logEntryErr := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "E",
			"location":     "parser.go:221",
			"message":      "goquery parse error",
			"data": map[string]interface{}{
				"islandID":   islandID,
				"error":      err.Error(),
				"htmlLength": len(html),
			},
			"timestamp": time.Now().UnixMilli(),
		}
		if logBytes, logErr := json.Marshal(logEntryErr); logErr == nil {
			if f, logErr := os.OpenFile("/Users/andrewgarcia/projects/cutlass-analytics/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
				f.Write(logBytes)
				f.WriteString("\n")
				f.Close()
			}
		}
		// #endregion
		return island, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract island name from <font size="+1"> tag
	doc.Find("font[size='+1']").Each(func(_ int, s *goquery.Selection) {
		if island.Name == "" {
			island.Name = strings.TrimSpace(s.Text())
		}
	})

	// Extract population from "Population: [number]"
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		re := regexp.MustCompile(`(?i)population[:\s]+(\d+)`)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			if pop, err := strconv.Atoi(strings.ReplaceAll(matches[1], ",", "")); err == nil {
				island.Population = pop
			}
		}
	})

	// Extract archipelago from "Located in the [name] archipelago."
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		re := regexp.MustCompile(`(?i)located\s+in\s+the\s+([A-Za-z\s]+?)\s+archipelago`)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			island.Archipelago = strings.TrimSpace(matches[1])
		}
	})

	// Extract governor name from "Governor: <a>...</a>"
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(strings.ToLower(text), "governor:") {
			s.Find("a[href*='pirate.wm']").Each(func(_ int, link *goquery.Selection) {
				if island.GovernorName == "" {
					island.GovernorName = strings.TrimSpace(link.Text())
				}
			})
		}
	})

	// Extract flag ID from "Ruled by <a href="/yoweb/flag/info.wm?flagid=...">...</a>"
	doc.Find("a[href*='flag/info.wm']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			re := regexp.MustCompile(`flagid=(\d+)`)
			matches := re.FindStringSubmatch(href)
			if len(matches) > 1 {
				if id, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
					island.GovernorFlag = fmt.Sprintf("%d", id)
				}
			}
		}
	})

	// Extract commodities from "Exports:" section
	bodyText := doc.Find("body").Text()
	exportsIdx := strings.Index(strings.ToLower(bodyText), "exports:")
	if exportsIdx >= 0 {
		exportsSection := bodyText[exportsIdx+8:] // Skip "Exports:"
		// Find the end of exports section (before next major section like "Colonized islands" or double newline)
		endIdx := strings.Index(strings.ToLower(exportsSection), "colonized islands")
		if endIdx < 0 {
			// Look for double newline or <br> tag
			endIdx = strings.Index(exportsSection, "\n\n")
		}
		if endIdx < 0 {
			endIdx = len(exportsSection)
		}
		exportsText := exportsSection[:endIdx]
		// Split by comma and clean up each commodity
		commodityParts := strings.Split(exportsText, ",")
		for _, part := range commodityParts {
			comm := strings.TrimSpace(part)
			// Remove any extra whitespace/newlines and normalize
			comm = regexp.MustCompile(`\s+`).ReplaceAllString(comm, " ")
			comm = strings.TrimSpace(comm)
			if comm != "" {
				island.Commodities = append(island.Commodities, comm)
			}
		}
	}

	// #region agent log
	logEntry3 := map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": "C",
		"location":     "parser.go:241",
		"message":      "After goquery parsing",
		"data": map[string]interface{}{
			"islandID":      islandID,
			"islandName":    island.Name,
			"population":    island.Population,
			"commodityCount": len(island.Commodities),
		},
		"timestamp": time.Now().UnixMilli(),
	}
	if logBytes, err2 := json.Marshal(logEntry3); err2 == nil {
		if f, err2 := os.OpenFile("/Users/andrewgarcia/projects/cutlass-analytics/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			f.Write(logBytes)
			f.WriteString("\n")
			f.Close()
		}
	}
	// #endregion

	return island, err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseTaxRates parses the tax rates page
// HTML structure: body > center > table (with 3 tds) > 1st td > table > tr (with th header, then tr with 2 tds)
func ParseTaxRates(html string, ocean types.Ocean) ([]TaxRateData, error) {
	var rates []TaxRateData

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return rates, fmt.Errorf("ParseTaxRates: failed to parse HTML: %w", err)
	}

	// Navigate to the nested table: body > center > table > first td > table
	doc.Find("body > center > table > td:first-child > table").Each(func(_ int, e *goquery.Selection) {
		e.Find("tr").Each(func(_ int, row *goquery.Selection) {
			// Skip header row
			if row.Find("th").Length() > 0 {
				return
			}

			var rate TaxRateData
			var cells []string
			row.Find("td").Each(func(_ int, cell *goquery.Selection) {
				cells = append(cells, strings.TrimSpace(cell.Text()))
			})

			if len(cells) >= 2 {
				rate.CommodityName = cells[0]
				// Extract tax value (remove any formatting)
				taxStr := strings.ReplaceAll(cells[1], ",", "")
				taxStr = regexp.MustCompile(`[^\d]`).ReplaceAllString(taxStr, "")
				if tax, err := strconv.Atoi(taxStr); err == nil {
					rate.TaxValue = tax
					rates = append(rates, rate)
				}
			}
		})
	})

	return rates, nil
}

// ParseCrewFameList parses the crew fame list page
// HTML structure: table > tr (first tr has th header) > tr (data rows with 3 tds)
//   - td[0]: Rank (integer)
//   - td[1]: <a href="/yoweb/crew/info.wm?crewid=CrewID">Name</a>
//   - td[2]: FameLevel (string matching FameLevel enum)
func ParseCrewFameList(html string, ocean types.Ocean) ([]CrewFameData, error) {
	var crews []CrewFameData

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return crews, fmt.Errorf("ParseCrewFameList: failed to parse HTML: %w", err)
	}

	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(_ int, row *goquery.Selection) {
			// Skip header row (contains th elements)
			if row.Find("th").Length() > 0 {
				return
			}

			cells := row.Find("td")
			// Expect exactly 3 cells: Rank, Name (with link), FameLevel
			if cells.Length() < 3 {
				return
			}

			var crew CrewFameData

			// td[0]: Rank
			rankStr := strings.TrimSpace(cells.Eq(0).Text())
			if rank, err := strconv.Atoi(rankStr); err == nil {
				crew.Rank = &rank
			}

			// td[1]: <a href="/yoweb/crew/info.wm?crewid=CrewID">Name</a>
			cells.Eq(1).Find("a").Each(func(_ int, link *goquery.Selection) {
				href, exists := link.Attr("href")
				if exists && strings.Contains(href, "crewid=") {
					re := regexp.MustCompile(`crewid=(\d+)`)
					matches := re.FindStringSubmatch(href)
					if len(matches) > 1 {
						if id, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
							crew.CrewID = id
							crew.Name = strings.TrimSpace(link.Text())
						}
					}
				}
			})

			// td[2]: FameLevel
			fameLevelText := strings.TrimSpace(cells.Eq(2).Text())
			fameLevelLower := strings.ToLower(fameLevelText)
			for _, level := range []types.FameLevel{
				types.FameLevelIllustrious, types.FameLevelRenowned, types.FameLevelEminent,
				types.FameLevelCelebrated, types.FameLevelDistinguished, types.FameLevelRecognized,
				types.FameLevelNoted, types.FameLevelRumored, types.FameLevelObscure,
			} {
				if strings.Contains(fameLevelLower, strings.ToLower(string(level))) {
					crew.FameLevel = level
					break
				}
			}

			if crew.CrewID > 0 && crew.Name != "" {
				crews = append(crews, crew)
			}
		})
	})

	return crews, nil
}

// ParseCrewInfo parses a crew info page
func ParseCrewInfo(html string, crewID uint64, ocean types.Ocean) (*CrewData, error) {
	crew := &CrewData{
		GameCrewID: crewID,
	}

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return crew, fmt.Errorf("ParseCrewInfo: failed to parse HTML: %w", err)
	}

	// Extract crew name
	doc.Find("h1, h2, title").Each(func(_ int, e *goquery.Selection) {
		if crew.Name == "" {
			text := strings.TrimSpace(e.Text())
			if text != "" && !strings.Contains(strings.ToLower(text), "puzzle pirates") {
				crew.Name = text
			}
		}
	})

	// Extract flag ID from links
	doc.Find("a[href*='flagid=']").Each(func(_ int, e *goquery.Selection) {
		href, exists := e.Attr("href")
		if exists {
			re := regexp.MustCompile(`flagid=(\d+)`)
			matches := re.FindStringSubmatch(href)
			if len(matches) > 1 {
				if id, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
					flagID := id
					crew.FlagID = &flagID
					crew.FlagName = strings.TrimSpace(e.Text())
				}
			}
		}
	})

	return crew, nil
}

// ParseCrewBattleInfo parses a crew battle info page
func ParseCrewBattleInfo(html string, crewID uint64) (*CrewBattleData, error) {
	battle := &CrewBattleData{}

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return battle, fmt.Errorf("ParseCrewBattleInfo: failed to parse HTML: %w", err)
	}

	// Extract crew rank
	doc.Find("*").Each(func(_ int, e *goquery.Selection) {
		text := strings.ToLower(e.Text())
		for _, rank := range []types.CrewRank{
			types.CrewRankImperials, types.CrewRankSeaLords, types.CrewRankDreadPirates,
			types.CrewRankBlaggards, types.CrewRankScoundrels, types.CrewRankScurvyDogs,
			types.CrewRankMostlyHarmless, types.CrewRankSailors,
		} {
			if strings.Contains(text, strings.ToLower(string(rank))) {
				battle.CrewRank = rank
				return
			}
		}
	})

	// Extract PVP wins/losses
	doc.Find("*").Each(func(_ int, e *goquery.Selection) {
		text := e.Text()
		// Look for patterns like "Wins: 123" or "PVP Wins: 123"
		re := regexp.MustCompile(`(?i)(?:pvp\s+)?wins?[:\s]+(\d+)`)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			if wins, err := strconv.Atoi(matches[1]); err == nil {
				battle.TotalPVPWins = wins
			}
		}

		re = regexp.MustCompile(`(?i)(?:pvp\s+)?losses?[:\s]+(\d+)`)
		matches = re.FindStringSubmatch(text)
		if len(matches) > 1 {
			if losses, err := strconv.Atoi(matches[1]); err == nil {
				battle.TotalPVPLosses = losses
			}
		}
	})

	return battle, nil
}

// ParseFlagFameList parses the flag fame list page
func ParseFlagFameList(html string, ocean types.Ocean) ([]FlagFameData, error) {
	var flags []FlagFameData

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return flags, fmt.Errorf("ParseFlagFameList: failed to parse HTML: %w", err)
	}

	doc.Find("table").Each(func(_ int, e *goquery.Selection) {
		e.Find("tr").Each(func(_ int, row *goquery.Selection) {
			// Skip header row
			if row.Find("th").Length() > 0 {
				return
			}

			var flag FlagFameData
			var rankStr string

			// Extract rank (usually first cell)
			row.Find("td").Each(func(i int, cell *goquery.Selection) {
				if i == 0 {
					rankStr = strings.TrimSpace(cell.Text())
					if rank, err := strconv.Atoi(rankStr); err == nil {
						flag.Rank = &rank
					}
				}
			})

			// Extract flag ID and name from links
			row.Find("a[href*='flagid=']").Each(func(_ int, link *goquery.Selection) {
				href, exists := link.Attr("href")
				if exists {
					re := regexp.MustCompile(`flagid=(\d+)`)
					matches := re.FindStringSubmatch(href)
					if len(matches) > 1 {
						if id, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
							flag.FlagID = id
							flag.Name = strings.TrimSpace(link.Text())
						}
					}
				}
			})

			// Extract fame level from text
			rowText := strings.ToLower(row.Text())
			for _, level := range []types.FameLevel{
				types.FameLevelIllustrious, types.FameLevelRenowned, types.FameLevelEminent,
				types.FameLevelCelebrated, types.FameLevelDistinguished, types.FameLevelRecognized,
				types.FameLevelNoted, types.FameLevelRumored, types.FameLevelObscure,
			} {
				if strings.Contains(rowText, strings.ToLower(string(level))) {
					flag.FameLevel = level
					break
				}
			}

			if flag.FlagID > 0 && flag.Name != "" {
				flags = append(flags, flag)
			}
		})
	})

	return flags, nil
}

// ParseFlagInfo parses a flag info page
func ParseFlagInfo(html string, flagID uint64, ocean types.Ocean) (*FlagData, error) {
	flag := &FlagData{
		GameFlagID: flagID,
	}

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return flag, fmt.Errorf("ParseFlagInfo: failed to parse HTML: %w", err)
	}

	// Extract flag name
	doc.Find("h1, h2, title").Each(func(_ int, e *goquery.Selection) {
		if flag.Name == "" {
			text := strings.TrimSpace(e.Text())
			if text != "" && !strings.Contains(strings.ToLower(text), "puzzle pirates") {
				flag.Name = text
			}
		}
	})

	return flag, nil
}
