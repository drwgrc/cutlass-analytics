package scraper

import (
	"fmt"
	"cutlass_analytics/internal/types"
)

// GetIslandListURL returns the URL for the island list page showing all islands
func GetIslandListURL(ocean types.Ocean) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/island/info.wm?showAll=true", ocean)
}

// GetIslandInfoURL returns the URL for a specific island's info page
func GetIslandInfoURL(ocean types.Ocean, islandID uint64) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/island/info.wm?islandid=%d", ocean, islandID)
}

// GetTaxRatesURL returns the URL for the tax rates page
func GetTaxRatesURL(ocean types.Ocean) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/econ/taxrates.wm", ocean)
}

// GetCrewInfoURL returns the URL for a specific crew's info page
func GetCrewInfoURL(ocean types.Ocean, crewID uint64) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/crew/info.wm?crewid=%d", ocean, crewID)
}

// GetCrewBattleInfoURL returns the URL for a specific crew's battle info page
func GetCrewBattleInfoURL(ocean types.Ocean, crewID uint64) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/crew/battleinfo.wm?crewid=%d&classic=$classic", ocean, crewID)
}

// GetCrewFameListURL returns the URL for the crew fame list
func GetCrewFameListURL(ocean types.Ocean) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/ratings/top_fame_97.html", ocean)
}

// GetFlagInfoURL returns the URL for a specific flag's info page
func GetFlagInfoURL(ocean types.Ocean, flagID uint64) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/yoweb/flag/info.wm?flagid=%d", ocean, flagID)
}

// GetFlagFameListURL returns the URL for the flag fame list
func GetFlagFameListURL(ocean types.Ocean) string {
	return fmt.Sprintf("https://%s.puzzlepirates.com/ratings/top_fame_112.html", ocean)
}
