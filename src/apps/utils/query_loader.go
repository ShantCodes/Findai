package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// Database connection (should be set during app initialization)
	db *sql.DB
	dbMutex sync.RWMutex
	
	// Cache for loaded SQL queries to avoid disk reads on every call
	queryCache = make(map[string]string)
	cacheMutex sync.RWMutex
)

// SetDB sets the database connection for QuerySelect to use
func SetDB(database *sql.DB) {
	dbMutex.Lock()
	db = database
	dbMutex.Unlock()
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return db
}

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
func QuerySelect(ctx context.Context, queryName string, args ...interface{}) (*sql.Row, error) {
	// 1. Load the SQL query from file
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		return nil, fmt.Errorf("failed to load query '%s': %v", queryName, err)
	}

	// 2. Get the database connection
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not set. Call utils.SetDB() first")
	}

	// 3. Execute the query and return the row
	// QueryRowContext executes the query and returns a single row
	// The caller will need to call .Scan() on the returned row
	return db.QueryRowContext(ctx, sqlQuery, args...), nil
}

// QuerySelectRows loads a SQL query from file AND executes it, returning multiple rows
func QuerySelectRows(ctx context.Context, queryName string, args ...interface{}) (*sql.Rows, error) {
	// 1. Load the SQL query from file
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		return nil, fmt.Errorf("failed to load query '%s': %v", queryName, err)
	}

	// 2. Get the database connection
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not set. Call utils.SetDB() first")
	}

	// 3. Execute the query and return multiple rows
	// QueryContext executes the query and returns multiple rows
	// The caller will need to iterate over the rows
	return db.QueryContext(ctx, sqlQuery, args...)
}

// ExecSelect loads a SQL query from file AND executes it (for INSERT, UPDATE, DELETE)
func ExecSelect(ctx context.Context, queryName string, args ...interface{}) (sql.Result, error) {
	// 1. Load the SQL query from file
	sqlQuery, err := loadQueryFromFile(queryName)
	if err != nil {
		return nil, fmt.Errorf("failed to load query '%s': %v", queryName, err)
	}

	// 2. Get the database connection
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection not set. Call utils.SetDB() first")
	}

	// 3. Execute the query (for INSERT, UPDATE, DELETE operations)
	// ExecContext executes the query and returns a Result with rows affected, last insert ID, etc.
	return db.ExecContext(ctx, sqlQuery, args...)
}

// ClearCache clears the query cache (useful for development when SQL files change)
func ClearCache() {
	cacheMutex.Lock()
	queryCache = make(map[string]string)
	cacheMutex.Unlock()
}