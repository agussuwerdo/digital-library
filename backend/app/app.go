package app

import (
	"log"
	"os"

	"digital-library/backend/config"
	"digital-library/backend/database"
	"digital-library/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupApp creates and configures a new Fiber app
func SetupApp() *fiber.App {
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
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, Access-Control-Allow-Origin",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin",
		MaxAge:           3600,
		Next:             nil,
		AllowOriginsFunc: nil,
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
