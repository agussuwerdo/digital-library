// @title Digital Library API
// @version 1.0
// @description This is the API documentation for the Digital Library application
// @host https://digital-library-backend.werdev.my.id
// @BasePath /api
package main

import (
	"log"
	"os"

	"digital-library/backend/app"
	"digital-library/backend/database"
	_ "digital-library/backend/docs" // Import generated docs

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading config from environment variables")
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Setup and start the application
	app := app.SetupApp()
	defer database.Close()

	log.Println("Starting server on port " + port + "...")
	log.Fatal(app.Listen(":" + port))
}
