package routes

import (
	"digital-library/backend/config"     // Import config
	"digital-library/backend/handlers"   // Import handlers
	"digital-library/backend/middleware" // Import middleware
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config) { // Accept config
	// Middleware
	app.Use(logger.New()) // Add basic request logging

	// Public routes
	api := app.Group("/api")
	api.Post("/register", handlers.Register) // Add registration route
	api.Post("/login", handlers.Login(cfg))  // Add login route, pass config

	// Serve Swagger documentation
	api.Get("/apidocs", func(c *fiber.Ctx) error {
		// Read the Swagger JSON file
		swaggerPath := filepath.Join("docs", "swagger.json")
		swaggerData, err := os.ReadFile(swaggerPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read Swagger documentation: " + err.Error(),
			})
		}

		// Set the content type and serve the Swagger JSON
		c.Set("Content-Type", "application/json")
		return c.Send(swaggerData)
	})

	// --- JWT Protected Routes ---
	// Apply JWT middleware to groups below
	protected := api.Group("", middleware.Protected(cfg)) // Create a group with the middleware

	// Book routes (now protected)
	book := protected.Group("/books")
	book.Post("/", handlers.CreateBook)      // Connect CreateBook handler
	book.Get("/", handlers.GetBooks)         // Connect GetBooks handler
	book.Get("/:id", handlers.GetBook)       // Connect GetBook handler
	book.Put("/:id", handlers.UpdateBook)    // Connect UpdateBook handler
	book.Delete("/:id", handlers.DeleteBook) // Connect DeleteBook handler

	// Lending routes (now protected)
	lending := protected.Group("/lending")
	lending.Post("/lend", handlers.LendBook)             // Connect LendBook handler
	lending.Post("/return/:id", handlers.ReturnBook)     // Connect ReturnBook handler
	lending.Get("/", handlers.GetLendingRecords)         // Connect GetLendingRecords handler
	lending.Delete("/:id", handlers.DeleteLendingRecord) // Connect DeleteLendingRecord handler

	// Analytics routes (now protected)
	analytics := protected.Group("/analytics")
	analytics.Get("/most-borrowed", handlers.GetMostBorrowedBooks)            // Connect analytics handler
	analytics.Get("/monthly-trends", handlers.GetMonthlyLendingTrends)        // Connect analytics handler
	analytics.Get("/category-distribution", handlers.GetCategoryDistribution) // Connect analytics handler

	// Health Check (optional - public)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
