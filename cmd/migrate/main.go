package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"findai/src/config"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if _, err := config.Init("config.yml"); err != nil {
		log.Fatalf("Failed to initialize config: %v. Make sure config.yml exists and is correctly formatted.", err)
	}

	// Check command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Expected 'new {name}' or 'up' command.")
	}

	command := os.Args[1]

	migrationPath := fmt.Sprintf("file://%s", config.Config.Database.Migrations)

	switch command {
	case "new":
		if len(os.Args) < 3 {
			log.Fatal("Expected a name for the migration.")
		}
		name := os.Args[2]
		err := createMigration(name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migration file created successfully!")
	case "up":
		m, err := migrate.New(migrationPath, config.Config.Database.URL)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations applied successfully!")
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Expected a version number for the 'force' command.")
		}
		version := os.Args[2]
		m, err := migrate.New(migrationPath, config.Config.Database.URL)
		if err != nil {
			log.Fatal(err)
		}
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			log.Fatal("Invalid version number.")
		}
		if err := m.Force(versionInt); err != nil {
			log.Fatal(err)
		}
		log.Printf("Forced migration to version %d", versionInt)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func createMigration(name string) error {
	timestamp := time.Now().Format("20060102150405") // e.g., 20230807123456
	upFilename := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	downFilename := fmt.Sprintf("%s_%s.down.sql", timestamp, name)

	migrationsDir := config.Config.Database.Migrations

	// Ensure the migrations directory exists
	err := os.MkdirAll(migrationsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	// Create the up file
	upFilePath := filepath.Join(migrationsDir, upFilename)
	upFile, err := os.Create(upFilePath)
	if err != nil {
		return fmt.Errorf("failed to create up migration file: %v", err)
	}
	defer upFile.Close()

	// Create the down file
	downFilePath := filepath.Join(migrationsDir, downFilename)
	downFile, err := os.Create(downFilePath)
	if err != nil {
		return fmt.Errorf("failed to create down migration file: %v", err)
	}
	defer downFile.Close()

	log.Printf("Created migration files: %s, %s", upFilePath, downFilePath)

	return nil
}
