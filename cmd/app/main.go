package main

import (
	"log"
	"findai/src/apps"
	"findai/src/config"
	"time"

	database "github.com/socious-io/pkg_database"
)

func main() {
	if _, err := config.Init("config.yml"); err != nil {
		log.Fatalf("Failed to initialize config: %v. Make sure config.yml exists and is correctly formatted.", err)
	}

	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 50,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	apps.Serve()
}

