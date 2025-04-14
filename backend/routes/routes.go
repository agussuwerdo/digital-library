package routes

import (
	// "digital-library/backend/handlers"
	// "digital-library/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// jwtware "github.com/gofiber/jwt/v3" // Import later when middleware is implemented
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Middleware
	app.Use(logger.New()) // Add basic request logging

	// Public routes
	api := app.Group("/api")
	// api.Post("/login", handlers.Login) // TODO: Implement auth handler

	// --- JWT Protected Routes ---
	// TODO: Apply JWT Middleware here once implemented
	// Example: app.Use(middleware.Protected())

	// Book routes
	book := api.Group("/books")
	// book.Post("/", handlers.CreateBook)    // TODO: Implement book handlers
	// book.Get("/", handlers.GetBooks)
	// book.Get("/:id", handlers.GetBook)
	// book.Put("/:id", handlers.UpdateBook)
	// book.Delete("/:id", handlers.DeleteBook)

	// Lending routes
	lending := api.Group("/lending")
	// lending.Post("/lend", handlers.LendBook)       // TODO: Implement lending handlers
	// lending.Post("/return/:id", handlers.ReturnBook) // :id is lending record id
	// lending.Get("/", handlers.GetLendingRecords)
	// lending.Delete("/:id", handlers.DeleteLendingRecord)

	// Analytics routes
	analytics := api.Group("/analytics")
	// analytics.Get("/most-borrowed", handlers.GetMostBorrowedBooks) // TODO: Implement analytics handlers
	// analytics.Get("/monthly-trends", handlers.GetMonthlyLendingTrends)
	// analytics.Get("/category-distribution", handlers.GetCategoryDistribution)

	// Health Check (optional)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
