package main

import (
	"log"
	"findai/src/config"
)

func main() {
	if _, err := config.Init("config.yml"); err != nil {
		log.Fatalf("Failed to initialize config: %v. Make sure config.yml exists and is correctly formatted.", err)
	}
	log.Println("Worker started. Implement your worker logic here.")	
}

