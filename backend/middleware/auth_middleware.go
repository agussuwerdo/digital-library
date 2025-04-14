package middleware

import (
	"digital-library/backend/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Protected creates a JWT middleware handler
func Protected(cfg *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,         // Specify the algorithm
			Key:    []byte(cfg.JWTSecret), // Get the secret from config
		},
		ErrorHandler: jwtError, // Custom error handler
		// ContextKey: "user", // Optional: Define the key to store the token in c.Locals
	})
}

// jwtError is a custom error handler for the JWT middleware
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or malformed JWT",
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid or expired JWT",
	})
}
