package main


import (
	"log"
	"findai/src/apps"
	"findai/src/apps/utils"
	"findai/src/config"
	"findai/src/database"
)

func main() {
	if _, err := config.Init("config.yml"); err != nil {
		log.Fatalf("Failed to initialize config: %v. Make sure config.yml exists and is correctly formatted.", err)
	}

	database.Connect(config.Config.Database.URL)
	utils.SetDB(database.DB().DB)

	apps.Serve()
}

