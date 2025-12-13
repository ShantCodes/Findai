package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	// Cache for loaded SQL queries to avoid disk reads on every call
	queryCache = make(map[string]string)
	cacheMutex sync.RWMutex
)

// loadQueryFromFile loads a SQL query from file with caching
func loadQueryFromFile(queryName string) (string, error) {
	// Add .sql extension if not present
	if !strings.HasSuffix(queryName, ".sql") {
		queryName = queryName + ".sql"
	}

	// Check cache first for performance
	cacheMutex.RLock()
	if query, ok := queryCache[queryName]; ok {
		cacheMutex.RUnlock()
		return query, nil
	}
	cacheMutex.RUnlock()

	// Search for the SQL file in the sql directory and its subdirectories
	basePath := "./src/sql"
	var foundPath string

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, only look for files
		if info.IsDir() {
			return nil
		}

		// Case-insensitive comparison for filename
		if strings.EqualFold(info.Name(), queryName) {
			foundPath = path
			// Stop walking further
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return "", fmt.Errorf("error searching for SQL file: %v", err)
	}

	if foundPath == "" {
		return "", fmt.Errorf("SQL file '%s' not found in %s or subdirectories", queryName, basePath)
	}

	// Read the SQL file
	data, err := os.ReadFile(foundPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SQL file %s: %v", foundPath, err)
	}

	query := strings.TrimSpace(string(data))

	// Cache the query for future use
	cacheMutex.Lock()
	queryCache[queryName] = query
	cacheMutex.Unlock()

	return query, nil
}

// QuerySelect loads a SQL query from file AND executes it, returning a single row
func QuerySelect(ctx context.Context, db *sqlx.DB, queryName string, args ...interface{}) *sqlx.Row {
	// 1. Load the SQL query from file
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		// This is not ideal, but we need to return a *sqlx.Row.
		// We can't return an error directly.
		// The caller will need to check row.Err().
		return &sqlx.Row{}
	}

	// 2. Execute the query and return the row
	return db.QueryRowxContext(ctx, sqlQuery, args...)
}

// QuerySelectRows loads a SQL query from file AND executes it, returning multiple rows
func QuerySelectRows(ctx context.Context, db *sqlx.DB, queryName string, args ...interface{}) (*sqlx.Rows, error) {
	// 1. Load the SQL query from file
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		return nil, fmt.Errorf("failed to load query '%s': %v", queryName, err)
	}

	// 2. Execute the query and return multiple rows
	return db.QueryxContext(ctx, sqlQuery, args...)
}

// ExecSelect loads a SQL query from file AND executes it (for INSERT, UPDATE, DELETE)
func ExecSelect(ctx context.Context, db *sqlx.DB, queryName string, args ...interface{}) (sql.Result, error) {
	// 1. Load the SQL query from file
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		return nil, fmt.Errorf("failed to load query '%s': %v", queryName, err)
	}

	// 2. Execute the query (for INSERT, UPDATE, DELETE operations)
	return db.ExecContext(ctx, sqlQuery, args...)
}

// ClearCache clears the query cache (useful for development when SQL files change)
func ClearCache() {
	cacheMutex.Lock()
	queryCache = make(map[string]string)
	cacheMutex.Unlock()
}