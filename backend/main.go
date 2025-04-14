package main

import (
	"log"

	"digital-library/backend/config"
	"digital-library/backend/database"
	"digital-library/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Frontend URL
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup Routes
	routes.SetupRoutes(app, cfg)

	log.Println("Starting server on port 3001...")
	if err := app.Listen(":3001"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
