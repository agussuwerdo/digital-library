package main

import (
	"log"
	"net/http"
	"os"

	"digital-library/backend/config"
	"digital-library/backend/database"
	"digital-library/backend/routes"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Create a new Fiber app
func createApp() *fiber.App {
	// Load Config
	cfg := config.LoadConfig()

	// Connect Database
	database.Connect(cfg)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Log the error
			log.Printf("Error: %v", err)
			// Return error response
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add recover middleware to catch panics
	app.Use(recover.New())

	// Configure CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Add logger middleware with more detailed configuration
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Setup Routes
	routes.SetupRoutes(app, cfg)

	return app
}

// Handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	app := createApp()
	// Add cleanup handler that executes after the request is complete
	defer func() {
		database.Close()
	}()

	handler := adaptor.FiberApp(app)
	handler.ServeHTTP(w, r)
}

// For local development
func main() {
	app := createApp()
	// Add cleanup handler
	defer database.Close()

	log.Println("Starting server on port 3001...")
	if err := app.Listen(":3001"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
