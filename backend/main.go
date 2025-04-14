package main

import (
	"log"

	"digital-library/backend/config"
	"digital-library/backend/database"
	"digital-library/backend/routes"

	"github.com/gofiber/fiber/v2"
	// Add routes import later
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Connect Database
	database.Connect(cfg)
	// Setup graceful shutdown
	// database.SetupCloseHandler() // Uncomment if you implement signal handling
	defer database.Close()

	app := fiber.New()

	// Setup Routes
	routes.SetupRoutes(app)

	log.Println("Starting server on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
