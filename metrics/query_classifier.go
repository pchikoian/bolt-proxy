package metrics

import (
	"strings"
)

// ClassifyQuery determines the CRUD operation type from a query string
func ClassifyQuery(query string) string {
	// Convert to uppercase for easier matching
	upperQuery := strings.ToUpper(strings.TrimSpace(query))

	// Check for Cypher/Neo4j query patterns
	if strings.HasPrefix(upperQuery, "CREATE") || strings.Contains(upperQuery, " CREATE ") {
		return "create"
	}
	if strings.HasPrefix(upperQuery, "MATCH") || strings.Contains(upperQuery, " MATCH ") ||
		strings.HasPrefix(upperQuery, "RETURN") || strings.Contains(upperQuery, " RETURN ") ||
		strings.HasPrefix(upperQuery, "WITH") {
		return "read"
	}
	if strings.HasPrefix(upperQuery, "MERGE") || strings.Contains(upperQuery, " MERGE ") ||
		strings.HasPrefix(upperQuery, "SET") || strings.Contains(upperQuery, " SET ") {
		return "update"
	}
	if strings.HasPrefix(upperQuery, "DELETE") || strings.Contains(upperQuery, " DELETE ") ||
		strings.HasPrefix(upperQuery, "DETACH DELETE") || strings.Contains(upperQuery, " DETACH DELETE ") ||
		strings.HasPrefix(upperQuery, "REMOVE") || strings.Contains(upperQuery, " REMOVE ") {
		return "delete"
	}

	// Check for SQL query patterns (for compatibility)
	if strings.HasPrefix(upperQuery, "INSERT") || strings.HasPrefix(upperQuery, "CREATE TABLE") {
		return "create"
	}
	if strings.HasPrefix(upperQuery, "SELECT") {
		return "read"
	}
	if strings.HasPrefix(upperQuery, "UPDATE") {
		return "update"
	}
	if strings.HasPrefix(upperQuery, "DELETE FROM") || strings.HasPrefix(upperQuery, "DROP") {
		return "delete"
	}

	// Default to "other" for unrecognized patterns
	return "other"
}
