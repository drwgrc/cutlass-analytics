package repositories

import (
	"strings"

	"gorm.io/gorm"
)

// applySorting applies sorting to a GORM query
func applySorting(query *gorm.DB, sortBy string, sortOrder string) *gorm.DB {
	if sortBy == "" {
		return query
	}

	// Default to descending if not specified
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Sanitize sortBy to prevent SQL injection
	// Only allow alphanumeric, underscore, and dot characters
	sortBy = strings.ReplaceAll(sortBy, ";", "")
	sortBy = strings.ReplaceAll(sortBy, "--", "")
	sortBy = strings.ReplaceAll(sortBy, "/*", "")
	sortBy = strings.ReplaceAll(sortBy, "*/", "")

	orderClause := sortBy
	if sortOrder == "asc" {
		orderClause += " ASC"
	} else {
		orderClause += " DESC"
	}

	return query.Order(orderClause)
}
