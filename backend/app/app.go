package app

import (
	"log"
	"os"
	"strings"

	"digital-library/backend/config"
	"digital-library/backend/database"
	"digital-library/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// QueryParamsMiddleware extracts query parameters from the request
func QueryParamsMiddleware(c *fiber.Ctx) error {
	// Get the original request
	req := c.Request()

	// Log the original request details
	log.Printf("QueryParamsMiddleware - Original URI: %s", string(req.RequestURI()))
	log.Printf("QueryParamsMiddleware - Raw Query: %s", string(req.URI().QueryString()))

	// If there's no query string in the URI, try to get it from the context
	if len(req.URI().QueryString()) == 0 {
		// Try to get query from context
		if queryStr, ok := c.Context().Value("query_string").(string); ok && queryStr != "" {
			// Set the query string in the request
			req.URI().SetQueryString(queryStr)
			log.Printf("QueryParamsMiddleware - Set query from context: %s", queryStr)
		}
	}

	return c.Next()
}

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

	// Add query params middleware
	app.Use(QueryParamsMiddleware)

	// Configure CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	allowOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"} // Default local development URLs
	if frontendURL != "" {
		allowOrigins = append(allowOrigins, frontendURL)
	}

	// Filter out empty strings and join the origins
	validOrigins := make([]string, 0)
	for _, origin := range allowOrigins {
		if origin != "" {
			validOrigins = append(validOrigins, origin)
		}
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(validOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, Access-Control-Allow-Origin",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin",
		MaxAge:           3600,
	}))

	// Add logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Setup Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup Routes
	routes.SetupRoutes(app, cfg)

	return app
}
