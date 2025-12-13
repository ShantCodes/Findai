package main

import (
	"findai/src/apps"
	"findai/src/config"
	"findai/src/database"
	"log"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	DB *sqlx.DB
}

func main() {
	if _, err := config.Init("config.yml"); err != nil {
		log.Fatalf("Failed to initialize config: %v. Make sure config.yml exists and is correctly formatted.", err)
	}

	database.Connect(config.Config.Database.URL)
	server := &Server{DB: database.DB()}

	apps.Serve(server.DB)
}

