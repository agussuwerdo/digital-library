package main

import (
	"log"

	"digital-library/backend/app"
	"digital-library/backend/database"
)

func main() {
	app := app.SetupApp()
	defer database.Close()

	log.Println("Starting server on port 3001...")
	if err := app.Listen(":3001"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
